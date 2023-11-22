package coordinatelist

import "github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinate"

type CoordinateList interface {
	Append(coordinate coordinate.Coordinate)
	Len() uint64
	Cap() uint64
	Get(index uint64) coordinate.Coordinate
	ToCoordinateArray() []coordinate.Coordinate
}

type CoordinateListImpl struct {
	lats []float64
	lons []float64
}

func NewCoordinateList(cap int) CoordinateList {
	return &CoordinateListImpl{
		lats: make([]float64, 0, cap),
		lons: make([]float64, 0, cap),
	}
}

func (c *CoordinateListImpl) Append(coordinate coordinate.Coordinate) {
	c.lats = append(c.lats, coordinate.Lat())
	c.lons = append(c.lons, coordinate.Lon())
}

func (c *CoordinateListImpl) Len() uint64 {
	return uint64(len(c.lats))
}

func (c *CoordinateListImpl) Cap() uint64 {
	return uint64(cap(c.lats))
}

func (c *CoordinateListImpl) Get(index uint64) coordinate.Coordinate {
	return coordinate.New(c.lons[index], c.lats[index])
}

func (c *CoordinateListImpl) ToCoordinateArray() []coordinate.Coordinate {
	cooArr := make([]coordinate.Coordinate, len(c.lats))
	for i := 0; i < len(c.lats); i++ {
		cooArr[i] = coordinate.New(c.lons[i], c.lats[i])
	}

	return cooArr
}
