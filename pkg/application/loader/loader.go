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
	Load(data, address bool) error
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

func (i *impl) Load(data bool, address bool) error {
	firstPassProcessor := newFirstPassProcessor(
		i.wayService,
		i.logger,
	)
	firstPassFilter := osmdatarepository.NewBinaryOsmDataFilter(
		true, false, true,
	)

	secondPassProcessor := newSecondPassProcessor(
		i.wayService,
		i.nodeService,
		i.logger,
	)
	secondPassFilter := osmdatarepository.NewBinaryOsmDataFilter(
		false, true, true,
	)

	thirdPassProcessor := newThirdPassProcessor(
		i.addressService,
		i.logger,
	)
	thirdPassFilter := osmdatarepository.NewBinaryOsmDataFilter(
		false, true, true,
	)

	if data {
		err := i.dataService.Process(firstPassProcessor, firstPassFilter)
		if err != nil {
			return fmt.Errorf("error while processing first pass: %s", err.Error())
		}

		err = i.dataService.Process(secondPassProcessor, secondPassFilter)
		if err != nil {
			return fmt.Errorf("error while processing second pass: %s", err.Error())
		}
	}

	if address {
		err := i.dataService.Process(thirdPassProcessor, thirdPassFilter)
		if err != nil {
			return fmt.Errorf("error while processing third pass: %s", err.Error())
		}
	}

	return nil
}
