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

func getAddressFromTags(tags map[string]string) (*address.Address, error) {
	var address address.Address

	anyAvailable := false

	if val, ok := tags["addr:street"]; ok {
		address.Street = val
		anyAvailable = true
	}

	if val, ok := tags["addr:housenumber"]; ok {
		address.Housenumber = val
		anyAvailable = true
	}

	if val, ok := tags["addr:city"]; ok {
		address.City = val
		anyAvailable = true
	}

	if val, ok := tags["addr:postcode"]; ok {
		address.Postcode = val
		anyAvailable = true
	}

	if val, ok := tags["addr:country"]; ok {
		address.Country = val
		anyAvailable = true
	}

	if val, ok := tags["addr:suburb"]; ok {
		address.Suburb = val
		anyAvailable = true
	}

	if val, ok := tags["addr:state"]; ok {
		address.State = val
		anyAvailable = true
	}

	if val, ok := tags["addr:province"]; ok {
		address.Province = val
		anyAvailable = true
	}

	if val, ok := tags["addr:floor"]; ok {
		address.Floor = val
		anyAvailable = true
	}

	if val, ok := tags["name"]; ok {
		address.Name = val
		anyAvailable = true
	}

	if !anyAvailable {
		return nil, fmt.Errorf("no address found")
	}

	return &address, nil
}
