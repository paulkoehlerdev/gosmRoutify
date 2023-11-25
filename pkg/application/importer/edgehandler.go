package importer

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinatelist"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/nodetags"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/osmid"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
)

type edgeHandler func(fromID, toID osmid.OsmID, nodeList coordinatelist.CoordinateList, tags []nodetags.NodeTags, way osmpbfreaderdata.Way)
