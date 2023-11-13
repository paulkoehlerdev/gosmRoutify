package router

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinate"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/route"
)

type Application interface {
	FindRoute(start coordinate.Coordinate, end coordinate.Coordinate) (route.Route, error)
}

type impl struct {
}

func New() Application {
	return &impl{}
}

func (i *impl) FindRoute(start coordinate.Coordinate, end coordinate.Coordinate) (route.Route, error) {
	return nil, fmt.Errorf("not implemented")
}
