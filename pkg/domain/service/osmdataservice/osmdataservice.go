package osmdataservice

import (
	"errors"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmdatarepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"io"
)

type OsmDataService interface {
	Process(processor osmdatarepository.OsmDataProcessor) error
}

type impl struct {
	logger            logging.Logger
	osmDataRepository osmdatarepository.OsmDataRepository
	filePaths         []string
}

func New(osmDataRepository osmdatarepository.OsmDataRepository, filePaths []string, logger logging.Logger) OsmDataService {
	return &impl{
		logger:            logger,
		osmDataRepository: osmDataRepository,
		filePaths:         filePaths,
	}
}

func (i *impl) Process(processor osmdatarepository.OsmDataProcessor) error {
	for _, filePath := range i.filePaths {
		err := i.osmDataRepository.Process(filePath, processor)
		if errors.Is(err, io.EOF) {
			continue
		}

		if err != nil {
			return fmt.Errorf("error while processing file: %s", err.Error())
		}
	}

	return nil
}
