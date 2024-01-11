package sphericmath

import (
	"math"
)

const (
	EarthRadius = 6378137.0
)

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

func (p Point) toLatLngRad() (float64, float64) {
	return p.Lat() * math.Pi / 180, p.Lon() * math.Pi / 180
}

func CalcDistanceInMeters(a, b Point) float64 {
	alat, alon := a.toLatLngRad()
	blat, blon := b.toLatLngRad()

	return RadToMeters(haversine(alat, alon, blat, blon))
}

func haversine(phiA, lambdaA, phiB, lambdaB float64) float64 {
	dPhi := phiA - phiB
	dLambda := lambdaA - lambdaB

	x1 := math.Pow(math.Sin(dPhi/2), 2)
	x2 := math.Cos(phiA) * math.Cos(lambdaA)
	x3 := math.Pow(math.Sin(dLambda/2), 2)

	x := x1 + x2*x3
	return 2 * math.Atan2(math.Sqrt(x), math.Sqrt(1-x))
}

func RadToMeters(rad float64) float64 {
	return rad * EarthRadius
}

func CalculateBearing(a, b Point) float64 {
	// Kurswinkelberechnung
	phiA, lambdaA := a.toLatLngRad()
	phiB, lambdaB := b.toLatLngRad()

	dLambda := lambdaB - lambdaA

	y := math.Sin(dLambda) * math.Cos(phiB)
	x := math.Cos(phiA)*math.Sin(phiB) - math.Sin(phiA)*math.Cos(phiB)*math.Cos(dLambda)

	return math.Atan2(y, x)
}
