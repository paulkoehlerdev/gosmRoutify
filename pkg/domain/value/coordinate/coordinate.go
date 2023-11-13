package coordinate

type Coordinate [2]float64

func New(lat float64, lon float64) Coordinate {
	return Coordinate{lat, lon}
}

func (c Coordinate) Lat() float64 {
	return c[0]
}

func (c Coordinate) Lon() float64 {
	return c[1]
}
