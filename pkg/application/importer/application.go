package importer

import (
	"errors"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/service/osmDataService"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"io"
)

const counterLogBreak = 1000000

type Importer interface {
	StartDataImport() error
}

type impl struct {
	osmDataService osmDataService.OsmDataService
	logger         logging.Logger
}

func New(osmDataService osmDataService.OsmDataService, logger logging.Logger) Importer {
	return &impl{
		osmDataService: osmDataService,
		logger:         logger,
	}
}

func (i *impl) StartDataImport() error {
	counter := 0
	for {
		counter++
		if counter%counterLogBreak == 0 {
			i.logger.Info().Msgf("read %d Mio. elements", counter/counterLogBreak)
		}

		_, err := i.osmDataService.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return err
			}
			return fmt.Errorf("error while reading data: %s", err.Error())
		}
	}
}
