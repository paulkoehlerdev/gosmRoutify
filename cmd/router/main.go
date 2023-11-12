package main

import (
	"flag"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/config"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
)

func main() {
	configPath := flag.String("config", "config.json", "path to config file")
	graphPath := flag.String("graph", "./resources/graph/", "path to graph folder")
	flag.Parse()

	if configPath == nil {
		panic("no config file provided")
	}

	if graphPath == nil {
		panic("no graph folder provided")
	}

	config, err := config.FromFile(*configPath)
	if err != nil {
		panic(fmt.Errorf("error while loading config: %s", err.Error()))
	}

	logger := setupLogger(config.LoggerConfig)

	logger.Info().Msg("starting routing engine")
}

func setupLogger(loggerConfig *config.LoggerConfig) logging.Logger {
	if loggerConfig == nil {
		panic("no logger config provided")
	}

	writer, fileWriterErr := logging.NewFileWriter(loggerConfig.FilePath)

	if loggerConfig.Console {
		if writer != nil {
			writer = logging.NewMultiWriter(logging.NewConsoleWriter(), writer)
		} else {
			writer = logging.NewConsoleWriter()
		}
	}

	logger := logging.New(logging.LogLevel(loggerConfig.Level), writer)

	if fileWriterErr != nil {
		logger.Warn().Msgf("no file writer configured: no logs file will be written: %s", fileWriterErr.Error())
	}

	return logger
}
