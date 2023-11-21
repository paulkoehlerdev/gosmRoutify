package importer

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/nodetype"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/nodeservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmfilterservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinate"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinatelist"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/nodetags"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
)

type GraphPassProcessor struct {
	logger           logging.Logger
	nodeService      nodeservice.NodeService
	filterService    osmfilterservice.OsmFilterService
	includedNodeTags map[string]struct{}
	edgeHandler      edgeHandler
	nodes            int
	acceptedNodes    int
	ways             int
	acceptedWays     int
}

func NewGraphPassProcessor(logger logging.Logger, nodeService nodeservice.NodeService, filterService osmfilterservice.OsmFilterService, edgeHandler edgeHandler) osmdatarepository.OsmDataProcessor {
	includedNodeTags := []string{"ford", "barrier", "highway", "crossing"}
	includedNodeTagsMap := make(map[string]struct{})
	for _, tag := range includedNodeTags {
		includedNodeTagsMap[tag] = struct{}{}
	}

	return &GraphPassProcessor{
		logger:           logger,
		nodeService:      nodeService,
		filterService:    filterService,
		includedNodeTags: includedNodeTagsMap,
		edgeHandler:      edgeHandler,
		nodes:            0,
		acceptedNodes:    0,
	}
}

func (g *GraphPassProcessor) ProcessNode(node *osmpbfreaderdata.Node) {
	g.nodes++
	if g.nodes%counterLogBreak == 0 {
		g.logger.Info().Msgf("Processed %d Mio. nodes, accepted nodes: %d", g.nodes/1000000, g.acceptedNodes)
	}

	if !g.filterService.NodeFilter(node) {
		return
	}

	nodeType := g.nodeService.SetCoordinate(node.ID, coordinate.New(node.Lat, node.Lon))
	if nodeType == nodetype.EMPTYNODE {
		return
	}

	g.acceptedNodes++

	if isBarrierNode(node) {
		if nodeType != nodetype.JUNCTIONNODE {
			g.nodeService.SetSplitNode(node.ID)
		}
	}

	for tag := range node.Tags {
		if _, ok := g.includedNodeTags[tag]; ok {
			err := g.nodeService.SetTags(node.ID, node.Tags)
			if err != nil {
				return
			}
			break
		}
	}
}

func isBarrierNode(node *osmpbfreaderdata.Node) bool {
	if _, ok := node.Tags["barrier"]; ok {
		return true
	}
	if _, ok := node.Tags["ford"]; ok {
		return true
	}
	return false
}

func (g *GraphPassProcessor) ProcessWay(way *osmpbfreaderdata.Way) {
	g.ways++
	if g.ways%counterLogBreak == 0 {
		g.logger.Info().Msgf("Processed %d Mio. ways, accepted ways: %d", g.ways/1000000, g.acceptedWays)
	}

	if !g.filterService.WayFilter(way) {
		return
	}

	g.acceptedWays++

	segment := make([]segmentNode, 0, len(way.NodeIDs))
	for _, nodeID := range way.NodeIDs {
		coo, err := g.nodeService.GetCoordinate(nodeID)
		if err != nil {
			coo = coordinate.New(0, 0)
		}

		tags, err := g.nodeService.GetTags(nodeID)
		if err != nil {
			tags = nil
		}

		segment = append(segment, segmentNode{
			nodeType:   g.nodeService.GetNodeType(nodeID),
			osmID:      nodeID,
			coordinate: coo,
			tags:       tags,
		})
	}

	g.splitWayAtJunctionsAndEmptySections(segment, way)
}

func (g *GraphPassProcessor) splitWayAtJunctionsAndEmptySections(fullSegment []segmentNode, way *osmpbfreaderdata.Way) {
	var segment []segmentNode
	for _, node := range fullSegment {
		if node.nodeType.IsEmpty() {
			if len(segment) > 1 {
				g.splitLoopSegments(segment, way)
			}
			segment = nil
			continue
		}

		if node.nodeType.IsTowerNode() {
			if len(segment) > 0 {
				segment = append(segment, node)
				g.splitLoopSegments(segment, way)
				segment = nil
			}
			segment = append(segment, node)
			continue
		}

		segment = append(segment, node)
	}

	if len(segment) > 1 {
		g.splitLoopSegments(segment, way)
	}
}

func (g *GraphPassProcessor) splitLoopSegments(segmentNodes []segmentNode, way *osmpbfreaderdata.Way) {
	if len(segmentNodes) < 2 {
		panic(fmt.Errorf("segmentNodes must have at least 2 elements, but was: %d", len(segmentNodes)))
	}

	isLoop := segmentNodes[0].osmID == segmentNodes[len(segmentNodes)-1].osmID
	if isLoop && len(segmentNodes) == 2 {
		return
	}

	if isLoop {
		g.splitSegmentAtSplitNodes(segmentNodes[:len(segmentNodes)-2], way)
		g.splitSegmentAtSplitNodes(segmentNodes[len(segmentNodes)-3:], way)
		return
	}

	g.splitSegmentAtSplitNodes(segmentNodes, way)
}

func (g *GraphPassProcessor) splitSegmentAtSplitNodes(segmentNodes []segmentNode, way *osmpbfreaderdata.Way) {
	var segment []segmentNode
	for _, node := range segmentNodes {
		if g.nodeService.IsSplitNode(node.osmID) {
			g.nodeService.UnsetSplitNode(node.osmID)

			barrierFrom := node
			barrierTo := segmentNode{
				osmID:      -node.osmID,
				coordinate: node.coordinate,
				nodeType:   node.nodeType,
			}

			if len(segment) > 0 {
				segment = append(segment, barrierFrom)
				g.handleSegment(segment, way)
				segment = nil
			}

			way.Tags["barrier"] = "yes"
			segment = append(segment, barrierFrom)
			segment = append(segment, barrierTo)
			g.handleSegment(segment, way)
			delete(way.Tags, "barrier")

			segment = nil
			segment = append(segment, barrierTo)

			continue
		}

		segment = append(segment, node)
	}

	if len(segment) > 1 {
		g.handleSegment(segment, way)
	}
}

func (g *GraphPassProcessor) handleSegment(segmentNodes []segmentNode, way *osmpbfreaderdata.Way) {

	coordinateList := coordinatelist.NewCoordinateList(len(segmentNodes))

	first := segmentNodes[0]
	last := segmentNodes[len(segmentNodes)-1]

	nodeTags := make([]nodetags.NodeTags, 0, len(segmentNodes))
	for _, node := range segmentNodes {
		coordinateList.Append(node.coordinate)
		nodeTags = append(nodeTags, node.tags)
	}

	g.edgeHandler(first.osmID, last.osmID, coordinateList, nodeTags, way)
}

func (g *GraphPassProcessor) ProcessRelation(*osmpbfreaderdata.Relation) {
}

func (g *GraphPassProcessor) OnFinish() {
	g.logger.Info().Msgf("Processed %d nodes, accepted: %d", g.nodes, g.acceptedNodes)
}
