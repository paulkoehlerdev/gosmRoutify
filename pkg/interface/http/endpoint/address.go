package endpoint

import (
	"encoding/json"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/application/router"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"net/http"
)

type AddressEndpointHandler struct {
	application router.Application
	logger      logging.Logger
}

func NewAddressEndpointHandler(application router.Application, logger logging.Logger) http.Handler {
	return AddressEndpointHandler{
		application: application,
		logger:      logger,
	}
}

func (r AddressEndpointHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query().Get("q")

	addresses, err := r.application.FindAddresses(query)
	if err != nil {
		r.logger.Debug().Msgf("error while finding addresses: %s", err.Error())
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(addresses)
	if err != nil {
		r.logger.Debug().Msgf("error while encoding addresses: %s", err.Error())
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
