package router

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/address"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/node"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/addressService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/graphService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/nodeService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/astar"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geojson"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"time"
)

type RouteSegmentInfo struct {
	LengthInMeters float64         `json:"distance"`
	LengthInTime   int64           `json:"time"`
	GeoJson        geojson.GeoJson `json:"geojson"`
}

type Application interface {
	FindRoute(points []geojson.Point) ([]RouteSegmentInfo, error)
	FindAddresses(query string) ([]*address.Address, error)
	LocateAddressByID(id int64) (geojson.Point, error)
}

const (
	maxVisitedNodes = 500000
)

type impl struct {
	logger         logging.Logger
	graphService   graphService.GraphService
	addressService addressService.AddressService
	nodeService    nodeService.NodeService
}

func New(graphService graphService.GraphService, addressService addressService.AddressService, nodeService nodeService.NodeService, logger logging.Logger) Application {
	return &impl{
		logger:         logger,
		graphService:   graphService,
		addressService: addressService,
		nodeService:    nodeService,
	}
}

func (i *impl) FindRoute(points []geojson.Point) ([]RouteSegmentInfo, error) {
	startTime := time.Now()

	out := make([]RouteSegmentInfo, 0, len(points)-1)

	nodes := make([]*node.Node, len(points))
	for index, point := range points {
		node, err := i.graphService.GetNearestNode(point.Lon(), point.Lat())
		if err != nil {
			return nil, fmt.Errorf("error while finding nearest node to [%f, %f]: %s", point.Lat(), point.Lon(), err.Error())
		}

		nodes[index] = node
	}

	i.logger.Debug().Msgf("calculated nearest node in %s", time.Since(startTime).String())

	start := nodes[0]
	for index, end := range nodes[1:] {
		path, length, err := astar.AStar[int64, float64](start.OsmID, end.OsmID, i.graphService.GetEdges(*end), i.graphService.GetHeuristic(*end), maxVisitedNodes)
		if err != nil {
			return nil, fmt.Errorf("error while routing: %s", err.Error())
		}

		nodePoints, lengthInMeters, err := i.graphService.CalculatePathInformation(path)
		if err != nil {
			return nil, fmt.Errorf("error while building geojson line: %s", err.Error())
		}

		nodePoints = append(
			[]geojson.Point{points[index]},
			nodePoints...,
		)

		nodePoints = append(
			nodePoints,
			points[index+1],
		)

		geometry := geojson.LineString(nodePoints).ToGeometry()

		geoJson := geojson.NewEmptyGeoJson()
		geoJson.AddFeature(geojson.Feature{
			Type:       "Feature",
			Geometry:   geometry,
			Properties: nil,
		})

		out = append(out, RouteSegmentInfo{
			LengthInMeters: lengthInMeters,
			LengthInTime:   int64(length),
			GeoJson:        geoJson,
		})

		start = end
	}

	return out, nil
}

func (i *impl) FindAddresses(query string) ([]*address.Address, error) {
	return i.addressService.GetSearchResultsFromAddress(query)
}

func (i *impl) LocateAddressByID(id int64) (geojson.Point, error) {
	lat, lon, err := i.nodeService.LocateOsmID(id)
	if err != nil {
		return geojson.Point{}, fmt.Errorf("error while locating address: %s", err.Error())
	}

	return geojson.NewPoint(lon, lat), nil
}
