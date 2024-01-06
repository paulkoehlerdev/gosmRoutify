package router

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/graphservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/weightingservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinate"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/osmid"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/astar"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geodistance"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"math"
	"time"
)

type Application interface {
	FindRoute(start coordinate.Coordinate, end coordinate.Coordinate) ([]coordinate.Coordinate, error)
}

type impl struct {
	logger           logging.Logger
	graphService     graphservice.GraphService
	weightingService weightingservice.WeightingService
}

func New(graphService graphservice.GraphService, weightingService weightingservice.WeightingService, logger logging.Logger) Application {
	return &impl{
		logger:           logger,
		graphService:     graphService,
		weightingService: weightingService,
	}
}

type GraphInfo struct {
	coo coordinate.Coordinate
}

func (i *impl) FindRoute(start coordinate.Coordinate, end coordinate.Coordinate) ([]coordinate.Coordinate, error) {
	startTime := time.Now()
	startTile := i.graphService.GetTile(start)
	endTile := i.graphService.GetTile(end)
	i.logger.Debug().Msgf("calculated nearest tiles in %s", time.Since(startTime).String())

	if startTile == nil {
		return nil, fmt.Errorf("start tile not found")
	}

	if endTile == nil {
		return nil, fmt.Errorf("end tile not found")
	}

	startTime = time.Now()
	startID := i.graphService.FindNearest(start)
	endID := i.graphService.FindNearest(end)
	i.logger.Debug().Msgf("calculated nearest nodes in %s", time.Since(startTime).String())

	if startID == nil {
		return nil, fmt.Errorf("start node not found")
	}

	if endID == nil {
		return nil, fmt.Errorf("end node not found")
	}

	infoMap := make(map[osmid.OsmID]*GraphInfo)

	infoMap[*startID] = &GraphInfo{
		coo: start,
	}

	infoMap[*endID] = &GraphInfo{
		coo: end,
	}

	startTime = time.Now()
	path, length, err := astar.AStar(*startID, *endID, i.GetNeighbors(infoMap), GetHeuristic(*endID, infoMap))
	if err != nil {
		return nil, fmt.Errorf("error while routing: %s", err.Error())
	}
	i.logger.Debug().Msgf("calculated route in %s", time.Since(startTime).String())

	duration := time.Duration(length) * time.Second
	i.logger.Debug().Msgf("calculated route takes %s", duration.String())

	var route []coordinate.Coordinate
	for _, id := range path {
		info, ok := infoMap[id]
		if !ok {
			continue
		}
		route = append(route, info.coo)
	}

	return route, nil
}

func GetHeuristic(end osmid.OsmID, infoMap map[osmid.OsmID]*GraphInfo) func(id osmid.OsmID) float64 {
	return func(node osmid.OsmID) float64 {
		nodeInfo, ok := infoMap[node]
		if !ok {
			return math.Inf(1)
		}

		endInfo, ok := infoMap[end]
		if !ok {
			return math.Inf(1)
		}

		distance := geodistance.CalcDistanceInMeters(nodeInfo.coo, endInfo.coo)
		return distance / 1000 //(highwaytype.Motorway.DefaultMaxSpeed() * 2)
	}
}

func (i *impl) GetNeighbors(infoMap map[osmid.OsmID]*GraphInfo) func(info osmid.OsmID) map[osmid.OsmID]float64 {
	return func(id osmid.OsmID) map[osmid.OsmID]float64 {
		info, ok := infoMap[id]
		if !ok {
			return nil
		}

		tile := i.graphService.GetTile(info.coo)
		if tile == nil {
			return nil
		}

		neighbors := make(map[osmid.OsmID]float64)

		for neighbourID, edgeID := range tile.Graph[id] {
			edge := tile.Edges[edgeID]

			neighborCoo := edge.Nodes[len(edge.Nodes)-1]
			if edge.StartID == neighbourID {
				neighborCoo = edge.Nodes[0]
			}

			infoMap[neighbourID] = &GraphInfo{
				coo: neighborCoo,
			}

			neighbors[neighbourID] = i.weightingService.Weight(edge)
		}

		return neighbors
	}
}
