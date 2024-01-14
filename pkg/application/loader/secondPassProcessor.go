package loader

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/address"
	nodeModel "github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/node"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/addressService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/nodeService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/wayService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
)

type secondPassProcessor struct {
	wayService        wayService.WayService
	nodeService       nodeService.NodeService
	addressService    addressService.AddressService
	logger            logging.Logger
	nodeCount         int
	acceptedNodeCount int
}

func newSecondPassProcessor(wayService wayService.WayService, nodeService nodeService.NodeService, addressService addressService.AddressService, logger logging.Logger) osmdatarepository.OsmDataProcessor {
	return &secondPassProcessor{
		wayService:     wayService,
		nodeService:    nodeService,
		addressService: addressService,
		logger:         logger,
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
	if i.nodeCount%1000000 == 0 {
		i.logger.Info().Msgf("Inserted %d nodes, accepted %d", i.nodeCount, i.acceptedNodeCount)
	}

	address, addrErr := i.getAddressFromNode(node)
	ways, wayErr := i.wayService.SelectWayIDsFromNode(newNode.OsmID)
	if !((wayErr == nil && len(ways) != 0) || (addrErr == nil && address != nil)) {
		return
	}

	if addrErr == nil && address != nil {
		err := i.addressService.InsertAddressBulk(*address)
		if err != nil {
			i.logger.Error().Msgf("Error while inserting address: %s", err.Error())
			return
		}
	}

	i.acceptedNodeCount++

	err := i.nodeService.InsertNodeBulk(newNode)
	if err != nil {
		i.logger.Error().Msgf("Error while inserting node: %s", err.Error())
		return
	}
}

func (i *secondPassProcessor) ProcessWay(_ osmpbfreaderdata.Way) {}

func (i *secondPassProcessor) ProcessRelation(_ osmpbfreaderdata.Relation) {
}

func (i *secondPassProcessor) OnFinish() {
	err := i.addressService.CommitBulkInsert()
	if err != nil {
		i.logger.Error().Msgf("Error while committing bulk insert: %s", err.Error())
	}

	err = i.nodeService.CommitBulkInsert()
	if err != nil {
		i.logger.Error().Msgf("Error while committing bulk insert: %s", err.Error())
		return
	}

	err = i.nodeService.CreateIndices()
	if err != nil {
		i.logger.Error().Msgf("Error while creating indices: %s", err.Error())
		return
	}

	i.logger.Info().Msgf("Inserted %d nodes, accepted %d", i.nodeCount, i.acceptedNodeCount)
}

func (i *secondPassProcessor) getAddressFromNode(node osmpbfreaderdata.Node) (*address.Address, error) {
	address, err := getAddressFromTags(node.Tags)
	if err != nil {
		return nil, fmt.Errorf("error while getting address from node tags: %s", err.Error())
	}

	address.OsmID = node.ID

	return address, nil
}
