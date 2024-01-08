package node

type Node struct {
	OsmID int64
	Lat   float64
	Lon   float64
	Tags  map[string]string
}
