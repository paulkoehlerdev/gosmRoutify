package frontend

import "embed"

//go:embed src
var frontendFS embed.FS

func GetFrontendFS() embed.FS {
	return frontendFS
}
