package loader

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/address"
	wayModel "github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/way"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/addressService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/wayService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
)

type firstPassProcessor struct {
	wayService       wayService.WayService
	addressService   addressService.AddressService
	logger           logging.Logger
	wayCount         int
	acceptedWayCount int
}

func newFirstPassProcessor(wayService wayService.WayService, addressService addressService.AddressService, logger logging.Logger) osmdatarepository.OsmDataProcessor {
	return &firstPassProcessor{
		wayService:     wayService,
		addressService: addressService,
		logger:         logger,
		wayCount:       0,
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
	if i.wayCount%1000000 == 0 {
		i.logger.Info().Msgf("Inserted %dM ways, accepted %d", i.wayCount/1000000, i.acceptedWayCount)
	}

	address, err := i.getAddressFromWay(way)
	if _, ok := way.Tags["highway"]; !(ok || (err == nil && address != nil)) {
		return
	}

	if err == nil && address != nil {
		err = i.addressService.InsertAddressBulk(*address)
		if err != nil {
			i.logger.Error().Msgf("Error while inserting address: %s", err.Error())
			return
		}
	}

	i.acceptedWayCount++

	err = i.wayService.InsertWayBulk(newWay)
	if err != nil {
		i.logger.Error().Msgf("Error while inserting way: %s", err.Error())
		return
	}
}

func (i *firstPassProcessor) ProcessRelation(_ osmpbfreaderdata.Relation) {
}

func (i *firstPassProcessor) OnFinish() {
	err := i.addressService.CommitBulkInsert()
	if err != nil {
		i.logger.Error().Msgf("Error while committing bulk insert: %s", err.Error())
	}

	err = i.wayService.CommitBulkInsert()
	if err != nil {
		i.logger.Error().Msgf("Error while committing bulk insert: %s", err.Error())
	}

	i.logger.Info().Msgf("Creating Way Indices!")

	err = i.wayService.CreateIndices()
	if err != nil {
		i.logger.Error().Msgf("Error while creating indices: %s", err.Error())
	}

	i.logger.Info().Msgf("Updating Crossings!")

	err = i.wayService.UpdateCrossings()
	if err != nil {
		i.logger.Error().Msgf("Error while updating crossings: %s", err.Error())
	}

	i.logger.Info().Msgf("Inserted %dM ways, accepted %d", i.wayCount/1000000, i.acceptedWayCount)
}

func (i *firstPassProcessor) getAddressFromWay(way osmpbfreaderdata.Way) (*address.Address, error) {
	address, err := getAddressFromTags(way.Tags)
	if err != nil {
		return nil, fmt.Errorf("error while getting address from node tags: %s", err.Error())
	}

	address.OsmID = way.ID

	return address, nil
}
