package weightRepository

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/way"
	"math"
	"strconv"
	"strings"
)

type VehicleType int

const (
	Car VehicleType = iota
	Bike
	Pedestrian
)

var maxVehicleTypeSpeed = map[VehicleType]float64{
	Car:        math.Inf(1),
	Bike:       30,
	Pedestrian: walkingSpeedBias,
}

func (v VehicleType) isWayTypeAllowed(way way.Way) bool {
	wayRoadType := getRoadType(way)
	switch v {
	case Car:
		return wayRoadType > unknown
	case Bike:
		return wayRoadType < ruralDual
	case Pedestrian:
		return wayRoadType < rural
	default:
		return false
	}
}

func (v VehicleType) maxmimumWayFactor() float64 {
	return 1 / (maxVehicleTypeSpeed[v] * 3.6)
}

func (v VehicleType) calcWayFactor(way way.Way) float64 {
	maxWaySpeed := calcMaxWaySpeed(way)

	maxWaySpeed = math.Min(maxWaySpeed, maxVehicleTypeSpeed[v])

	return 1 / (maxWaySpeed * 3.6)
}

func (v VehicleType) String() string {
	switch v {
	case Car:
		return "car"
	case Bike:
		return "bike"
	case Pedestrian:
		return "pedestrian"
	default:
		return "unknown"
	}
}

func calcMaxWaySpeed(way way.Way) float64 {
	speedBias := minimumSpeedBias
	if v, ok := way.Tags["maxspeed"]; ok && v != "" {
		speedBias = calcMaxSpeed(way.Tags["maxspeed"])
	} else if v, ok := way.Tags["highway"]; ok && v != "" {
		speedBias = calcMaxSpeedFromRoadType(way)
	}

	if speedBias > maximumSpeedBias {
		speedBias = maximumSpeedBias
	}

	return speedBias
}

func calcMaxSpeed(maxSpeed string) float64 {
	maxSpeed = strings.SplitN(maxSpeed, " ", 2)[0]

	if speed, err := strconv.ParseFloat(maxSpeed, 64); err == nil {
		return speed
	}

	if maxSpeed == "walk" {
		return walkingSpeedBias
	}

	if maxSpeed == "none" {
		return maximumSpeedBias
	}

	return minimumSpeedBias
}
