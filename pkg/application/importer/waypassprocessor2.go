package importer

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/nodetype"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/noderepository2"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmfilterservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/osmid"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
)

type WayPassProcessor2 struct {
	logger         logging.Logger
	filterService  osmfilterservice.OsmFilterService
	ways           int
	acceptedWays   int
	noderepository noderepository2.NodeRepository
}

func NewWayPassProcessor2(logger logging.Logger, filterService osmfilterservice.OsmFilterService) osmdatarepository.OsmDataProcessor {
	repo, err := noderepository2.New("test.db")
	if err != nil {
		panic(err)
	}

	return &WayPassProcessor2{
		logger:         logger,
		filterService:  filterService,
		ways:           0,
		acceptedWays:   0,
		noderepository: repo,
	}
}

func (w *WayPassProcessor2) ProcessNode(_ osmpbfreaderdata.Node) {
	w.logger.Warn().Msgf("WayPassProcessor2.ProcessNode() should not be called")
}

func (w *WayPassProcessor2) ProcessWay(way osmpbfreaderdata.Way) {
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

		data, err := w.noderepository.GetData(osmid.OsmID(nodeID))
		if err == nil {
			if data.Type == nodetype.ENDNODE && isEnd {
				nodeType = nodetype.CONNECTIONNODE
			} else {
				nodeType = nodetype.JUNCTIONNODE
			}
		}

		if data == nil {
			data = new(noderepository2.NodeData)
		}
		data.Type = nodeType

		err = w.noderepository.SetData(osmid.OsmID(nodeID), *data)
		if err != nil {
			panic(err)
		}
	}
}

func (w *WayPassProcessor2) ProcessRelation(_ osmpbfreaderdata.Relation) {
	w.logger.Warn().Msgf("WayPassProcessor2.ProcessRelation() should not be called")
}

func (w *WayPassProcessor2) OnFinish() {
	w.logger.Info().Msgf("Processed %d Mio. ways, accepted ways: %d", w.ways/1000000, w.acceptedWays)
}
