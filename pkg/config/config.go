package config

import (
	"encoding/json"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"os"
)

type Config struct {
	LoggerConfig   *LoggerConfig   `json:"logging"`
	DatabaseConfig *DatabaseConfig `json:"database"`
	ServerConfig   *ServerConfig   `json:"server"`
}

type LoggerConfig struct {
	Level    string `json:"level"`
	FilePath string `json:"file"`
	Console  bool   `json:"console"`
}

type DatabaseConfig struct {
	FilePath string `json:"file"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func FromFile(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error while reading file: %s", err.Error())
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling json config: %s", err.Error())
	}

	return &config, nil
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
