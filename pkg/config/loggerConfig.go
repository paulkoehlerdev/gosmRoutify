package config

type LoggerConfig struct {
	Level    string `json:"level"`
	FilePath string `json:"file"`
	Console  bool   `json:"console"`
}
