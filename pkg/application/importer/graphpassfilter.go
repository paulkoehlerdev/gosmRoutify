package importer

import "github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"

type GraphPassFilter struct {
}

func NewGraphPassFilter() osmdatarepository.OsmDataFilter {
	return &GraphPassFilter{}
}

func (w *GraphPassFilter) FilterNodes() bool {
	return false
}

func (w *GraphPassFilter) FilterWays() bool {
	return false
}

func (w *GraphPassFilter) FilterRelations() bool {
	return true
}
