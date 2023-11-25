package importer

import "github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"

type WayPassFilter2 struct {
}

func NewWayPassFilter2() osmdatarepository.OsmDataFilter {
	return &WayPassFilter2{}
}

func (n *WayPassFilter2) FilterNodes() bool {
	return true
}

func (n *WayPassFilter2) FilterWays() bool {
	return false
}

func (n *WayPassFilter2) FilterRelations() bool {
	return true
}
