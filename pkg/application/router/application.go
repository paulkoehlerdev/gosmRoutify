package router

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/graphService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/nodeService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/astar"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geojson"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"time"
)

type Application interface {
	FindRoute(start geojson.Point, end geojson.Point) ([]geojson.Point, error)
}

const (
	maxVisitedNodes = 500000
)

type impl struct {
	logger       logging.Logger
	graphService graphService.GraphService
	nodeService  nodeService.NodeService
}

func New(graphService graphService.GraphService, nodeService nodeService.NodeService, logger logging.Logger) Application {
	return &impl{
		logger:       logger,
		graphService: graphService,
		nodeService:  nodeService,
	}
}

func (i *impl) FindRoute(start geojson.Point, end geojson.Point) ([]geojson.Point, error) {
	startTime := time.Now()

	startNode, err := i.graphService.GetNearestNode(start.Lon(), start.Lat())
	if err != nil {
		return nil, fmt.Errorf("error while finding nearest node at start: %s", err.Error())
	}

	endNode, err := i.graphService.GetNearestNode(end.Lon(), end.Lat())
	if err != nil {
		return nil, fmt.Errorf("error while finding nearest node at end: %s", err.Error())
	}

	i.logger.Debug().Msgf("calculated nearest node in %s", time.Since(startTime).String())

	if endNode == nil {
		return nil, fmt.Errorf("no end node found")
	}

	path, length, err := astar.AStar[int64, float64](startNode.OsmID, endNode.OsmID, i.graphService.GetEdges(*endNode), i.graphService.GetHeuristic(*endNode), maxVisitedNodes)
	if err != nil {
		return nil, fmt.Errorf("error while routing: %s", err.Error())
	}

	nodePoints, lengthInMeters, err := i.graphService.CalculatePathInformation(path)
	if err != nil {
		return nil, fmt.Errorf("error while building geojson line: %s", err.Error())
	}

	i.logger.Debug().Msgf("calculated route in %s", time.Since(startTime).String())
	i.logger.Debug().Msgf("route length (time): %s", time.Duration(length)*time.Second)
	i.logger.Debug().Msgf("route length (distance): %.3fkm", lengthInMeters/1000)
	i.logger.Debug().WithAttrs("elements", path).Msgf("route length (elements): %v", len(path))

	nodePoints = append(
		[]geojson.Point{start},
		nodePoints...,
	)

	nodePoints = append(
		nodePoints,
		end,
	)

	return nodePoints, nil
}
