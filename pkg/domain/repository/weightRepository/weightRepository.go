package weightRepository

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/node"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/way"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/arrayutil"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geodistance"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"math"
)

const (
	maximumSpeedBias = 180.0
	walkingSpeedBias = 5
	minimumSpeedBias = 2.5

	laneFactor      = 0.5849 // 2^0.5849 ~ 1.5
	maximumLaneBias = 3.0
)

type WeightRepository interface {
	IsWayAllowed(way way.Way, vehicleType VehicleType) bool
	MaximumWayFactor(vehicleType VehicleType) float64
	CalculateWeights(from *node.Node, over *way.Way, to []*node.Node, vehicleType VehicleType) map[int64]float64
}

type impl struct {
	logger logging.Logger
}

func New(logger logging.Logger) WeightRepository {
	return &impl{
		logger: logger,
	}
}

func (i *impl) IsWayAllowed(way way.Way, vehicleType VehicleType) bool {
	return vehicleType.isWayTypeAllowed(way)
}

func (i *impl) MaximumWayFactor(vehicleType VehicleType) float64 {
	return vehicleType.maxmimumWayFactor()
}

func (i *impl) CalculateWeights(from *node.Node, over *way.Way, to []*node.Node, vehicleType VehicleType) map[int64]float64 {
	if from == nil {
		i.logger.Error().Msg("from node is nil")
		return make(map[int64]float64)
	}

	if over == nil {
		i.logger.Error().Msg("over way is nil")
		return make(map[int64]float64)
	}

	if oneway, ok := over.Tags["oneway"]; ok && !(oneway == "no" || oneway == "false" || oneway == "0") {
		fromIndex := -1
		for i, n := range to {
			if n.OsmID == from.OsmID {
				fromIndex = i
				break
			}
		}

		if fromIndex == -1 {
			i.logger.Error().Msg("from node not found in to nodes")
			return make(map[int64]float64)
		}

		if oneway == "yes" || oneway == "true" || oneway == "1" {
			to = to[fromIndex:]
		}

		if oneway == "-1" || oneway == "reverse" {
			to = to[:fromIndex+1]
		}
	}

	distances := i.calculateDistances(*from, to)
	out := make(map[int64]float64, len(to))

	for iter, node := range to {
		out[node.OsmID] = distances[iter] * vehicleType.calcWayFactor(*over)
	}

	return out
}

func (i *impl) calculateDistances(from node.Node, to []*node.Node) []float64 {
	out := make([]float64, len(to))

	nodeIndex := findNodeIndex(from, to)
	if nodeIndex == -1 {
		i.logger.Error().Msgf("node %d not found in to nodes", from.OsmID)
		for i := range out {
			out[i] = math.Inf(1)
		}
		return out
	}

	// forward
	currentDistance := 0.0
	previousNode := from

	for i, cNode := range to[nodeIndex+1:] {
		if cNode == nil {
			for i := range out {
				out[i] = math.Inf(1)
			}
			continue
		}

		additionalDistance := geodistance.CalcDistanceInMeters(
			geodistance.NewPoint(previousNode.Lat, previousNode.Lon),
			geodistance.NewPoint(cNode.Lat, cNode.Lon),
		)
		currentDistance += additionalDistance

		out[i] = currentDistance
		previousNode = *cNode
	}

	// backward
	currentDistance = 0.0
	previousNode = from

	for i, cNode := range arrayutil.Reverse(to[:nodeIndex]) {
		if cNode == nil {
			for i := range out {
				out[i] = math.Inf(1)
			}
			continue
		}

		additionalDistance := geodistance.CalcDistanceInMeters(
			geodistance.NewPoint(previousNode.Lat, previousNode.Lon),
			geodistance.NewPoint(cNode.Lat, cNode.Lon),
		)
		currentDistance += additionalDistance

		out[nodeIndex-i] = currentDistance
		previousNode = *cNode
	}

	return out
}

func findNodeIndex(node node.Node, nodes []*node.Node) int {
	for i, n := range nodes {
		if n.OsmID == node.OsmID {
			return i
		}
	}
	return -1
}
