package loader

import (
	nodeModel "github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/node"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/nodeService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/wayService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
)

type secondPassProcessor struct {
	wayService        wayService.WayService
	nodeService       nodeService.NodeService
	logger            logging.Logger
	nodeCount         int
	acceptedNodeCount int
}

func newSecondPassProcessor(wayService wayService.WayService, nodeService nodeService.NodeService, logger logging.Logger) osmdatarepository.OsmDataProcessor {
	return &secondPassProcessor{
		wayService:  wayService,
		nodeService: nodeService,
		logger:      logger,
	}
}

func (i *secondPassProcessor) ProcessNode(node osmpbfreaderdata.Node) {
	newNode := nodeModel.Node{
		OsmID: node.ID,
		Lat:   node.Lat,
		Lon:   node.Lon,
		Tags:  node.Tags,
	}

	i.nodeCount++
	if i.nodeCount%100000 == 0 {
		i.logger.Info().Msgf("Inserted %d nodes, accepted %d", i.nodeCount, i.acceptedNodeCount)
	}

	ways, err := i.wayService.SelectWayIDsFromNode(newNode.OsmID)
	if err != nil || len(ways) == 0 {
		return
	}

	i.acceptedNodeCount++

	err = i.nodeService.InsertNodeBulk(newNode)
	if err != nil {
		i.logger.Error().Msgf("Error while inserting node: %s", err.Error())
		return
	}
}

func (i *secondPassProcessor) ProcessWay(_ osmpbfreaderdata.Way) {}

func (i *secondPassProcessor) ProcessRelation(_ osmpbfreaderdata.Relation) {
}

func (i *secondPassProcessor) OnFinish() {
	err := i.nodeService.CommitBulkInsert()
	if err != nil {
		i.logger.Error().Msgf("Error while committing bulk insert: %s", err.Error())
		return
	}
	i.logger.Info().Msgf("Inserted %d nodes, accepted %d", i.nodeCount, i.acceptedNodeCount)
}
