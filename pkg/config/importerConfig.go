package config

type ImporterConfig struct {
	Enable            bool   `json:"enable"`
	FilePath          string `json:"path"`
	TmpFilePath       string `json:"tmpPath"`
	EnableDiskStorage bool   `json:"enableDiskStorage"`
}
