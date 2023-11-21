package importer

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/nodeservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmdataservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmfilterservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinatelist"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/nodetags"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
)

const counterLogBreak = 1000000

type Importer interface {
	RunDataImport() error
}

type impl struct {
	osmdataService   osmdataservice.OsmDataService
	osmfilterService osmfilterservice.OsmFilterService
	nodeService      nodeservice.NodeService
	logger           logging.Logger
}

func New(osmDataService osmdataservice.OsmDataService, nodeService nodeservice.NodeService, osmfilterService osmfilterservice.OsmFilterService, logger logging.Logger) Importer {
	return &impl{
		logger:           logger,
		osmdataService:   osmDataService,
		osmfilterService: osmfilterService,
		nodeService:      nodeService,
	}
}

func (i *impl) RunDataImport() error {
	wayPassProcessor := NewWayPassProcessor(i.logger, i.nodeService, i.osmfilterService)

	err := i.osmdataService.Process(wayPassProcessor)
	if err != nil {
		return fmt.Errorf("error while processing way pass: %s", err.Error())
	}

	i.nodeService.PrintNodeTypeStatistics()

	graphPassProcessor := NewGraphPassProcessor(i.logger, i.nodeService, i.osmfilterService, i.simpleEdgeHandler)

	err = i.osmdataService.Process(graphPassProcessor)
	if err != nil {
		return fmt.Errorf("error while processing graph pass: %s", err.Error())
	}

	return nil
}

func (i *impl) simpleEdgeHandler(fromID, toID int64, list coordinatelist.CoordinateList, _ []nodetags.NodeTags, _ *osmpbfreaderdata.Way) {
	i.logger.Debug().Msgf("Adding edge from %d to %d with %d nodes", fromID, toID, list.Len())
}
