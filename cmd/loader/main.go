//go:build json && fts5

package main

import (
	"flag"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/application/loader"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/addressRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/nodeRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/wayRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/addressService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/nodeService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmdataservice"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/wayService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/database"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"os"
	"runtime"
)

func main() {
	importFile := flag.String("import", "", "import file")
	databaseFile := flag.String("database", "", "database file")

	importData := flag.Bool("data", false, "import data")
	importAddress := flag.Bool("address", false, "import address")

	flag.Parse()

	if importFile == nil || databaseFile == nil {
		panic("no import or database file provided")
	}

	logger := logging.New(logging.LevelDebug, os.Stdout)

	if importData == nil {
		importData = new(bool)
		*importData = false
	}

	if importAddress == nil {
		importAddress = new(bool)
		*importAddress = false
	}

	if *importData == false && *importAddress == false {
		logger.Error().Msgf("no import type provided")
		return
	}

	db, err := database.New(*databaseFile)
	if err != nil {
		logger.Error().Msgf("error while creating database: %s", err.Error())
		return
	}
	defer db.Close()

	osmdataRepo := osmdatarepository.New(runtime.GOMAXPROCS(-1))
	osmdataSvc := osmdataservice.New(
		osmdataRepo,
		[]string{*importFile},
		logger.WithAttrs("service", "osmdata"),
	)

	wayRepo := wayRepository.New(db)
	err = wayRepo.Init()
	if err != nil {
		logger.Error().Msgf("error while initializing node repository: %s", err.Error())
		return
	}

	waySvc := wayService.New(wayRepo)

	nodeRepo := nodeRepository.New(db)
	err = nodeRepo.Init()
	if err != nil {
		logger.Error().Msgf("error while initializing node repository: %s", err.Error())
		return
	}

	nodeSvc := nodeService.New(nodeRepo, logger.WithAttrs("service", "node"))

	addrRepo := addressRepository.New(db)
	err = addrRepo.Init()
	if err != nil {
		logger.Error().Msgf("error while initializing address repository: %s", err.Error())
		return
	}

	addrSvc := addressService.New(addrRepo, logger.WithAttrs("service", "address"))

	application := loader.New(osmdataSvc, nodeSvc, waySvc, addrSvc, logger.WithAttrs("application", "loader"))

	err = application.Load(*importData, *importAddress)
	if err != nil {
		logger.Error().Msgf("error while loading data: %s", err.Error())
		return
	}
}
