package weightRepository

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/node"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/way"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/sphericmath"
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
	Car:        160,
	Bike:       30,
	Pedestrian: walkingSpeedBias,
}

func (v VehicleType) isWayTypeAllowed(way way.Way) bool {
	if way.OsmID == 293556313 {
		println("wayType", getRoadType(way))
	}

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
	return 1 / (maxVehicleTypeSpeed[v] / 3.6)
}

func (v VehicleType) calcWayFactor(way way.Way) float64 {
	maxWaySpeed := calcMaxWaySpeed(way)

	if way.OsmID == 159692417 {
		println("maxWaySpeed", maxWaySpeed, "wayType", getRoadType(way))
	}

	if way.OsmID == 389295370 {
		println("maxWaySpeed", maxWaySpeed, "wayType", getRoadType(way))
	}

	maxWaySpeed = math.Min(maxWaySpeed, maxVehicleTypeSpeed[v])

	return 1 / (maxWaySpeed / 3.6)
}

const (
	tenDegree = 10 * math.Pi / 180
)

func (v VehicleType) calcCrossingFactor(prev, curr, next *node.Node) float64 {
	if prev == nil || curr == nil || next == nil {
		return 0
	}

	if curr.OsmID == 16526339 {
		println("curr", curr.OsmID)
	}

	if v != Car {
		return 0
	}

	phi1 := sphericmath.CalculateBearing(
		sphericmath.NewPoint(curr.Lat, curr.Lon),
		sphericmath.NewPoint(prev.Lat, prev.Lon),
	)

	phi2 := sphericmath.CalculateBearing(
		sphericmath.NewPoint(curr.Lat, curr.Lon),
		sphericmath.NewPoint(next.Lat, next.Lon),
	)

	if math.IsNaN(phi1) || math.IsNaN(phi2) {
		return 0
	}

	phi := phi2 - phi1 + math.Pi
	phi = math.Mod(phi+math.Pi, 2*math.Pi) - math.Pi

	// straight
	if math.Abs(phi) < (math.Pi / 2) {
		return 0
	}

	// uturn
	if math.Abs(phi) > (math.Pi - tenDegree) {
		return 15
	}

	// left
	if phi < 0 {
		return 10
	}

	// right
	return 10
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
