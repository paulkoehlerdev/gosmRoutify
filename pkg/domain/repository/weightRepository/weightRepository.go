package weightRepository

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/node"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/way"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/arrayutil"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geodistance"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"math"
	"strconv"
	"strings"
)

const (
	maximumSpeedBias = 180.0
	minimumSpeedBias = 5.0

	walkinSpeedBias = 15.0
)

type WeightRepository interface {
	CalculateWeights(from *node.Node, over *way.Way, to []*node.Node) map[int64]float64
}

type impl struct {
	logger logging.Logger
}

func New(logger logging.Logger) WeightRepository {
	return &impl{
		logger: logger,
	}
}

func (i *impl) CalculateWeights(from *node.Node, over *way.Way, to []*node.Node) map[int64]float64 {
	if from == nil {
		i.logger.Error().Msg("from node is nil")
		return make(map[int64]float64)
	}

	if over == nil {
		i.logger.Error().Msg("over way is nil")
		return make(map[int64]float64)
	}

	distances := i.calculateDistances(*from, to)
	out := make(map[int64]float64, len(to))

	wayFactor := 1.0

	if v, ok := over.Tags["maxspeed"]; ok && v != "" {
		wayFactor /= calcMaxSpeed(over.Tags["maxspeed"]) * 0.9
	} else if v, ok := over.Tags["highway"]; ok && v != "" {
		wayFactor /= calcMaxSpeedFromRoadType(over.Tags["highway"])
	}

	for i, node := range to {
		out[node.OsmID] = distances[i] * wayFactor
	}

	return out
}

func calcMaxSpeed(maxSpeed string) float64 {
	maxSpeed = strings.SplitN(maxSpeed, " ", 2)[0]

	if speed, err := strconv.ParseFloat(maxSpeed, 64); err == nil {
		return speed
	}

	if maxSpeed == "walk" {
		return walkinSpeedBias
	}

	if maxSpeed == "none" {
		return maximumSpeedBias
	}

	return minimumSpeedBias
}

func calcMaxSpeedFromRoadType(highwayClass string) float64 {
	return fClassToSteetType(highwayClass).DefaultMaxSpeed()
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
