package blobReader

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/dataDecoder"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/osmpbfData"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/osmproto"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/valueErrPair"
	"io"
	"sync"
)

const expectedBlobHeaderType = "OSMHeader"
const expectedBlobDataType = "OSMData"
const maxBlobHeaderSize = 32 * 1024
const maxBlobSize = 64 * 1024 * 1024
const degMultiplier = 1e-9

type BlobReader interface {
	Read(blobs chan valueErrPair.Pair[*osmproto.Blob])
	Header() (*osmpbfData.Header, error)
}

type impl struct {
	reader        io.Reader
	buffer        *bytes.Buffer
	osmHeaderOnce sync.Once
	osmHeader     *osmpbfData.Header
}

func NewBlobReader(reader io.Reader) BlobReader {
	return &impl{
		reader: reader,
		buffer: bytes.NewBuffer(make([]byte, 0, maxBlobSize)),
	}
}

func (i *impl) Header() (*osmpbfData.Header, error) {
	return i.osmHeader, i.readHeaderBlock()
}

func (i *impl) Read(blobs chan valueErrPair.Pair[*osmproto.Blob]) {
	err := i.readHeaderBlock()
	if err != nil {
		blobs <- valueErrPair.Pair[*osmproto.Blob]{Err: fmt.Errorf("error reading header: %s", err.Error())}
		close(blobs)
		return
	}

	for {
		blobHeader, blob, err := i.readBlock()
		if err != nil || blobHeader.GetType() != expectedBlobDataType {
			if blobHeader.GetType() != expectedBlobDataType && err == nil {
				err = fmt.Errorf("invalid type: \"%s\"", blobHeader.GetType())
			}
			blobs <- valueErrPair.Pair[*osmproto.Blob]{Err: EOFor(err, fmt.Errorf("error reading blob: %s", err.Error()))}
			close(blobs)
			return
		}

		blobs <- valueErrPair.Pair[*osmproto.Blob]{Value: blob}
	}
}

func (i *impl) readHeaderBlock() error {
	var err error

	i.osmHeaderOnce.Do(func() {
		var blobHeader *osmproto.BlobHeader
		var blob *osmproto.Blob

		blobHeader, blob, err = i.readBlock()
		if err != nil || blobHeader.GetType() != expectedBlobHeaderType {
			if blobHeader.GetType() != expectedBlobHeaderType {
				err = fmt.Errorf("invalid type: %s", blobHeader.GetType())
			}
			return
		}

		i.osmHeader, err = dataDecoder.DecodeHeaderBlock(blob)
	})

	return nil
}

func (i *impl) readBlock() (*osmproto.BlobHeader, *osmproto.Blob, error) {
	blobHeaderSize, err := i.readBlobHeaderSize()
	if err != nil {
		return nil, nil, EOFor(err, fmt.Errorf("error reading blob header size: %s", err.Error()))
	}

	blobHeader, err := i.readBlobHeader(blobHeaderSize)
	if err != nil {
		return nil, nil, EOFor(err, fmt.Errorf("error reading blob header: %s", err.Error()))
	}

	blob, err := i.readBlob(blobHeader)
	if err != nil {
		return nil, nil, EOFor(err, fmt.Errorf("error reading blob: %s", err.Error()))
	}

	return blobHeader, blob, nil
}

func (i *impl) readBlobHeaderSize() (uint32, error) {
	i.buffer.Reset()
	if _, err := io.CopyN(i.buffer, i.reader, 4); err != nil {
		return 0, EOFor(err, fmt.Errorf("error reading blob header size: %s", err.Error()))
	}

	size := binary.BigEndian.Uint32(i.buffer.Bytes())

	if size >= maxBlobHeaderSize {
		return 0, errors.New("blobHeader size >= 64Kb")
	}
	return size, nil
}

func (i *impl) readBlobHeader(size uint32) (*osmproto.BlobHeader, error) {
	i.buffer.Reset()
	if _, err := io.CopyN(i.buffer, i.reader, int64(size)); err != nil {
		return nil, EOFor(err, fmt.Errorf("error reading blob header: %s", err.Error()))
	}

	blobHeader := new(osmproto.BlobHeader)
	if err := proto.Unmarshal(i.buffer.Bytes(), blobHeader); err != nil {
		return nil, fmt.Errorf("error unmarshalling blob header: %s", err.Error())
	}

	if blobHeader.GetDatasize() >= maxBlobSize {
		return nil, errors.New("blob size >= 32Mb")
	}

	return blobHeader, nil
}

func (i *impl) readBlob(blobHeader *osmproto.BlobHeader) (*osmproto.Blob, error) {
	i.buffer.Reset()
	if _, err := io.CopyN(i.buffer, i.reader, int64(blobHeader.GetDatasize())); err != nil {
		return nil, EOFor(err, fmt.Errorf("error reading blob: %s", err.Error()))
	}

	blob := new(osmproto.Blob)
	if err := proto.Unmarshal(i.buffer.Bytes(), blob); err != nil {
		return nil, fmt.Errorf("error unmarshalling blob: %s", err.Error())
	}

	return blob, nil
}

func EOFor(cmp error, err error) error {
	if errors.Is(cmp, io.EOF) {
		return io.EOF
	}
	return err
}
