package relation

type Relation struct {
	OsmID int64
	Tags  map[string]string
	Nodes []int64
	Ways  []int64
}
