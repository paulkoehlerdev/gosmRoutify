package kdtree

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinate"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geodistance"
	"math"
)

type KdTree[T any] interface {
	Insert(coo coordinate.Coordinate, data T)
	FindNearest(coo coordinate.Coordinate) (float64, *T)
}

type impl[T any] struct {
	Root *node[T]
}

func New[T any]() KdTree[T] {
	return &impl[T]{}
}

func (i *impl[T]) Insert(coo coordinate.Coordinate, data T) {
	if i.Root == nil {
		i.Root = &node[T]{
			Point: coo,
			Data:  data,
			Dim:   false,
		}
	}

	i.Root.insert(coo, data)
}

func (i *impl[T]) FindNearest(coo coordinate.Coordinate) (float64, *T) {
	if i.Root == nil {
		return math.NaN(), nil
	}

	dist, data := i.Root.findNearest(coo)

	return dist, &data
}

type node[T any] struct {
	Point coordinate.Coordinate
	Data  T
	Dim   bool
	A     *node[T]
	B     *node[T]
}

func (n *node[T]) compare(p coordinate.Coordinate) bool {
	comp := n.Point.Lat() < p.Lat()
	if n.Dim {
		comp = n.Point.Lon() < p.Lon()
	}
	return comp
}

func (n *node[T]) checkHyperplane(p coordinate.Coordinate, dist float64) bool {
	if n.Dim {
		return geodistance.CalcDistanceInMeters(
			coordinate.New(n.Point.Lat(), p.Lon()),
			p,
		) < dist
	}
	return geodistance.CalcDistanceInMeters(
		coordinate.New(p.Lat(), n.Point.Lon()),
		p,
	) < dist
}

func (n *node[T]) insert(newP coordinate.Coordinate, data T) {
	if n.compare(newP) {
		if n.A == nil {
			n.A = &node[T]{
				Point: newP,
				Dim:   !n.Dim,
				Data:  data,
			}
		} else {
			n.A.insert(newP, data)
		}
	} else {
		if n.B == nil {
			n.B = &node[T]{
				Point: newP,
				Dim:   !n.Dim,
				Data:  data,
			}
		} else {
			n.B.insert(newP, data)
		}
	}
}

func (n *node[T]) checkLeaf(p coordinate.Coordinate) (float64, T) {
	return geodistance.CalcDistanceInMeters(n.Point, p), n.Data
}

func (n *node[T]) checkTree(p coordinate.Coordinate, first *node[T], second *node[T]) (float64, T) {
	dist, data := first.findNearest(p)

	if d := geodistance.CalcDistanceInMeters(n.Point, p); d < dist {
		dist = d
		data = n.Data
	}

	if second != nil && n.checkHyperplane(p, dist) {
		tempDist, tempData := second.findNearest(p)
		if tempDist < dist {
			dist = tempDist
			data = tempData
		}
	}

	return dist, data
}

func (n *node[T]) findNearest(p coordinate.Coordinate) (float64, T) {
	if n.compare(p) {
		if n.A == nil {
			return n.checkLeaf(p)
		} else {
			return n.checkTree(p, n.A, n.B)
		}
	} else {
		if n.B == nil {
			return n.checkLeaf(p)
		} else {
			return n.checkTree(p, n.B, n.A)
		}
	}
}
