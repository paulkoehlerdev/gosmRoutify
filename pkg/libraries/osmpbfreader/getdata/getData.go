package getdata

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmproto"
)

func GetData(blob *osmproto.Blob) ([]byte, error) {
	switch blob.Data.(type) {
	case *osmproto.Blob_Raw:
		return blob.GetRaw(), nil

	case *osmproto.Blob_ZlibData:
		r, err := zlib.NewReader(bytes.NewReader(blob.GetZlibData()))
		if err != nil {
			return nil, err
		}
		buf := bytes.NewBuffer(make([]byte, 0, blob.GetRawSize()+bytes.MinRead))
		_, err = buf.ReadFrom(r)
		if err != nil {
			return nil, err
		}
		if buf.Len() != int(blob.GetRawSize()) {
			err = fmt.Errorf("raw blob data size %d but expected %d", buf.Len(), blob.GetRawSize())
			return nil, err
		}
		return buf.Bytes(), nil

	default:
		return nil, fmt.Errorf("unhandled blob data type %T", blob.Data)
	}
}
