package osmdatarepository

import (
	"errors"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/filter"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
	"io"
)

type OsmDataProcessor interface {
	ProcessNode(node osmpbfreaderdata.Node)
	ProcessWay(way osmpbfreaderdata.Way)
	ProcessRelation(relation osmpbfreaderdata.Relation)

	OnFinish()
}

type OsmDataFilter interface {
	filter.Filter
}

type BinaryOsmDataFilter struct {
	filterNodes     bool
	filterWays      bool
	filterRelations bool
}

func NewBinaryOsmDataFilter(filterNodes bool, filterWays bool, filterRelations bool) OsmDataFilter {
	return &BinaryOsmDataFilter{
		filterNodes:     filterNodes,
		filterWays:      filterWays,
		filterRelations: filterRelations,
	}
}

func (i *BinaryOsmDataFilter) FilterNodes() bool {
	return i.filterNodes
}

func (i *BinaryOsmDataFilter) FilterWays() bool {
	return i.filterWays
}

func (i *BinaryOsmDataFilter) FilterRelations() bool {
	return i.filterRelations
}

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
