package osmDataPreprocessorService

import "github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/osmpbfData"

type OsmDataPreprocessorService interface {
	Filter(data any) bool
}

type impl struct {
}

func New() OsmDataPreprocessorService {
	return &impl{}
}

// Filter Returns true if the data should be filtered out
func (i impl) Filter(data any) bool {
	switch data.(type) {
	case *osmpbfData.Node:
		return i.filterNode(data.(*osmpbfData.Node))
	case *osmpbfData.Way:
		return i.filterWay(data.(*osmpbfData.Way))
	case *osmpbfData.Relation:
		return i.filterRelation(data.(*osmpbfData.Relation))
	default:
		return true
	}
}

// filterNode Returns true if the data should be filtered out
func (i impl) filterNode(node *osmpbfData.Node) bool {
	// filter non-street nodes
	if _, ok := node.Tags["highway"]; !ok {
		return true
	}
	return false
}

// filterWay Returns true if the data should be filtered out
func (i impl) filterWay(way *osmpbfData.Way) bool {
	// filter non-street ways
	if _, ok := way.Tags["highway"]; !ok {
		return true
	}
	return false
}

// filterRelation Returns true if the data should be filtered out
func (i impl) filterRelation(relation *osmpbfData.Relation) bool {
	// filter all relations for now, as we don't need them for cars
	return true
}
