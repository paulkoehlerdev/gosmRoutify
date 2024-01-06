package importer

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/nodetype"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinate"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/osmid"
)

type segmentNode struct {
	nodeType   nodetype.NodeType
	osmID      osmid.OsmID
	coordinate coordinate.Coordinate
	tags       map[string]string
}
