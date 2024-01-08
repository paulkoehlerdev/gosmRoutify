package graphService

import (
	"encoding/json"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/node"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/nodeRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/wayRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geodistance"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geojson"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"math"
)

const nearNodesApproxDistance = 0.001 // approx. 1km

type GraphService interface {
	GetEdges(id int64) map[int64]float64
	GetNearestNode(lat float64, lon float64) (*node.Node, error)
}

type impl struct {
	nodeRepository nodeRepository.NodeRepository
	wayRepository  wayRepository.WayRepository
	logger         logging.Logger
}

func New(nodeRepository nodeRepository.NodeRepository, wayRepository wayRepository.WayRepository, logger logging.Logger) GraphService {
	return &impl{
		nodeRepository: nodeRepository,
		wayRepository:  wayRepository,
		logger:         logger,
	}
}

func (i *impl) GetEdges(id int64) map[int64]float64 {
	return make(map[int64]float64)
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
