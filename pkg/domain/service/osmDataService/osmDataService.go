package osmDataService

import (
	"errors"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/osmDataRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"io"
)

type OsmDataService interface {
	Next() (any, error)
}

type impl struct {
	logger            logging.Logger
	osmDataRepository osmDataRepository.OsmDataRepository
	filePaths         []string
	index             int
}

func New(osmDataRepository osmDataRepository.OsmDataRepository, filePaths []string, logger logging.Logger) (OsmDataService, error) {
	service := &impl{
		logger:            logger,
		osmDataRepository: osmDataRepository,
		filePaths:         filePaths,
		index:             -1,
	}

	err := service.selectNextFile()
	if err != nil {
		return nil, fmt.Errorf("error while selecting a siutable file: %s", err.Error())
	}

	return service, nil
}

func (i *impl) selectNextFile() error {
	for {
		i.index++
		if i.index >= len(i.filePaths) {
			return io.EOF
		}

		err := i.osmDataRepository.Read(i.filePaths[i.index])
		if err != nil {
			i.logger.Error().Msgf("error while reading data: %s", err.Error())
			continue
		}

		return nil
	}
}

func (i *impl) Next() (any, error) {
	data, err := i.osmDataRepository.Next()
	if errors.Is(err, io.EOF) {
		if err := i.selectNextFile(); errors.Is(err, io.EOF) {
			return nil, io.EOF
		}

		return i.Next()
	}

	if err != nil {
		i.logger.Error().Msgf("error while reading data: %s", err.Error())
		return i.Next()
	}

	return data, nil
}
