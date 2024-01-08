package loader

import (
	wayModel "github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/way"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/nodeService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/wayService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
)

type firstPassProcessor struct {
	wayService       wayService.WayService
	nodeService      nodeService.NodeService
	logger           logging.Logger
	wayCount         int
	acceptedWayCount int
}

func newFirstPassProcessor(wayService wayService.WayService, logger logging.Logger) osmdatarepository.OsmDataProcessor {
	return &firstPassProcessor{
		wayService: wayService,
		logger:     logger,
		wayCount:   0,
	}
}

func (i *firstPassProcessor) ProcessNode(_ osmpbfreaderdata.Node) {}

func (i *firstPassProcessor) ProcessWay(way osmpbfreaderdata.Way) {
	newWay := wayModel.Way{
		OsmID: way.ID,
		Tags:  way.Tags,
		Nodes: way.NodeIDs,
	}

	i.wayCount++
	if i.wayCount%100000 == 0 {
		i.logger.Info().Msgf("Inserted %d ways, accepted %d", i.wayCount, i.acceptedWayCount)
	}

	if _, ok := way.Tags["highway"]; !ok {
		return
	}

	i.acceptedWayCount++

	err := i.wayService.InsertWayBulk(newWay)
	if err != nil {
		i.logger.Error().Msgf("Error while inserting way: %s", err.Error())
		return
	}
}

func (i *firstPassProcessor) ProcessRelation(_ osmpbfreaderdata.Relation) {
}

func (i *firstPassProcessor) OnFinish() {
	err := i.wayService.CommitBulkInsert()
	if err != nil {
		i.logger.Error().Msgf("Error while committing bulk insert: %s", err.Error())
	}
}
