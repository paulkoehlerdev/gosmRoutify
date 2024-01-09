package weightRepository

type highwayClass int

const (
	motorway highwayClass = iota
	link
	trunk
	primary
	secondary
	tertiary
	residential
	service
	livingStreet
	road
	unclassified
	unknown
)

var fClassToType = map[string]highwayClass{
	"motorway":       motorway,
	"trunk":          trunk,
	"primary":        primary,
	"secondary":      secondary,
	"tertiary":       tertiary,
	"residential":    residential,
	"motorway_link":  link,
	"trunk_link":     link,
	"primary_link":   link,
	"secondary_link": link,
	"tertiary_link":  link,
	"living_street":  livingStreet,
	"road":           road,
	"service":        service,
	"unclassified":   unclassified,
}

var typeToMaxSpeed = map[highwayClass]float64{
	motorway:     180,
	trunk:        100,
	primary:      100,
	secondary:    100,
	tertiary:     50,
	residential:  30,
	link:         30,
	service:      10,
	livingStreet: 10,
	road:         50,
	unclassified: 50,
	unknown:      minimumSpeedBias,
}

func fClassToSteetType(fClass string) highwayClass {
	out, ok := fClassToType[fClass]
	if !ok {
		return unknown
	}
	return out
}

func (t highwayClass) DefaultMaxSpeed() float64 {
	return typeToMaxSpeed[t]
}
