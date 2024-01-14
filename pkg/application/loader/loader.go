package loader

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/addressService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/nodeService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmdataservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/wayService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
)

type Loader interface {
	Load() error
}

type impl struct {
	dataService    osmdataservice.OsmDataService
	nodeService    nodeService.NodeService
	addressService addressService.AddressService
	wayService     wayService.WayService
	logger         logging.Logger

	nodeCount int
	wayCount  int
}

func New(dataService osmdataservice.OsmDataService, nodeService nodeService.NodeService, wayService wayService.WayService, addressService addressService.AddressService, logger logging.Logger) Loader {
	return &impl{
		dataService:    dataService,
		nodeService:    nodeService,
		addressService: addressService,
		wayService:     wayService,
		logger:         logger,
		nodeCount:      0,
	}
}

func (i *impl) Load() error {
	firstPassProcessor := newFirstPassProcessor(
		i.wayService,
		i.addressService,
		i.logger,
	)
	firstPassFilter := osmdatarepository.NewBinaryOsmDataFilter(
		true, false, true,
	)

	secondPassProcessor := newSecondPassProcessor(
		i.wayService,
		i.nodeService,
		i.addressService,
		i.logger,
	)
	secondPassFilter := osmdatarepository.NewBinaryOsmDataFilter(
		false, true, true,
	)

	i.logger.Info().Msgf("Starting import!")
	i.logger.Info().Msgf("First pass: inserting ways and addresses")

	err := i.dataService.Process(firstPassProcessor, firstPassFilter)
	if err != nil {
		return fmt.Errorf("error while processing first pass: %s", err.Error())
	}

	i.logger.Info().Msgf("Second pass: inserting nodes")

	err = i.dataService.Process(secondPassProcessor, secondPassFilter)
	if err != nil {
		return fmt.Errorf("error while processing second pass: %s", err.Error())
	}

	return nil
}
