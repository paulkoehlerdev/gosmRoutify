package endpoint

import (
	"github.com/paulkoehlerdev/gosmRoutify/frontend"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"net/http"
)

func NewFrontendEndpointHandler(logger logging.Logger, prefix string) http.Handler {
	fs, err := frontend.GetFrontendFS()
	if err != nil {
		logger.Error().Msgf("Unable to load static filesystem. Won't be serving static files: %s", err.Error())
	}

	return http.StripPrefix(
		prefix,
		http.FileServer(http.FS(fs)),
	)
}
