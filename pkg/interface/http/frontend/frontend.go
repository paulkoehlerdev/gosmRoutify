package frontend

import (
	"encoding/json"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/frontend"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/application/router"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geojson"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"html/template"
	"net/http"
	"strconv"
)

type FrontendHandler interface {
	GetMux() *http.ServeMux
	Init() (FrontendHandler, error)
}

type impl struct {
	mux         *http.ServeMux
	application router.Application
	template    template.Template
	logger      logging.Logger
}

func New(application router.Application, logger logging.Logger) (FrontendHandler, error) {
	return (&impl{
		application: application,
		logger:      logger,
	}).Init()
}

func (i *impl) Init() (FrontendHandler, error) {
	err := i.initTemplates()
	if err != nil {
		return nil, fmt.Errorf("error while initializing templates: %s", err.Error())
	}

	i.initRoutes()

	return i, nil
}

func (i *impl) initRoutes() {
	i.mux = http.NewServeMux()

	i.mux.HandleFunc("/detail", i.Detail)
	i.mux.HandleFunc("/search", i.Search)
	i.mux.HandleFunc("/", i.Root)
}

func (i *impl) initTemplates() error {
	templateFS, err := frontend.GetTemplatesFS()
	if err != nil {
		return fmt.Errorf("error while getting templatesFS: %s", err.Error())
	}

	template, err := template.ParseFS(templateFS, "*.html")
	if err != nil {
		return fmt.Errorf("error while parsing templates: %s", err.Error())
	}

	i.template = *template

	return nil
}

func (i *impl) GetMux() *http.ServeMux {
	return i.mux
}

func (i *impl) Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := i.template.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (i *impl) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	query := r.URL.Query().Get("q")

	adresses, err := i.application.FindAddresses(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Results": adresses,
		"Query":   query,
	}

	if !isHxRequest(r) {
		err = i.template.ExecuteTemplate(w, "base", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = i.template.ExecuteTemplate(w, "results", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (i *impl) Detail(w http.ResponseWriter, r *http.Request) {
	osmID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isHxRequest(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	address, err := i.application.FindAddressByID(osmID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"focusCoordinate": geojson.NewPoint(address.Lat, address.Lon),
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", string(dataBytes))
	w.WriteHeader(http.StatusOK)
}

func isHxRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
