package loader

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/address"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/addressService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
)

type thirdPassProcessor struct {
	addressService       addressService.AddressService
	logger               logging.Logger
	addressCount         int
	acceptedAddressCount int
}

func newThirdPassProcessor(addressService addressService.AddressService, logger logging.Logger) osmdatarepository.OsmDataProcessor {
	return &thirdPassProcessor{
		addressService: addressService,
		logger:         logger,
	}
}

func (i *thirdPassProcessor) ProcessNode(node osmpbfreaderdata.Node) {
	i.addressCount++
	if i.addressCount%1000000 == 0 {
		i.logger.Info().Msgf("Inserted %d addresses, accepted %d", i.addressCount, i.acceptedAddressCount)
	}

	address, err := i.getAddressFromNode(node)
	if err != nil {
		return
	}

	if address == nil {
		return
	}

	i.acceptedAddressCount++

	err = i.addressService.InsertAddressBulk(*address)
	if err != nil {
		i.logger.Error().Msgf("Error while inserting address: %s", err.Error())
		return
	}
}

func (i *thirdPassProcessor) ProcessWay(way osmpbfreaderdata.Way) {
	i.addressCount++
	if i.addressCount%1000000 == 0 {
		i.logger.Info().Msgf("Inserted %d addresses, accepted %d", i.addressCount, i.acceptedAddressCount)
	}

	address, err := i.getAddressFromWay(way)
	if err != nil {
		return
	}

	if address == nil {
		return
	}

	i.acceptedAddressCount++

	err = i.addressService.InsertAddressBulk(*address)
	if err != nil {
		i.logger.Error().Msgf("Error while inserting address: %s", err.Error())
		return
	}
}

func (i *thirdPassProcessor) ProcessRelation(_ osmpbfreaderdata.Relation) {
}

func (i *thirdPassProcessor) OnFinish() {
	err := i.addressService.CommitBulkInsert()
	if err != nil {
		i.logger.Error().Msgf("Error while committing bulk insert: %s", err.Error())
	}
}

func (i *thirdPassProcessor) getAddressFromNode(node osmpbfreaderdata.Node) (*address.Address, error) {
	address, err := getAddressFromTags(node.Tags)
	if err != nil {
		return nil, fmt.Errorf("error while getting address from node tags: %s", err.Error())
	}

	address.OsmID = node.ID

	return address, nil
}

func (i *thirdPassProcessor) getAddressFromWay(way osmpbfreaderdata.Way) (*address.Address, error) {
	address, err := getAddressFromTags(way.Tags)
	if err != nil {
		return nil, fmt.Errorf("error while getting address from node tags: %s", err.Error())
	}

	address.OsmID = way.ID

	return address, nil
}
