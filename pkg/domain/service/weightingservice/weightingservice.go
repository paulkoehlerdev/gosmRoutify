package weightingservice

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/graphtile"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/highwaytype"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geodistance"
	"math"
	"strconv"
)

type WeightingService interface {
	Weight(info graphtile.EdgeInfo) float64
}

type impl struct {
}

func New() WeightingService {
	return &impl{}
}

func (i *impl) Weight(info graphtile.EdgeInfo) float64 {
	length := 0.0
	prev := info.Nodes[0]
	for _, node := range info.Nodes[1:] {
		length += geodistance.CalcDistanceInMeters(prev, node)
		prev = node
	}

	highwayClass, ok := info.Tags["highway"]
	if !ok {
		highwayClass = "unclassified"
	}

	speed := 0.0
	speedStr, ok := info.Tags["maxspeed"]
	if !ok {
		speed = getDefaultSpeed(highwayClass)
	} else {
		speed = parseSpeed(speedStr)
		if speed == 0.0 || math.IsNaN(speed) {
			speed = getDefaultSpeed(highwayClass)
		}
	}

	if speed == 0.0 {
		speed = 1
	}

	return length / speed
}

func parseSpeed(speedStr string) float64 {
	speed, err := strconv.ParseFloat(speedStr, 64)
	if err != nil {
		return 0.0
	}
	return speed
}

func getDefaultSpeed(highwayClass string) float64 {
	return highwaytype.FClassToSteetType(highwayClass).DefaultMaxSpeed()
}
