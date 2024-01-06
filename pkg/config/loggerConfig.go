package config

import "github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"

type LoggerConfig struct {
	Level    string `json:"level"`
	FilePath string `json:"file"`
	Console  bool   `json:"console"`
}

func (lc *LoggerConfig) SetupLogger() logging.Logger {
	if lc == nil {
		panic("no logger config provided")
	}

	writer, fileWriterErr := logging.NewFileWriter(lc.FilePath)

	if lc.Console {
		if writer != nil {
			writer = logging.NewMultiWriter(logging.NewConsoleWriter(), writer)
		} else {
			writer = logging.NewConsoleWriter()
		}
	}

	logger := logging.New(logging.LogLevel(lc.Level), writer)

	if fileWriterErr != nil {
		logger.Warn().Msgf("no file writer configured: no logs file will be written: %s", fileWriterErr.Error())
	}

	return logger
}
