package filter

type Filter interface {
	FilterNodes() bool
	FilterWays() bool
	FilterRelations() bool

	// Possibly add these later, as they could improve memory usage and pointer allocations, as they allocate strings
	// FilterTags() bool
}
