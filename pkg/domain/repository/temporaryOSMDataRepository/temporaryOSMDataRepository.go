package temporaryOSMDataRepository

import "github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/osmpbfData"

type TemporaryOSMDataRepository interface {
	AddOSMData(data any) error
	Next() (*osmpbfData.Way, error)
	FindNode(osmID int64) *osmpbfData.Node
	Cleanup()
}
