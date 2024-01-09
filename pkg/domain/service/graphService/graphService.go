package graphService

import (
	"encoding/json"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/node"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/way"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/nodeRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/wayRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/weightRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/arrayutil"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geodistance"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geojson"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"math"
)

const nearNodesApproxDistance = 0.001 // approx. 1km

type GraphService interface {
	GetEdges(id int64) map[int64]float64
	GetHeuristic(end *node.Node) func(id int64) float64
	BuildGeojsonLineFromPath(path []int64) ([]geojson.Point, error)
	GetNearestNode(lat float64, lon float64) (*node.Node, error)
}

type impl struct {
	nodeRepository   nodeRepository.NodeRepository
	wayRepository    wayRepository.WayRepository
	weightRepository weightRepository.WeightRepository
	logger           logging.Logger
}

func New(nodeRepository nodeRepository.NodeRepository, wayRepository wayRepository.WayRepository, weightRepository weightRepository.WeightRepository, logger logging.Logger) GraphService {
	return &impl{
		nodeRepository:   nodeRepository,
		wayRepository:    wayRepository,
		weightRepository: weightRepository,
		logger:           logger,
	}
}

func (i *impl) GetEdges(id int64) map[int64]float64 {
	way, err := i.wayRepository.SelectWaysFromNode(id)
	if err != nil {
		i.logger.Error().Msgf("error while selecting ways from node: %s", err.Error())
		return make(map[int64]float64)
	}

	out := make(map[int64]float64)

	var fromNode *node.Node
	for _, w := range way {
		nodes, err := i.nodeRepository.SelectNodesFromWay(w.OsmID)
		if err != nil {
			i.logger.Error().Msgf("error while selecting nodes from way: %s", err.Error())
			continue
		}

		if fromNode == nil {
			for _, n := range nodes {
				if n.OsmID == id {
					fromNode = n
					break
				}
			}
		}

		weights := i.weightRepository.CalculateWeights(fromNode, w, nodes)
		for k, v := range weights {
			if prevV, ok := out[k]; ok && prevV < v {
				continue
			}
			out[k] = v
		}
	}

	return out
}

func (i *impl) GetHeuristic(end *node.Node) func(id int64) float64 {
	if end == nil {
		return func(id int64) float64 {
			return 0
		}
	}

	return func(nodeId int64) float64 {
		node, err := i.nodeRepository.SelectNodeFromID(nodeId)
		if err != nil {
			i.logger.Error().Msgf("error while selecting node from id: %s", err.Error())
			return 0
		}

		return geodistance.CalcDistanceInMeters(geodistance.NewPoint(end.Lat, end.Lon), geodistance.NewPoint(node.Lat, node.Lon))
	}
}

func (i *impl) BuildGeojsonLineFromPath(path []int64) ([]geojson.Point, error) {
	prevNode, err := i.nodeRepository.SelectNodeFromID(path[0])
	if err != nil {
		return nil, fmt.Errorf("error while selecting node from id: %s", err.Error())
	}

	var points []geojson.Point

	for _, nodeId := range path[1:] {
		n, err := i.nodeRepository.SelectNodeFromID(nodeId)
		if err != nil {
			return nil, fmt.Errorf("error while selecting node from id: %s", err.Error())
		}

		ways, err := i.wayRepository.SelectWaysFromTwoNodeIDs(prevNode.OsmID, n.OsmID)
		if err != nil {
			return nil, fmt.Errorf("error while selecting ways from two nodes: %s", err.Error())
		}

		var way *way.Way
		weight := math.Inf(1)
		for _, w := range ways {
			weights := i.weightRepository.CalculateWeights(prevNode, w, []*node.Node{prevNode, n})
			if weights[n.OsmID] < weight {
				weight = weights[n.OsmID]
				way = w
			}
		}

		if way == nil {
			return nil, fmt.Errorf("no way found between node %d and %d", prevNode.OsmID, n.OsmID)
		}

		pathNodes, err := i.nodeRepository.SelectNodesFromWay(way.OsmID)
		if err != nil {
			return nil, fmt.Errorf("error while selecting nodes from way: %s", err.Error())
		}

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
			return nil, fmt.Errorf("node %d or %d not found in way %d", prevNode.OsmID, n.OsmID, way.OsmID)
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

	return points, nil
}

func (i *impl) GetNearestNode(lat float64, lon float64) (*node.Node, error) {
	nodes, err := i.nodeRepository.SelectNearNodesApprox(lat, lon, nearNodesApproxDistance)
	if err != nil {
		return nil, fmt.Errorf("error while selecting near nodes: %s", err.Error())
	}

	geojsonObj := geojson.NewEmptyGeoJson()
	var multiLineString geojson.MultiLineString
	for _, node := range nodes {
		var lineString geojson.LineString
		lineString = append(lineString, geojson.NewPoint(lon, lat))
		lineString = append(lineString, geojson.NewPoint(node.Lon, node.Lat))
		multiLineString = append(multiLineString, lineString)
	}
	geojsonObj.AddFeature(geojson.NewFeature(multiLineString.ToGeometry()))

	geoJsonBytes, err := json.Marshal(geojsonObj)
	if err != nil {
		i.logger.Error().Msgf("error while marshalling geojson: %s", err.Error())
	}

	i.logger.Debug().WithAttrs("geojson", string(geoJsonBytes)).Msgf("found %d near nodes", len(nodes))

	searchPoint := geodistance.NewPoint(lat, lon)
	var nearestNode *node.Node
	var nearestNodeDistance float64

	for _, node := range nodes {
		nodePoint := geodistance.NewPoint(node.Lat, node.Lon)
		distance := geodistance.CalcDistanceInMeters(searchPoint, nodePoint)
		if nearestNode == nil || distance < nearestNodeDistance {
			nearestNode = node
			nearestNodeDistance = distance
		}
	}

	if nearestNode == nil {
		return nil, fmt.Errorf("no near node found (in %f meters)", math.Tan(nearNodesApproxDistance*math.Pi/180)*geodistance.EarthRadius)
	}

	return nearestNode, nil
}
