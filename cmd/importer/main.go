package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/application/importer"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/graphService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmDataPreprocessorService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmDataService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/infrastructure/graphRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/infrastructure/osmDataRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/infrastructure/temporaryOSMDataRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"io"
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

	osmDataRepo := osmDataRepository.New(runtime.GOMAXPROCS(-1))

	osmDataSvc, err := osmDataService.New(osmDataRepo, []string{*dataPath}, logger.WithAttrs("service", "osmDataService"))
	if err != nil {
		logger.Error().Msgf("error while initializing data service %s", err.Error())
		return
	}

	preprocessorSvc := osmDataPreprocessorService.New()

	tempOSMRepo := temporaryOSMDataRepository.New(fmt.Sprintf("%s/tmp/", *graphPath), logger.WithAttrs("service", "temporaryOSMDataRepository"))
	graphRepo := graphRepository.New(*graphPath, logger.WithAttrs("service", "graphRepository"))

	graphSvc := graphService.New(graphRepo, tempOSMRepo, logger.WithAttrs("service", "graphService"))

	application := importer.New(osmDataSvc, preprocessorSvc, graphSvc, logger.WithAttrs("service", "application"))
	err = application.StartDataImport()
	if errors.Is(err, io.EOF) {
		logger.Info().Msg("finished importing data")
	} else if err != nil {
		logger.Error().Msgf("error while importing data %s", err.Error())
		return
	}

	application.BuildGraph()

	logger.Info().Msg("stopping service")

	osmDataRepo.Stop()
}
