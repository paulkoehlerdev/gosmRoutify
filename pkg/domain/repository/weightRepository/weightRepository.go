package weightRepository

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/crossing"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/node"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/way"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/sphericmath"
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
	CalculateWeights(prevNode *node.Node, from *crossing.Crossing, over *way.Way, to []*crossing.Crossing, end node.Node, vehicleType VehicleType) map[int64]float64
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

func (i *impl) CalculateWeights(prevNode *node.Node, from *crossing.Crossing, over *way.Way, to []*crossing.Crossing, end node.Node, vehicleType VehicleType) map[int64]float64 {
	if from == nil {
		i.logger.Error().Msg("from node is nil")
		return make(map[int64]float64)
	}

	if over == nil {
		i.logger.Error().Msg("over way is nil")
		return make(map[int64]float64)
	}

	to = i.cutOneway(*from, *over, to)
	if to == nil {
		i.logger.Error().Msg("to nodes are nil after cutting oneway")
		return make(map[int64]float64)
	}

	to = i.cutCrossing(*from, to)
	if to == nil {
		i.logger.Error().Msg("to nodes are nil after cutting crossing")
		return make(map[int64]float64)
	}

	distancesToCrossings := i.calculateDistances(*from, to, end)

	out := make(map[int64]float64)
	for crossing, length := range distancesToCrossings {
		out[crossing.OsmID] = length*
			vehicleType.calcWayFactor(*over) +
			vehicleType.calcCrossingFactor(prevNode, &from.Node, &crossing.Node)
	}

	return out
}

func (i *impl) cutOneway(from crossing.Crossing, over way.Way, to []*crossing.Crossing) []*crossing.Crossing {
	fromIndex := -1
	for i, n := range to {
		if n.OsmID == from.OsmID {
			fromIndex = i
			break
		}
	}

	if fromIndex == -1 {
		return nil
	}

	if oneway, ok := over.Tags["oneway"]; ok && !(oneway == "no" || oneway == "false" || oneway == "0") {
		if oneway == "yes" || oneway == "true" || oneway == "1" {
			return to[fromIndex:]
		}

		if oneway == "-1" || oneway == "reverse" {
			return to[:fromIndex+1]
		}
	}

	if j, ok := over.Tags["junction"]; ok && (j == "roundabout" || j == "circular") {
		return to[fromIndex:]
	}

	return to
}

func (i *impl) cutCrossing(from crossing.Crossing, to []*crossing.Crossing) []*crossing.Crossing {
	cutFrom := -1
	cutTo := math.MaxInt64

	for i, n := range to {
		if n.OsmID == from.OsmID {
			if cutFrom == -1 {
				cutFrom = i
			}
			cutTo = i
			continue
		}

		if cutTo < i {
			if n.IsCrossing {
				cutTo = i
				break
			}
			continue
		}

		if n.IsCrossing {
			cutFrom = i
		}
	}

	if cutFrom == -1 || cutTo == math.MaxInt64 {
		return nil
	}

	return to[cutFrom : cutTo+1]
}

func (i *impl) calculateDistances(from crossing.Crossing, to []*crossing.Crossing, end node.Node) map[*crossing.Crossing]float64 {
	out := make(map[*crossing.Crossing]float64)

	fullLength := 0.0
	leftLength := 0.0
	endLength := math.NaN()

	prevNode := to[0]

	for _, n := range to[1:] {
		dist := sphericmath.CalcDistanceInMeters(
			sphericmath.NewPoint(prevNode.Lat, prevNode.Lon),
			sphericmath.NewPoint(n.Lat, n.Lon),
		)

		fullLength += dist

		if n.OsmID == end.OsmID {
			endLength = fullLength
		}

		if from.OsmID == n.OsmID {
			leftLength = fullLength
		}

		prevNode = n
	}

	if to[0].OsmID != from.OsmID {
		out[to[0]] = leftLength
	}

	if to[len(to)-1].OsmID != from.OsmID {
		out[to[len(to)-1]] = fullLength - leftLength
	}

	if !math.IsNaN(endLength) {
		out[&crossing.Crossing{Node: end}] = math.Abs(leftLength - endLength)
	}

	return out
}
