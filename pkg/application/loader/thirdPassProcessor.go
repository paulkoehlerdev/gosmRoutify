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

func (i *thirdPassProcessor) getAddressFromNode(node osmpbfreaderdata.Node) (*address.Address, error) {
	var address address.Address

	anyAvailable := false

	if val, ok := node.Tags["addr:street"]; ok {
		address.Street = val
		anyAvailable = true
	}

	if val, ok := node.Tags["addr:housenumber"]; ok {
		address.Housenumber = val
		anyAvailable = true
	}

	if val, ok := node.Tags["addr:city"]; ok {
		address.City = val
		anyAvailable = true
	}

	if val, ok := node.Tags["addr:postcode"]; ok {
		address.Postcode = val
		anyAvailable = true
	}

	if val, ok := node.Tags["addr:country"]; ok {
		address.Country = val
		anyAvailable = true
	}

	if val, ok := node.Tags["addr:suburb"]; ok {
		address.Suburb = val
		anyAvailable = true
	}

	if val, ok := node.Tags["addr:state"]; ok {
		address.State = val
		anyAvailable = true
	}

	if val, ok := node.Tags["addr:province"]; ok {
		address.Province = val
		anyAvailable = true
	}

	if val, ok := node.Tags["addr:floor"]; ok {
		address.Floor = val
		anyAvailable = true
	}

	if val, ok := node.Tags["name"]; ok {
		address.Name = val
		anyAvailable = true
	}

	if !anyAvailable {
		return nil, fmt.Errorf("no address found")
	}

	address.Lat = node.Lat
	address.Lon = node.Lon

	return &address, nil
}

func (i *thirdPassProcessor) ProcessWay(_ osmpbfreaderdata.Way) {}

func (i *thirdPassProcessor) ProcessRelation(_ osmpbfreaderdata.Relation) {
}

func (i *thirdPassProcessor) OnFinish() {
	err := i.addressService.CommitBulkInsert()
	if err != nil {
		i.logger.Error().Msgf("Error while committing bulk insert: %s", err.Error())
	}
}
