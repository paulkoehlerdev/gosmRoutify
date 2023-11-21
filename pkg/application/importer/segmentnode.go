package importer

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/nodetype"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinate"
)

type segmentNode struct {
	nodeType   nodetype.NodeType
	osmID      int64
	coordinate coordinate.Coordinate
	tags       map[string]string
}
