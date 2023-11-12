package importer

import (
	"errors"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/graphService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmDataPreprocessorService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmDataService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"io"
)

const counterLogBreak = 1000000

type Importer interface {
	StartDataImport() error
	BuildGraph()
}

type impl struct {
	osmDataService osmDataService.OsmDataService
	preprocessor   osmDataPreprocessorService.OsmDataPreprocessorService
	graphService   graphService.GraphService
	logger         logging.Logger
}

func New(osmDataService osmDataService.OsmDataService, preprocessor osmDataPreprocessorService.OsmDataPreprocessorService, graphService graphService.GraphService, logger logging.Logger) Importer {
	return &impl{
		osmDataService: osmDataService,
		logger:         logger,
		preprocessor:   preprocessor,
		graphService:   graphService,
	}
}

func (i *impl) StartDataImport() error {
	counter := 0

	for {
		counter++
		if counter%counterLogBreak == 0 {
			i.logger.Info().Msgf("read %d Mio. elements", counter/counterLogBreak)
		}

		data, err := i.osmDataService.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("error while reading data: %s", err.Error())
		}

		if i.preprocessor.Filter(data) {
			continue
		}

		err = i.graphService.AddOSMData(data)
		if err != nil {
			return fmt.Errorf("error while adding data to graph: %s", err.Error())
		}
	}

	return io.EOF
}

func (i *impl) BuildGraph() {
	i.graphService.BuildGraph()
}
