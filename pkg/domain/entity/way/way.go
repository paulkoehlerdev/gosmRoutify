package way

type Way struct {
	OsmID int64
	Tags  map[string]string
	Nodes []int64
}
