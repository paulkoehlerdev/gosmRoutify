package importer

import "github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"

type GraphPassFilter struct {
}

func NewWayGraphPassFilter() osmdatarepository.OsmDataFilter {
	return &GraphPassFilter{}
}

func (w *GraphPassFilter) FilterNodes() bool {
	return true
}

func (w *GraphPassFilter) FilterWays() bool {
	return false
}

func (w *GraphPassFilter) FilterRelations() bool {
	return true
}
