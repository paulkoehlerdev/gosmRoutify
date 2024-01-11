package http

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/application/router"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/config"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/interface/http/endpoint"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"net/http"
)

func NewHttpServer(
	logger logging.Logger,
	application router.Application,
	serverConfig *config.ServerConfig,
) (*http.Server, error) {
	serveMux := http.NewServeMux()

	routeEndpoint := endpoint.NewRouteEndpointHandler(application, logger.WithAttrs("endpoint", "route"))
	serveMux.Handle("/api/route", routeEndpoint)

	addressEndpoint := endpoint.NewAddressEndpointHandler(application, logger.WithAttrs("endpoint", "address"))
	serveMux.Handle("/api/address", addressEndpoint)

	frontendEndpoint := endpoint.NewFrontendEndpointHandler(logger.WithAttrs("endpoint", "frontend"), "")
	serveMux.Handle("/", frontendEndpoint)

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
		Handler: serveMux,
	}, nil
}
