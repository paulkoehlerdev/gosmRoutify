package osmdatarepository

import (
	"errors"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/filter"
	"io"
	"os"
)

type osmReader struct {
	parallelization int
	decoder         osmpbfreader.Decoder
	file            *os.File
}

func (o *osmReader) Stop() {
	if o.file != nil {
		_ = o.file.Close()
	}
}

func (o *osmReader) Read(filePath string, filter filter.Filter) error {
	var err error
	o.file, err = os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err.Error())
	}

	o.decoder = osmpbfreader.New(o.file)
	err = o.decoder.Start(osmpbfreader.ProcsCount(o.parallelization), filter)
	if err != nil {
		return fmt.Errorf("error while starting decoder: %s", err.Error())
	}
	return nil
}

func (o *osmReader) Next() (any, error) {
	if o.decoder == nil {
		return nil, errors.New("no decoder loaded: you need to call Read() before you can call Next()")
	}

	data, err := o.decoder.Decode()
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("error while decoding: %s", err.Error())
	}
	return data, err
}
