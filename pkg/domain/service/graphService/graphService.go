package graphService

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/crossing"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/node"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/way"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/crossingRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/nodeRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/wayRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/weightRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/arrayutil"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geojson"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/sphericmath"
	"math"
)

const (
	nearNodesApproxDistance = 0.001 // approx. 1km
	defaultVehicleType      = weightRepository.Car
)

type GraphService interface {
	GetEdges(end node.Node) func(prevId, id int64) map[int64]float64
	GetHeuristic(end node.Node) func(id int64) float64
	CalculatePathInformation(path []int64) (way []geojson.Point, lengthInMeters float64, err error)
	GetNearestNode(lat float64, lon float64) (*node.Node, error)
}

type impl struct {
	nodeRepository     nodeRepository.NodeRepository
	crossingRepository crossingRepository.CrossingRepository
	wayRepository      wayRepository.WayRepository
	weightRepository   weightRepository.WeightRepository
	logger             logging.Logger

	visitedNodes int
}

func New(nodeRepository nodeRepository.NodeRepository, crossingRepository crossingRepository.CrossingRepository, wayRepository wayRepository.WayRepository, weightRepository weightRepository.WeightRepository, logger logging.Logger) GraphService {
	return &impl{
		nodeRepository:     nodeRepository,
		crossingRepository: crossingRepository,
		wayRepository:      wayRepository,
		weightRepository:   weightRepository,
		logger:             logger,
	}
}

func (i *impl) GetEdges(end node.Node) func(prevId, id int64) map[int64]float64 {
	return func(prevId int64, id int64) map[int64]float64 {
		return i.getEdges(prevId, id, end)
	}
}

func (i *impl) getEdges(prevId, id int64, end node.Node) map[int64]float64 {
	ways, err := i.wayRepository.SelectWaysFromNode(id)
	if err != nil {
		i.logger.Error().Msgf("error while selecting ways from node: %s", err.Error())
		return make(map[int64]float64)
	}

	prevNode, err := i.nodeRepository.SelectNodeFromID(prevId)
	if err != nil {
		i.logger.Error().Msgf("error while selecting node from id: %s", err.Error())
	}

	out := make(map[int64]float64)
	for _, w := range ways {
		if !i.weightRepository.IsWayAllowed(*w, defaultVehicleType) {
			continue
		}

		crossings, err := i.crossingRepository.SelectCrossingsFromWayID(w.OsmID)
		if err != nil {
			i.logger.Error().Msgf("error while selecting nodes from way: %s", err.Error())
			continue
		}

		var fromCrossing *crossing.Crossing
		for _, n := range crossings {

			if n.OsmID == id {
				fromCrossing = n
				break
			}
		}

		if fromCrossing == nil {
			i.logger.Error().Msgf("from crossing not found")
			continue
		}

		weights := i.weightRepository.CalculateWeights(prevNode, fromCrossing, w, crossings, end, defaultVehicleType)
		for k, v := range weights {
			if prevV, ok := out[k]; ok && prevV < v {
				continue
			}
			out[k] = v
		}
	}

	return out
}

func (i *impl) GetHeuristic(end node.Node) func(id int64) float64 {
	return func(nodeId int64) float64 {
		node, err := i.nodeRepository.SelectNodeFromID(nodeId)
		if err != nil {
			i.logger.Error().Msgf("error while selecting node from id: %s", err.Error())
			return 0
		}

		return sphericmath.CalcDistanceInMeters(
			sphericmath.NewPoint(end.Lat, end.Lon),
			sphericmath.NewPoint(node.Lat, node.Lon),
		) * (i.weightRepository.MaximumWayFactor(defaultVehicleType) * 2)
	}
}

