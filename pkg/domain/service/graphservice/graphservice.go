package graphservice

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinatelist"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/nodetags"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
)

type GraphService interface {
	AddEdge(fromID, toID int64, nodeList coordinatelist.CoordinateList, tags []nodetags.NodeTags, way *osmpbfreaderdata.Way)
}
