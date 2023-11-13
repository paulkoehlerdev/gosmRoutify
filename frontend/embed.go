package frontend

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed src
var frontendFS embed.FS

func GetFrontendFS() (fs.FS, error) {
	fs, err := fs.Sub(frontendFS, "src")
	if err != nil {
		return nil, fmt.Errorf("error reading frontendFS subfolder: %s", err.Error())
	}

	return fs, nil
}
