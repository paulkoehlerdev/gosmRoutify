package highwaytype

type Type int

const (
	Motorway Type = iota
	Link
	Trunk
	Primary
	Secondary
	Tertiary
	Residential
	Service
	LivingStreet
	Road
	Unclassified
	Unknown
)

var fClassToType = map[string]Type{
	"motorway": Motorway, "trunk": Trunk, "primary": Primary, "secondary": Secondary, "tertiary": Tertiary, "residential": Residential,
	"motorway_link": Link, "trunk_link": Link, "primary_link": Link, "secondary_link": Link, "tertiary_link": Link,
	"living_street": LivingStreet, "road": Road,
	// "service": Service, "unclassified": Unclassified,
}

var typeToMaxSpeed = map[Type]float64{
	Motorway:     180,
	Trunk:        100,
	Primary:      100,
	Secondary:    100,
	Tertiary:     50,
	Residential:  30,
	Link:         30,
	Service:      10,
	LivingStreet: 10,
	Road:         50,
	Unclassified: 50,
	Unknown:      1,
}

func FClassToSteetType(fClass string) Type {
	out, ok := fClassToType[fClass]
	if !ok {
		return Unknown
	}
	return out
}

func (t Type) DefaultMaxSpeed() float64 {
	return typeToMaxSpeed[t] / 3.6
}
