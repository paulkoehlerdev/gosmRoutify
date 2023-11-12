package config

type graphRepositoryConfig struct {
	Path           string `json:"path"`
	MaxFileHandles int    `json:"maxFileHandles"`
}
