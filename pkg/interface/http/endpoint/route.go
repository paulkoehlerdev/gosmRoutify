package endpoint

import (
	"encoding/json"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/application/router"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geojson"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"net/http"
)

type RouteEndpointHandler struct {
	application router.Application
	logger      logging.Logger
}

func NewRouteEndpointHandler(application router.Application, logger logging.Logger) http.Handler {
	return RouteEndpointHandler{
		application: application,
		logger:      logger,
	}
}

func (r RouteEndpointHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	startStr := request.URL.Query().Get("start")
	endStr := request.URL.Query().Get("end")

	var start geojson.Point
	var end geojson.Point

	err := json.Unmarshal([]byte(startStr), &start)
	if err != nil {
		r.logger.Debug().Msgf("error while unmarshalling start coordinate: %s", err.Error())
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal([]byte(endStr), &end)
	if err != nil {
		r.logger.Debug().Msgf("error while unmarshalling end coordinate: %s", err.Error())
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	route, err := r.application.FindRoute(start, end)
	if err != nil {
		r.logger.Debug().Msgf("error while finding route: %s", err.Error())
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	gJson := geojson.NewEmptyGeoJson()
	lineString := geojson.LineString(route)
	gJson.AddFeature(geojson.NewFeature(lineString.ToGeometry()))

	geoJSONBytes, err := json.Marshal(gJson)
	if err != nil {
		r.logger.Debug().Msgf("error while marshalling geojson: %s", err.Error())
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/geo+json")
	_, err = writer.Write(geoJSONBytes)
	if err != nil {
		r.logger.Debug().Msgf("error while writing response: %s", err.Error())
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
