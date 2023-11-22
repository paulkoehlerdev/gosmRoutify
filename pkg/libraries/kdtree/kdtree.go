package kdtree

import "math"

type KdTree[T any] interface {
	Insert(x float64, y float64, data T)
	SearchNearest(x float64, y float64) *T
}

type impl[T any] struct {
	Root *node[T]
}

func New[T any]() KdTree[T] {
	return &impl[T]{}
}

func (i *impl[T]) Insert(x, y float64, data T) {
	if i.Root == nil {
		i.Root = &node[T]{
			Point: point{
				X: x,
				Y: y,
			},
			Data: data,
			Dim:  false,
		}
	}

	i.Root.insert(point{
		X: x,
		Y: y,
	}, data)
}

func (i *impl[T]) SearchNearest(x, y float64) *T {
	if i.Root == nil {
		return nil
	}

	_, data := i.Root.findNearest(point{
		X: x,
		Y: y,
	})

	return &data
}

type node[T any] struct {
	Point point
	Data  T
	Dim   bool
	A     *node[T]
	B     *node[T]
}

func (n *node[T]) compare(p point) bool {
	comp := n.Point.X < p.X
	if n.Dim {
		comp = n.Point.Y < p.Y
	}
	return comp
}

func (n *node[T]) checkHyperplane(p point, dist float64) bool {
	if n.Dim {
		return math.Abs(n.Point.Y-p.Y) < dist
	}
	return math.Abs(n.Point.X-p.X) < dist
}

func (n *node[T]) insert(newP point, data T) {
	if n.compare(newP) {
		if n.A == nil {
			n.A = &node[T]{
				Point: point{
					X: newP.X,
					Y: newP.Y,
				},
				Dim:  !n.Dim,
				Data: data,
			}
		} else {
			n.A.insert(newP, data)
		}
	} else {
		if n.B == nil {
			n.B = &node[T]{
				Point: point{
					X: newP.X,
					Y: newP.Y,
				},
				Dim:  !n.Dim,
				Data: data,
			}
		} else {
			n.B.insert(newP, data)
		}
	}
}

func (n *node[T]) checkLeaf(p point) (float64, T) {
	return n.Point.dist(p), n.Data
}

func (n *node[T]) checkTree(p point, first *node[T], second *node[T]) (float64, T) {
	dist, data := first.findNearest(p)

	if d := n.Point.dist(p); d < dist {
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

func (n *node[T]) findNearest(p point) (float64, T) {
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

type point struct {
	X float64
	Y float64
}

func (p point) dist(other point) float64 {
	return math.Sqrt(math.Pow(p.X-other.X, 2) + math.Pow(p.Y-other.Y, 2))
}

func (p point) distXY(other point) point {
	return point{
		X: math.Abs(p.X - other.X),
		Y: math.Abs(p.Y - other.Y),
	}
}
