package graphRepository

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/graph"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/osmpbfData"
)

type GraphRepository interface {
	AddWay(ways *osmpbfData.Way, tID graph.TileID) graph.GraphID
	AddIntersection(node *osmpbfData.Node, a graph.GraphID, b graph.GraphID) graph.GraphID
	GetNode(gID graph.GraphID) (*osmpbfData.Node, error)
	GetWay(gID graph.GraphID) (*osmpbfData.Way, error)
}
