package osmpbf

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/blobReader"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/dataDecoder"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/osmproto"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/valueErrPair"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/workerpool"
	"io"
)

type ProcsCount int

type Decoder interface {
	Start(ProcsCount) error
	Decode() (any, error)
}

const standartPrimitiveFeatureCount = 8000

type impl struct {
	blobReader blobReader.BlobReader
	serializer chan valueErrPair.Pair[any]
	done       chan struct{}
}

func New(reader io.Reader) Decoder {
	return &impl{
		blobReader: blobReader.NewBlobReader(reader),
		serializer: make(chan valueErrPair.Pair[any], standartPrimitiveFeatureCount),
	}
}

func (i *impl) Start(count ProcsCount) error {
	i.done = make(chan struct{})
	blobs := make(chan valueErrPair.Pair[*osmproto.Blob], count)
	pool := workerpool.New[valueErrPair.Pair[[]interface{}]](workerpool.ProcsCount(count))
	pool.Start()

	go i.blobReader.Read(blobs)
	go i.decoderHandler(pool, blobs)
	go i.decodedBlobHandler(pool)

	go func() {
		<-i.done
		pool.Stop()
	}()

	return nil
}

func (i *impl) decodedBlobHandler(pool workerpool.Pool[valueErrPair.Pair[[]interface{}]]) {
	for {
		result, err := pool.Result()
		if err != nil {
			i.serializer <- valueErrPair.Pair[any]{Err: err}
			return
		}

		if result.Err != nil {
			i.serializer <- valueErrPair.Pair[any]{Err: result.Err}
			return
		}

		for _, object := range result.Value {
			i.serializer <- valueErrPair.Pair[any]{Value: object}
		}
	}
}

func (i *impl) decoderHandler(pool workerpool.Pool[valueErrPair.Pair[[]interface{}]], blobs chan valueErrPair.Pair[*osmproto.Blob]) {
	for {
		blob, ok := <-blobs
		if !ok {
			return
		}

		if blob.Err != nil {
			pool.Submit(func(int) valueErrPair.Pair[[]interface{}] {
				return valueErrPair.Pair[[]interface{}]{Err: blobReader.EOFor(blob.Err, fmt.Errorf("error reading blob: %s", blob.Err.Error()))}
			})
			return
		}

		pool.Submit(func(int) valueErrPair.Pair[[]interface{}] {
			data, err := dataDecoder.NewDataDecoder().Decode(blob.Value)
			if err != nil {
				return valueErrPair.Pair[[]interface{}]{Err: err}
			}
			return valueErrPair.Pair[[]interface{}]{Value: data}
		})
	}
}

func (i *impl) Decode() (any, error) {
	pair, ok := <-i.serializer
	if !ok {
		close(i.done)
		return nil, io.EOF
	}
	return pair.Value, pair.Err
}
