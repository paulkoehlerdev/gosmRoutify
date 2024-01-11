package weightRepository

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/way"
	"math"
	"strconv"
)

type roadType int

const (
	unknown roadType = iota
	livingStreet
	cycleStreet
	urban
	rural
	ruralDual
	motorway
)

func (r roadType) String() string {
	switch r {
	case unknown:
		return "unknown"
	case livingStreet:
		return "livingStreet"
	case cycleStreet:
		return "cycleStreet"
	case urban:
		return "urban"
	case rural:
		return "rural"
	case ruralDual:
		return "ruralDual"
	case motorway:
		return "motorway"
	default:
		return "unknown"
	}
}

var typeToMaxSpeed = map[roadType]float64{
	unknown:      walkingSpeedBias,
	livingStreet: walkingSpeedBias,
	cycleStreet:  30,
	urban:        50,
	rural:        100,
	ruralDual:    130,
	motorway:     130,
}

func getRoadType(way way.Way) roadType {
	if isMotorway(way) {
		return motorway
	}

	if isRuralWithTwoLanes(way) {
		return ruralDual
	}

	if isRuralDualCarriageWay(way) {
		return ruralDual
	}

	if isRuralStreet(way) {
		return rural
	}

	if isUrbanStreet(way) {
		return urban
	}

	if isCycleStreet(way) {
		return cycleStreet
	}

	if isLivingStreet(way) {
		return livingStreet
	}

	if isUrbanStreetFuzzy(way) {
		return urban
	}

	return unknown
}

func calcMaxSpeedFromRoadType(way way.Way) float64 {
	return typeToMaxSpeed[getRoadType(way)]
}

func isLivingStreet(way way.Way) bool {
	// highway=living_street or living_street=yes
	if ls, ok := way.Tags["living_street"]; ok && (ls == "yes" || ls == "true" || ls == "1") {
		return true
	}

	if ls, ok := way.Tags["highway"]; ok && ls == "living_street" {
		return true
	}

	return false
}

func isCycleStreet(way way.Way) bool {
	// bicycle_road=yes or cyclestreet=yes
	if cs, ok := way.Tags["cyclestreet"]; ok && (cs == "yes" || cs == "true" || cs == "1") {
		return true
	}

	if br, ok := way.Tags["bicycle_road"]; ok && (br == "yes" || br == "true" || br == "1") {
		return true
	}

	return false
}

func isUrbanStreet(way way.Way) bool {
	// source:maxspeed~.*urban or maxspeed:type~.*urban or zone:maxspeed~.*urban or zone:traffic~.*urban or maxspeed~.*urban or HFCS~.*Urban.* or rural=no,
	// the maxspeed filters are ignored here, as they are processed in the maxspeed filter

	if m, ok := way.Tags["maxspeed"]; ok {
		maxspeed, err := strconv.ParseFloat(m, 64)
		if err == nil {
			if math.Abs(maxspeed-50) < 0.000001 {
				return true
			}
		}
	}

	if z, ok := way.Tags["maxspeed:type"]; ok && (z == "urban" || z == "DE:urban") {
		return true
	}

	if r, ok := way.Tags["rural"]; ok && r == "no" {
		return true
	}

	return false
}

func isUrbanStreetFuzzy(way way.Way) bool {
	// highway~living_street|residential or lit=yes or {sidewalk~yes|both|left|right|separate or sidewalk:left~yes|separate or sidewalk:right~yes|separate or sidewalk:both~yes|separate}

	if ls, ok := way.Tags["highway"]; ok && (ls == "living_street" || ls == "residential") {
		return true
	}

	if l, ok := way.Tags["lit"]; ok && (l == "yes" || l == "true" || l == "1") {
		return true
	}

	if s, ok := way.Tags["sidewalk"]; ok && (s == "yes" || s == "true" || s == "1" || s == "left" || s == "right" || s == "both" || s == "separate") {
		return true
	}

	if s, ok := way.Tags["sidewalk:left"]; ok && (s == "yes" || s == "true" || s == "1" || s == "separate") {
		return true
	}

	if s, ok := way.Tags["sidewalk:right"]; ok && (s == "yes" || s == "true" || s == "1" || s == "separate") {
		return true
	}

	return false
}

func isRuralStreet(way way.Way) bool {
	// source:maxspeed~.*rural or maxspeed:type~.*rural or zone:maxspeed~.*rural or zone:traffic~.*rural or maxspeed~.*rural or HFCS~.*Rural.* or rural=yes
	// maxspeed filters are ignored here, as they are processed in the maxspeed filter

	if m, ok := way.Tags["maxspeed"]; ok {
		maxspeed, err := strconv.ParseFloat(m, 64)
		if err == nil {
			if maxspeed > 50 {
				return true
			}
		}
	}

	if z, ok := way.Tags["maxspeed:type"]; ok && (z == "rural" || z == "DE:rural") {
		return true
	}

	if r, ok := way.Tags["rural"]; ok && r == "yes" {
		return true
	}

	return false
}

func isRuralDualCarriageWay(way way.Way) bool {
	// dual_carriageway=yes or maxspeed:type~\".*nsl_dual\""
	// maxspeed filters are ignored here, as they are processed in the maxspeed filter

	dualCarriage := false
	if dc, ok := way.Tags["dual_carriageway"]; ok && (dc == "yes" || dc == "true" || dc == "1") {
		dualCarriage = true
	}

	// {rural} and {dual carriageway}
	return isRuralStreet(way) && dualCarriage
}

func isRuralWithTwoLanes(way way.Way) bool {
	// oneway~yes|-1 or junction~roundabout|circular
	oneway := false
	if o, ok := way.Tags["oneway"]; ok && (o == "yes" || o == "true" || o == "1" || o == "-1") {
		oneway = true
	}

	if j, ok := way.Tags["junction"]; ok && (j == "roundabout" || j == "circular") {
		oneway = true
	}

	// ((!{oneway} and lanes>=4) or ({oneway} and lanes>=2))
	hasTwoOrMoreLanes := false
	if l, ok := way.Tags["lanes"]; ok {
		lanes, err := strconv.ParseFloat(l, 64)
		if err != nil {
			if !oneway && lanes >= 4 {
				hasTwoOrMoreLanes = true
			}

			if oneway && lanes >= 2 {
				hasTwoOrMoreLanes = true
			}
		}
	}

	// {rural} and {road with 2 or more lanes in each direction}
	return isRuralStreet(way) && hasTwoOrMoreLanes
}

func isMotorway(way way.Way) bool {
	// highway~motorway|motorway_link
	if h, ok := way.Tags["highway"]; ok && (h == "motorway" || h == "motorway_link") {
		return true
	}

	return false
}
