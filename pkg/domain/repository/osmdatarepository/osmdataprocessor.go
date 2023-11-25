package osmdatarepository

import "github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"

type OsmDataProcessor interface {
	ProcessNode(node osmpbfreaderdata.Node)
	ProcessWay(way osmpbfreaderdata.Way)
	ProcessRelation(relation osmpbfreaderdata.Relation)

	OnFinish()
}
