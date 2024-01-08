package geodistance

import (
	"math"
)

const (
	EarthRadius = 6378137.0

	wgs84A = 6377397.1550
	wgs84B = 6356078.9629

	// wgs84A  = 6378137.000
	// wgs84B  = 6356752.31424518
	wgs84E2 = (wgs84A*wgs84A - wgs84B*wgs84B) / (wgs84A * wgs84A)
)

type Coordinate interface {
	Lat() float64
	Lon() float64
}

type Point [2]float64

func NewPoint(lat float64, lon float64) Point {
	return Point{lat, lon}
}

func (p Point) Lat() float64 {
	return p[0]
}

func (p Point) Lon() float64 {
	return p[1]
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

	return RadToMeters(2 * math.Asin(math.Sqrt(x1+x2*x3)))
}

func RadToMeters(rad float64) float64 {
	return rad * EarthRadius
}

/*
	x, y, z := convertToLocalCoordinates(
		a.Lat()*math.Pi/180,
		a.Lon()*math.Pi/180,
		b.Lat()*math.Pi/180,
		b.Lon()*math.Pi/180,
	)
	return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2) + math.Pow(z, 2))
}

func convertToLocalCoordinates(centerPhi, centerLambda, phi, lambda float64) (float64, float64, float64) {
	centerH := 100.0

	N := wgs84A / (math.Sqrt(1 - (wgs84E2 * math.Pow(math.Sin(centerPhi), 2))))
	centerX := (N + centerH) * math.Cos(centerPhi) * math.Cos(centerLambda)
	centerY := (N + centerH) * math.Cos(centerPhi) * math.Sin(centerLambda)
	centerZ := ((wgs84B*wgs84B)/(wgs84A*wgs84A)*N + centerH) * math.Sin(centerPhi)

	x := -math.Sin(phi)*math.Cos(lambda)*centerX - math.Sin(phi)*math.Sin(lambda)*centerY + math.Cos(phi)*centerZ
	y := -math.Sin(lambda)*centerX + math.Cos(lambda)*centerY
	z := math.Cos(phi)*math.Cos(lambda)*centerX + math.Cos(phi)*math.Sin(lambda)*centerY + math.Sin(phi)*centerZ

	x = centerX - x
	y = centerY - y
	z = centerZ - z

	return x, y, z
}

*/
