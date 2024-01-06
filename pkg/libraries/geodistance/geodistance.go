package geodistance

import (
	"math"
)

const earthRadius = 6378137.0

type Coordinate interface {
	Lat() float64
	Lon() float64
}

func CalcDistanceInMeters(a, b Coordinate) float64 {
	alat := a.Lat() * math.Pi / 180
	alon := a.Lon() * math.Pi / 180

	blat := b.Lat() * math.Pi / 180
	blon := b.Lon() * math.Pi / 180

	dlat := alat - blat
	dlon := alon - blon

	x1 := math.Pow(math.Sin(dlat/2), 2)
	x2 := 1 - x1 - math.Pow(math.Sin((alat+blat)/2), 2)
	x3 := math.Pow(math.Sin(dlon/2), 2)

	return 2 * math.Asin(math.Sqrt(x1+x2*x3)) * earthRadius
}
