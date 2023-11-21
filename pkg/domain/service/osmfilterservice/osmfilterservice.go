package osmfilterservice

import "github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"

type OsmFilterService interface {
	NodeFilter(node *osmpbfreaderdata.Node) bool
	WayFilter(way *osmpbfreaderdata.Way) bool
	RelationFilter(relation *osmpbfreaderdata.Relation) bool
}

type impl struct {
	ignoredHighwayTypes map[string]struct{}
}

func New(ignoredHighwayTypes []string) OsmFilterService {
	ignoredHighwayTypesMap := make(map[string]struct{})
	for _, highwayType := range ignoredHighwayTypes {
		ignoredHighwayTypesMap[highwayType] = struct{}{}
	}

	return &impl{
		ignoredHighwayTypes: ignoredHighwayTypesMap,
	}
}

func (i *impl) NodeFilter(*osmpbfreaderdata.Node) bool {
	return true
}

func (i *impl) WayFilter(way *osmpbfreaderdata.Way) bool {
	if len(way.NodeIDs) < 2 {
		return false
	}

	if len(way.Tags) == 0 {
		return false
	}

	highwayType, ok := way.Tags["highway"]
	if !ok {
		return false
	}

	if _, ok := i.ignoredHighwayTypes[highwayType]; ok {
		return false
	}

	return true
}

func (i *impl) RelationFilter(*osmpbfreaderdata.Relation) bool {
	return true
}
