//go:build json && fts5
// +build json,fts5

package main

import (
	"flag"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/application/router"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/config"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/addressRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/crossingRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/nodeRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/wayRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/weightRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/addressService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/graphService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/interface/http"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/database"
)

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

	logger.Info().Msgf("Database file: %s", config.DatabaseConfig.FilePath)
	db, err := database.New(config.DatabaseConfig.FilePath)
	if err != nil {
		logger.Error().Msgf("error while creating database: %s", err.Error())
		return
	}
	defer db.Close()

	nodeRepo := nodeRepository.New(db)
	err = nodeRepo.Init()
	if err != nil {
		logger.Error().Msgf("error while initializing node repository: %s", err.Error())
		return
	}

	// nodeSvc := nodeService.New(nodeRepo, logger.WithAttrs("service", "node"))

	wayRepo := wayRepository.New(db)
	err = wayRepo.Init()
	if err != nil {
		logger.Error().Msgf("error while initializing node repository: %s", err.Error())
		return
	}

	// waySvc := wayService.New(wayRepo)

	crossingRepo := crossingRepository.New(db)
	err = crossingRepo.Init()
	if err != nil {
		logger.Error().Msgf("error while initializing node repository: %s", err.Error())
		return
	}

	weightRepo := weightRepository.New(logger.WithAttrs("repository", "weight"))

	graphSvc := graphService.New(nodeRepo, crossingRepo, wayRepo, weightRepo, logger.WithAttrs("service", "graph"))

	addrRepo := addressRepository.New(db)
	err = addrRepo.Init()
	if err != nil {
		logger.Error().Msgf("error while initializing address repository: %s", err.Error())
		return
	}

	addrSvc := addressService.New(addrRepo, logger.WithAttrs("service", "address"))

	application := router.New(graphSvc, addrSvc, logger.WithAttrs("application", "loader"))

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
