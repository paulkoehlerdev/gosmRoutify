package importer

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/graphservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/nodeservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmdataservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmfilterservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
)

const counterLogBreak = 1000000

type Importer interface {
	RunDataImport() error
}

type impl struct {
	osmdataService   osmdataservice.OsmDataService
	osmfilterService osmfilterservice.OsmFilterService
	graphService     graphservice.GraphService
	nodeService      nodeservice.NodeService
	logger           logging.Logger
}

func New(
	osmDataService osmdataservice.OsmDataService,
	nodeService nodeservice.NodeService,
	osmfilterService osmfilterservice.OsmFilterService,
	graphService graphservice.GraphService,
	logger logging.Logger,
) Importer {
	return &impl{
		logger:           logger,
		osmdataService:   osmDataService,
		osmfilterService: osmfilterService,
		nodeService:      nodeService,
		graphService:     graphService,
	}
}

func (i *impl) RunDataImport() error {
	wayPassProcessor := NewWayPassProcessor(i.logger.WithAttrs("processor", "waypass"), i.nodeService, i.osmfilterService)

	err := i.osmdataService.Process(wayPassProcessor)
	if err != nil {
		return fmt.Errorf("error while processing way pass: %s", err.Error())
	}

	i.nodeService.PrintNodeTypeStatistics()

	graphPassProcessor := NewGraphPassProcessor(
		i.logger.WithAttrs("processor", "graphpass"),
		i.nodeService,
		i.osmfilterService,
		i.graphService.AddEdge,
	)

	err = i.osmdataService.Process(graphPassProcessor)
	if err != nil {
		return fmt.Errorf("error while processing graph pass: %s", err.Error())
	}

	return nil
}
