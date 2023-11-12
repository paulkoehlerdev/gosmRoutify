package graphRepository

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/graph"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/osmpbfData"
)

type tile struct {
	Data  []any
	Graph map[graph.GraphID]map[graph.GraphID]graph.GraphID
}

func (t *tile) getNode(gID graph.GraphID) (*osmpbfData.Node, error) {
	node, ok := t.Data[gID.ObjectID()].(*osmpbfData.Node)
	if !ok {
		return nil, fmt.Errorf("object is not a node")
	}
	return node, nil
}

func (t *tile) getWay(gID graph.GraphID) (*osmpbfData.Way, error) {
	node, ok := t.Data[gID.ObjectID()].(*osmpbfData.Way)
	if !ok {
		return nil, fmt.Errorf("object is not a way")
	}
	return node, nil
}

func (t *tile) addNode(tID graph.TileID, lID graph.LevelID, node *osmpbfData.Node) graph.GraphID {
	return t.add(tID, lID, node)
}

func (t *tile) addWay(tID graph.TileID, lID graph.LevelID, way *osmpbfData.Way) graph.GraphID {
	gID := t.add(tID, lID, way)
	return gID
}

func (t *tile) add(tID graph.TileID, lID graph.LevelID, data any) graph.GraphID {
	oID := graph.ObjectID(len(t.Data))
	t.Data = append(t.Data, data)
	return graph.NewGraphID(lID, tID, oID)
}

func (t *tile) addRelation(a graph.GraphID, b graph.GraphID, connector graph.GraphID) {
	if t.Graph == nil {
		t.Graph = make(map[graph.GraphID]map[graph.GraphID]graph.GraphID)
	}

	if _, ok := t.Graph[a]; !ok {
		t.Graph[a] = make(map[graph.GraphID]graph.GraphID)
	}
	t.Graph[a][b] = connector

	if _, ok := t.Graph[b]; !ok {
		t.Graph[b] = make(map[graph.GraphID]graph.GraphID)
	}
	t.Graph[b][a] = connector
}
