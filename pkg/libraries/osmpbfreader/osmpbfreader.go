package osmpbfreader

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/blobreader"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/datadecoder"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmproto"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/valueerrpair"
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
	blobReader blobreader.BlobReader
	serializer chan valueerrpair.Pair[any]
	done       chan struct{}
}

func New(reader io.Reader) Decoder {
	return &impl{
		blobReader: blobreader.NewBlobReader(reader),
		serializer: make(chan valueerrpair.Pair[any], standartPrimitiveFeatureCount),
	}
}

func (i *impl) Start(count ProcsCount) error {
	i.done = make(chan struct{})
	blobs := make(chan valueerrpair.Pair[*osmproto.Blob], count)
	pool := workerpool.New[valueerrpair.Pair[[]interface{}]](workerpool.ProcsCount(count))
	pool.Start()

	go i.blobReader.Read(blobs)
	go decoderHandler(pool, blobs)
	go i.decodedBlobHandler(pool)

	go func() {
		<-i.done
		pool.Stop()
	}()

	return nil
}

func (i *impl) decodedBlobHandler(pool workerpool.Pool[valueerrpair.Pair[[]interface{}]]) {
	for {
		result, err := pool.Result()
		if err != nil {
			i.serializer <- valueerrpair.Pair[any]{Err: err}
			return
		}

		if result.Err != nil {
			i.serializer <- valueerrpair.Pair[any]{Err: result.Err}
			return
		}

		for _, object := range result.Value {
			i.serializer <- valueerrpair.Pair[any]{Value: object}
		}
	}
}

func decoderHandler(pool workerpool.Pool[valueerrpair.Pair[[]interface{}]], blobs chan valueerrpair.Pair[*osmproto.Blob]) {
	for {
		blob, ok := <-blobs
		if !ok {
			return
		}

		if blob.Err != nil {
			pool.Submit(func(int) valueerrpair.Pair[[]interface{}] {
				return valueerrpair.Pair[[]interface{}]{Err: blobreader.EOFor(blob.Err, fmt.Errorf("error reading blob: %s", blob.Err.Error()))}
			})
			return
		}

		pool.Submit(func(int) valueerrpair.Pair[[]interface{}] {
			data, err := datadecoder.NewDataDecoder().Decode(blob.Value)
			if err != nil {
				return valueerrpair.Pair[[]interface{}]{Err: err}
			}
			return valueerrpair.Pair[[]interface{}]{Value: data}
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