func (i *impl) CalculatePathInformation(path []int64) (outPath []geojson.Point, lengthInMeters float64, err error) {
	prevNode, err := i.nodeRepository.SelectNodeFromID(path[0])
	if err != nil {
		return nil, 0.0, fmt.Errorf("error while selecting node from id: %s", err.Error())
	}

	var points []geojson.Point

	for _, nodeId := range path[1:] {
		n, err := i.nodeRepository.SelectNodeFromID(nodeId)
		if err != nil {
			return nil, 0.0, fmt.Errorf("error while selecting node from id: %s", err.Error())
		}

		ways, err := i.wayRepository.SelectWaysFromTwoNodeIDs(prevNode.OsmID, n.OsmID)
		if err != nil {
			return nil, 0.0, fmt.Errorf("error while selecting ways from two nodes: %s", err.Error())
		}

		if len(ways) == 0 {
			return nil, 0.0, fmt.Errorf("no way found between node %d and %d", prevNode.OsmID, n.OsmID)
		}

		var way *way.Way
		var pathNodes []*crossing.Crossing
		shortestLength := math.Inf(1)

		if len(ways) > 1 {
			i.logger.Debug().Msgf("found %d ways between node %d and %d", len(ways), prevNode.OsmID, n.OsmID)
		}

		for _, w := range ways {
			cPathNodes, err := i.crossingRepository.SelectCrossingsFromWayID(w.OsmID)
			if err != nil {
				return nil, 0.0, fmt.Errorf("error while selecting nodes from way: %s", err.Error())
			}

			cPathNodes = i.weightRepository.CutPathNodes(&crossing.Crossing{Node: *prevNode}, w, cPathNodes)

			dist := i.weightRepository.CalculateDistances(prevNode, w, cPathNodes, n)
			if dist < shortestLength {
				shortestLength = dist
				way = w
				pathNodes = cPathNodes
			}
		}

		lengthInMeters += shortestLength

		startIndex := -1
		endIndex := -1
		for i, node := range pathNodes {
			if node.OsmID == prevNode.OsmID {
				startIndex = i
			}

			if node.OsmID == n.OsmID {
				endIndex = i
			}
		}

		if startIndex == -1 || endIndex == -1 {
			return nil, 0.0, fmt.Errorf("node %d or %d not found in way %d", prevNode.OsmID, n.OsmID, way.OsmID)
		}

		if startIndex > endIndex {
			pathNodes = arrayutil.Reverse(pathNodes[endIndex:(startIndex + 1)])
		} else {
			pathNodes = pathNodes[startIndex:(endIndex + 1)]
		}

		for _, node := range pathNodes {
			points = append(points, geojson.NewPoint(node.Lon, node.Lat))
		}

		prevNode = n
	}

	return points, lengthInMeters, nil
}

func (i *impl) GetNearestNode(lat float64, lon float64) (*node.Node, error) {
	nodes, err := i.nodeRepository.SelectNearNodesApprox(lat, lon, nearNodesApproxDistance)
	if err != nil {
		return nil, fmt.Errorf("error while selecting near nodes: %s", err.Error())
	}

	i.logger.Debug().Msgf("found %d near nodes", len(nodes))

	searchPoint := sphericmath.NewPoint(lat, lon)
	var nearestNode *node.Node
	var nearestNodeDistance float64

	var skippedNodes []int64

	for _, node := range nodes {
		if !i.hasEdges(node.OsmID) {
			skippedNodes = append(skippedNodes, node.OsmID)
			continue
		}

		nodePoint := sphericmath.NewPoint(node.Lat, node.Lon)
		distance := sphericmath.CalcDistanceInMeters(searchPoint, nodePoint)
		if nearestNode == nil || distance < nearestNodeDistance {
			nearestNode = node
			nearestNodeDistance = distance
		}
	}

	if nearestNode == nil {
		return nil, fmt.Errorf("no near node found (in %f meters)", math.Tan(nearNodesApproxDistance*math.Pi/180)*sphericmath.EarthRadius)
	}

	i.logger.WithAttrs("skipped", skippedNodes).Debug().Msgf("skipped %d nodes without edges", len(skippedNodes))

	return nearestNode, nil
}

func (i *impl) hasEdges(id int64) bool {
	ways, err := i.wayRepository.SelectWaysFromNode(id)
	if err != nil {
		i.logger.Error().Msgf("error while selecting ways from node: %s", err.Error())
		return false
	}

	for _, w := range ways {
		if !i.weightRepository.IsWayAllowed(*w, defaultVehicleType) {
			continue
		}

		crossings, err := i.crossingRepository.SelectCrossingsFromWayID(w.OsmID)
		if err != nil {
			i.logger.Error().Msgf("error while selecting nodes from way: %s", err.Error())
			continue
		}

		for _, c := range crossings {
			if c.IsCrossing {
				return true
			}
		}
	}

	return false
}
