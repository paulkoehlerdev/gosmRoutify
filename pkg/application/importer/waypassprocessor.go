package importer

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/nodetype"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/nodeservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmfilterservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/osmid"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
)

type WayPassProcessor struct {
	logger        logging.Logger
	nodeService   nodeservice.NodeService
	filterService osmfilterservice.OsmFilterService
	ways          int
	acceptedWays  int
	nodes         int
}

func NewWayPassProcessor(logger logging.Logger, nodeService nodeservice.NodeService, filterService osmfilterservice.OsmFilterService) osmdatarepository.OsmDataProcessor {
	return &WayPassProcessor{
		logger:        logger,
		filterService: filterService,
		nodeService:   nodeService,
		ways:          0,
		acceptedWays:  0,
	}
}

func (w *WayPassProcessor) ProcessNode(osmpbfreaderdata.Node) {
	w.logger.Warn().Msgf("WayPassProcessor.ProcessNode() should not be called")
}

func (w *WayPassProcessor) ProcessWay(way osmpbfreaderdata.Way) {
	w.ways++
	if w.ways%counterLogBreak == 0 {
		w.logger.Info().Msgf("Processed %d Mio. ways, accepted ways: %d", w.ways/1000000, w.acceptedWays)
	}

	if !w.filterService.WayFilter(way) {
		return
	}
	w.acceptedWays++

	for index, nodeID := range way.NodeIDs {
		isEnd := index == 0 || index == len(way.NodeIDs)-1

		nodeType := nodetype.ENDNODE
		if !isEnd {
			nodeType = nodetype.INTERMEDIATENODE
		}

		w.nodeService.AddOrUpdate(osmid.OsmID(nodeID), nodeType, func(prev nodetype.NodeType) nodetype.NodeType {
			if prev == nodetype.ENDNODE && isEnd {
				return nodetype.CONNECTIONNODE
			}
			return nodetype.JUNCTIONNODE
		})
	}
}

func (w *WayPassProcessor) ProcessRelation(osmpbfreaderdata.Relation) {
	w.logger.Warn().Msgf("WayPassProcessor.ProcessRelation() should not be called")
}

func (w *WayPassProcessor) OnFinish() {
	w.logger.Info().Msgf("Processed %d ways, accepted: %d", w.ways, w.acceptedWays)
}
