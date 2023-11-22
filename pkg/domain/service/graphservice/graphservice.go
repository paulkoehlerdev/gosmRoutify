package graphservice

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/graphtile"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/tilerepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinate"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinatelist"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/nodetags"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/osmid"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
)

type GraphService interface {
	AddEdge(fromID, toID osmid.OsmID, nodeList coordinatelist.CoordinateList, tags []nodetags.NodeTags, way *osmpbfreaderdata.Way)
	GetTile(coo coordinate.Coordinate) *graphtile.GraphTile
}

type impl struct {
	tileRepo tilerepository.TileRepository
	logger   logging.Logger
}

func New(tileRepo tilerepository.TileRepository, logger logging.Logger) GraphService {
	return &impl{
		tileRepo: tileRepo,
		logger:   logger,
	}
}

func (i *impl) AddEdge(fromID, toID osmid.OsmID, nodeList coordinatelist.CoordinateList, tags []nodetags.NodeTags, way *osmpbfreaderdata.Way) {
	startCoo := nodeList.Get(0)

	tile := i.tileRepo.GetTile(startCoo)
	if tile == nil {
		tile = graphtile.New()
	}

	edgeInfo := graphtile.EdgeInfo{
		StartID:  fromID,
		EndID:    toID,
		Tags:     way.Tags,
		NodeTags: tags,
		Nodes:    nodeList.ToCoordinateArray(),
	}

	if _, ok := tile.Graph[fromID]; !ok {
		tile.Graph[fromID] = make(map[osmid.OsmID]int)
	}

	if _, ok := tile.Graph[toID]; !ok {
		tile.Graph[toID] = make(map[osmid.OsmID]int)
	}

	edgeID := len(tile.Edges)
	tile.Edges = append(tile.Edges, edgeInfo)

	if isReversed(edgeInfo) {
		tile.Graph[toID][fromID] = edgeID
	} else if !isReversible(edgeInfo) {
		tile.Graph[fromID][toID] = edgeID
	} else {
		tile.Graph[fromID][toID] = edgeID
		tile.Graph[toID][fromID] = edgeID
	}

	tile.NodeInfos.Insert(startCoo, fromID)
	tile.NodeInfos.Insert(nodeList.Get(nodeList.Len()-1), toID)

	i.tileRepo.SetTile(startCoo, tile)
}

func isReversed(info graphtile.EdgeInfo) bool {
	value, ok := info.Tags["oneway"]
	if !ok {
		return false
	}

	return value == "-1"
}

func isReversible(info graphtile.EdgeInfo) bool {
	value, ok := info.Tags["oneway"]
	if !ok {
		return true
	}

	return value != "yes" && value != "1" && value != "true"
}

func (i *impl) GetTile(coo coordinate.Coordinate) *graphtile.GraphTile {
	return i.tileRepo.GetTile(coo)
}
