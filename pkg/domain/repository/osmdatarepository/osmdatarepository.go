package osmdatarepository

import (
	"errors"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
	"io"
)

type OsmDataRepository interface {
	Process(file string, processor OsmDataProcessor, filter OsmDataFilter) error
}

type impl struct {
	parallelization int
}

func New(parallelization int) OsmDataRepository {
	return &impl{
		parallelization: parallelization,
	}
}

func (o *impl) Process(file string, processor OsmDataProcessor, filter OsmDataFilter) error {
	reader := &osmReader{
		parallelization: o.parallelization,
	}
	err := reader.Read(file, filter)
	if err != nil {
		return fmt.Errorf("error while reading file: %s", err.Error())
	}

	for {
		data, err := reader.Next()
		if errors.Is(err, io.EOF) {
			processor.OnFinish()
			return nil
		}

		if err != nil {
			return fmt.Errorf("error while reading next data: %s", err.Error())
		}

		switch v := data.(type) {
		case osmpbfreaderdata.Node:
			processor.ProcessNode(v)
		case osmpbfreaderdata.Way:
			processor.ProcessWay(v)
		case osmpbfreaderdata.Relation:
			processor.ProcessRelation(v)
		default:
			return fmt.Errorf("unknown data type: %T", v)
		}
	}
}
