package graphtile

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinate"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/nodetags"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/osmid"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/kdtree"
)

type GraphTile struct {
	Graph     map[osmid.OsmID]map[osmid.OsmID]int
	Edges     []EdgeInfo
	NodeInfos NodeInfos
}

func New() *GraphTile {
	return &GraphTile{
		Graph:     make(map[osmid.OsmID]map[osmid.OsmID]int),
		NodeInfos: NodeInfos{Tree: kdtree.New[osmid.OsmID]()},
	}
}

type NodeInfos struct {
	Tree kdtree.KdTree[osmid.OsmID]
}

func (n NodeInfos) FindNearest(coordinate coordinate.Coordinate) *osmid.OsmID {
	return n.Tree.SearchNearest(coordinate.Lat(), coordinate.Lon())
}

func (n NodeInfos) Insert(coordinate coordinate.Coordinate, id osmid.OsmID) {
	n.Tree.Insert(coordinate.Lat(), coordinate.Lon(), id)
}

type EdgeInfo struct {
	StartID  osmid.OsmID
	EndID    osmid.OsmID
	Tags     map[string]string
	NodeTags []nodetags.NodeTags
	Nodes    []coordinate.Coordinate
}
