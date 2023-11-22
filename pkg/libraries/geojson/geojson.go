package geojson

type GeoJson struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

func (g *GeoJson) AddFeature(feature Feature) {
	g.Features = append(g.Features, feature)
}

func NewEmptyGeoJson() GeoJson {
	return GeoJson{
		Type:     "FeatureCollection",
		Features: []Feature{},
	}
}

type Feature struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

func NewFeature(geometry Geometry) Feature {
	return Feature{
		Type:       "Feature",
		Geometry:   geometry,
		Properties: make(Properties),
	}
}

type Geometry struct {
	Type        string        `json:"type"`
	Coordinates []interface{} `json:"coordinates"`
}

type Properties map[string]interface{}

type GeometryCollection []Geometry

func (g GeometryCollection) ToGeometry() Geometry {
	var out []interface{}
	for _, geometry := range g {
		out = append(out, geometry)
	}

	return Geometry{
		Type:        "GeometryCollection",
		Coordinates: out,
	}
}

type MultiPolygon []Polygon

func (m MultiPolygon) ToGeometry() Geometry {
	var out []interface{}
	for _, polygon := range m {
		out = append(out, polygon)
	}

	return Geometry{
		Type:        "MultiPolygon",
		Coordinates: out,
	}
}

type MultiLineString []LineString

func (m MultiLineString) ToGeometry() Geometry {
	var out []interface{}
	for _, lineString := range m {
		out = append(out, lineString)
	}

	return Geometry{
		Type:        "MultiLineString",
		Coordinates: out,
	}
}

type MultiPoint []Point

func (m MultiPoint) ToGeometry() Geometry {
	var out []interface{}
	for _, point := range m {
		out = append(out, point)
	}

	return Geometry{
		Type:        "MultiPoint",
		Coordinates: out,
	}
}

type Polygon [][]Point

func (p Polygon) ToGeometry() Geometry {
	var out []interface{}
	for _, lineString := range p {
		out = append(out, lineString)
	}

	return Geometry{
		Type:        "Polygon",
		Coordinates: out,
	}
}

type LineString []Point

func (l LineString) ToGeometry() Geometry {
	var out []interface{}
	for _, point := range l {
		out = append(out, point)
	}

	return Geometry{
		Type:        "LineString",
		Coordinates: out,
	}
}

type Point [2]float64

func NewPoint(lat float64, lon float64) Point {
	return Point{lon, lat}
}

func (p Point) Lat() float64 {
	return p[1]
}

func (p Point) Lon() float64 {
	return p[0]
}

func (p Point) ToGeometry() Geometry {
	return Geometry{
		Type:        "Point",
		Coordinates: []interface{}{p[0], p[1]},
	}
}
