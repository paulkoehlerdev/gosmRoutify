package importer

import "github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"

type WayPassFilter struct {
}

func NewWayPassFilter() osmdatarepository.OsmDataFilter {
	return &WayPassFilter{}
}

func (w *WayPassFilter) FilterNodes() bool {
	return true
}

func (w *WayPassFilter) FilterWays() bool {
	return false
}

func (w *WayPassFilter) FilterRelations() bool {
	return true
}
