package geoJson

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
	return Geometry{
		Type:        "GeometryCollection",
		Coordinates: []interface{}{g},
	}
}

type MultiPolygon []Polygon

func (m MultiPolygon) ToGeometry() Geometry {
	return Geometry{
		Type:        "MultiPolygon",
		Coordinates: []interface{}{m},
	}
}

type MultiLineString []LineString

func (m MultiLineString) ToGeometry() Geometry {
	return Geometry{
		Type:        "MultiLineString",
		Coordinates: []interface{}{m},
	}
}

type MultiPoint []Point

func (m MultiPoint) ToGeometry() Geometry {
	return Geometry{
		Type:        "MultiPoint",
		Coordinates: []interface{}{m},
	}
}

type Polygon [][]Point

func (p Polygon) ToGeometry() Geometry {
	return Geometry{
		Type:        "Polygon",
		Coordinates: []interface{}{p},
	}
}

type LineString []Point

func (l LineString) ToGeometry() Geometry {
	return Geometry{
		Type:        "LineString",
		Coordinates: []interface{}{l},
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
