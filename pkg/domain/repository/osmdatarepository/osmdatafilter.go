package osmdatarepository

import "github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/filter"

type OsmDataFilter interface {
	filter.Filter
}
