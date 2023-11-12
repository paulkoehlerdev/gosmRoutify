package osmDataRepository

import (
	"errors"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmDataRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf"
	"io"
	"os"
)

type OsmDataRepository struct {
	decoder         osmpbf.Decoder
	file            *os.File
	parallelization int
}

func New(parallelization int) osmDataRepository.OsmDataRepository {
	return &OsmDataRepository{
		decoder:         nil,
		parallelization: parallelization,
	}
}

func (o *OsmDataRepository) Read(filePath string) error {
	var err error
	o.file, err = os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error while opening file: %s", err.Error())
	}

	o.decoder = osmpbf.New(o.file)
	err = o.decoder.Start(osmpbf.ProcsCount(o.parallelization))
	if err != nil {
		return fmt.Errorf("error while starting decoder: %s", err.Error())
	}
	return nil
}

func (o *OsmDataRepository) Next() (any, error) {
	if o.decoder == nil {
		return nil, errors.New("no decoder loaded: you need to call Read() before you can call Next()")
	}

	data, err := o.decoder.Decode()
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("error while decoding: %s", err.Error())
	}
	return data, err
}

func (o *OsmDataRepository) Stop() {
	if o.file != nil {
		_ = o.file.Close()
	}
}
