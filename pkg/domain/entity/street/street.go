package street

import "github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/osmpbfData"

type Street struct {
	Way  *osmpbfData.Way
	Node []*osmpbfData.Node
}
