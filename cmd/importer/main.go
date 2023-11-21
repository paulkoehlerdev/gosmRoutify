package main

import (
	"flag"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/application/importer"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/noderepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/nodeservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmdataservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmfilterservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"runtime"
)

func main() {
	dataPath := flag.String("data", "./resources/sample/sample.pbf", "path to data file")
	graphPath := flag.String("graph", "./resources/graph/", "path to graph folder")
	flag.Parse()

	logger := logging.New(logging.LevelDebug, logging.NewConsoleWriter())

	if dataPath == nil {
		logger.Error().Msg("no data file provided")
		return
	}

	if graphPath == nil {
		logger.Error().Msg("no graph folder provided")
		return
	}

	osmdataRepo := osmdatarepository.New(runtime.GOMAXPROCS(-1))
	osmdataSvc := osmdataservice.New(osmdataRepo, []string{*dataPath}, logger.WithAttrs("service", "osmdata"))

	nodeRepo := noderepository.New(logger.WithAttrs("repository", "node"))
	nodeSvc := nodeservice.New(nodeRepo, logger.WithAttrs("service", "node"))

	osmfilterSvc := osmfilterservice.New([]string{})

	application := importer.New(osmdataSvc, nodeSvc, osmfilterSvc, logger.WithAttrs("application", "importer"))

	err := application.RunDataImport()
	if err != nil {
		logger.Error().Msgf("error while running data import: %s", err.Error())
		return
	}
}
