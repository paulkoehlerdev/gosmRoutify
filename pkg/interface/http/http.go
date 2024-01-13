package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/application/router"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/config"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geojson"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"net/http"
	"strconv"
)

type impl struct {
	logger      logging.Logger
	application router.Application
}

func NewHttpServer(
	logger logging.Logger,
	application router.Application,
	serverConfig *config.ServerConfig,
) (*http.Server, error) {
	mux := &http.ServeMux{}

	server := &impl{
		logger:      logger,
		application: application,
	}

	mux.HandleFunc("/api/route", server.route)
	mux.HandleFunc("/api/locate", server.locate)
	mux.HandleFunc("/api/search", server.search)

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
		Handler: mux,
	}, nil
}

func (i *impl) route(w http.ResponseWriter, r *http.Request) {
	routeQueryB64 := r.URL.Query().Get("r")

	routeQuery, err := base64.URLEncoding.DecodeString(routeQueryB64)
	if err != nil {
		http.Error(w, "invalid route query", http.StatusBadRequest)
		return
	}

	var points []geojson.Point
	err = json.Unmarshal(routeQuery, &points)
	if err != nil {
		http.Error(w, "invalid route query", http.StatusBadRequest)
		return
	}

	if len(points) < 2 {
		http.Error(w, "invalid route query", http.StatusBadRequest)
		return
	}

	route, err := i.application.FindRoute(points)
	if err != nil {
		i.logger.Error().Msgf("error while finding route: %s", err.Error())
		http.Error(w, fmt.Sprintf("error while finding route: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	routeBytes, err := json.Marshal(route)
	if err != nil {
		i.logger.Error().Msgf("error while marshalling route: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(routeBytes)
	if err != nil {
		i.logger.Error().Msgf("error while writing response: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (i *impl) locate(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	position, err := i.application.LocateAddressByID(idInt)
	if err != nil {
		i.logger.Error().Msgf("error while locating address: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	positionBytes, err := json.Marshal(position)
	if err != nil {
		i.logger.Error().Msgf("error while marshalling position: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(positionBytes)
	if err != nil {
		i.logger.Error().Msgf("error while writing response: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (i *impl) search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	addresses, err := i.application.FindAddresses(query)
	if err != nil {
		i.logger.Error().Msgf("error while searching for addresses: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	addressBytes, err := json.Marshal(addresses)
	if err != nil {
		i.logger.Error().Msgf("error while marshalling addresses: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(addressBytes)
	if err != nil {
		i.logger.Error().Msgf("error while writing response: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
