package main

import (
	"flag"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/application/router"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/config"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/tilerepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/graphservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/weightingservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/interface/http"
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

	logger.Info().Msg("starting routing engine")

	tileRepo := tilerepository.New(
		logger.WithAttrs("repository", "tile"),
		tileSize,
		config.GraphConfig.Path,
		maxCacheSize,
	)

	graphSvc := graphservice.New(tileRepo, logger.WithAttrs("service", "graph"))

	weightingSvc := weightingservice.New()

	application := router.New(graphSvc, weightingSvc, logger.WithAttrs("application", "router"))

	server, err := http.NewHttpServer(logger.WithAttrs("service", "interfaceHTTP"), application, config.ServerConfig)
	if err != nil {
		logger.Error().Msgf("error while creating http server: %s", err.Error())
	}

	logger.Info().Msg("loaded interfaceHTTP")
	err = server.ListenAndServe()
	if err != nil {
		logger.Info().Msgf("error while serving http server: %s", err.Error())
	}
}
