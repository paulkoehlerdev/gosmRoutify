package http

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/application/router"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/config"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/interface/http/frontend"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"net/http"
)

func NewHttpServer(
	logger logging.Logger,
	application router.Application,
	serverConfig *config.ServerConfig,
) (*http.Server, error) {

	frontendHandler, err := frontend.New(application, logger)
	if err != nil {
		return nil, fmt.Errorf("error while initializing frontend: %s", err.Error())
	}

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
		Handler: frontendHandler.GetMux(),
	}, nil
}
