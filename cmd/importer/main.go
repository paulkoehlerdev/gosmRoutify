package main

import (
	"flag"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/application/importer"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/config"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/noderepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/tilerepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/graphservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/nodeservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmdataservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmfilterservice"
	_ "net/http/pprof"
	"runtime"
)

const tileSize = 0.12 //degree
const maxCacheSize = 400

func main() {
	configPath := flag.String("config", "config.json", "path to config file")
	flag.Parse()

	if configPath == nil {
		panic("no config file provided")
	}

	config, err := config.FromFile(*configPath)
	if err != nil {
		panic(fmt.Errorf("error while loading config: %s", err.Error()))
	}

	logger := config.LoggerConfig.SetupLogger()

	logger.Info().Msg("starting importing engine")

	osmdataRepo := osmdatarepository.New(runtime.GOMAXPROCS(-1))
	osmdataSvc := osmdataservice.New(
		osmdataRepo,
		[]string{config.ImporterConfig.FilePath},
		logger.WithAttrs("service", "osmdata"),
	)

	nodeRepo, err := noderepository.New(
		fmt.Sprintf("%s/gosmRoutifyNodes.db", config.ImporterConfig.TmpFilePath),
		config.ImporterConfig.EnableDiskStorage,
		logger.WithAttrs("repository", "node"),
	)
	if err != nil {
		logger.Error().Msgf("error while creating node repository: %s", err.Error())
		return
	}

	nodeSvc := nodeservice.New(nodeRepo, logger.WithAttrs("service", "node"))

	osmfilterSvc := osmfilterservice.New([]string{})

	tileRepo := tilerepository.New(
		logger.WithAttrs("repository", "tile"),
		tileSize,
		config.GraphConfig.Path,
		maxCacheSize,
	)

	graphSvc := graphservice.New(tileRepo, logger.WithAttrs("service", "graph"))

	application := importer.New(
		osmdataSvc,
		nodeSvc,
		osmfilterSvc,
		graphSvc,
		logger.WithAttrs("application", "importer"),
	)

	err = application.RunDataImport()
	if err != nil {
		logger.Error().Msgf("error while running data import: %s", err.Error())
		return
	}

	err = nodeRepo.Close()
	if err != nil {
		logger.Error().Msgf("error while closing node repository: %s", err.Error())
	}

	tileRepo.Close()
}
