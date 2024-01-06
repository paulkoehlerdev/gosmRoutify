package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	LoggerConfig   *LoggerConfig   `json:"logging"`
	ImporterConfig *ImporterConfig `json:"import"`
	GraphConfig    *GraphConfig    `json:"graph"`
	ServerConfig   *ServerConfig   `json:"server"`
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
