package crossing

import "github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/node"

type Crossing struct {
	node.Node
	IsCrossing bool
}
