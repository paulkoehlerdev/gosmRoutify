package wayService

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/way"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/wayRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"sync"
)

type WayService interface {
	InsertWay(way way.Way) error

	InsertWayBulk(way way.Way) error
	CommitBulkInsert() error

	CreateIndices() error

	SelectWayIDsFromNode(nodeID int64) ([]int64, error)

	UpdateCrossings() error
}

const bulkInsertBufferSize = 1<<16 - 1

type impl struct {
	wayRepository       wayRepository.WayRepository
	bulkInsertBuffer    []way.Way
	bulkInsertWaitGroup *sync.WaitGroup
	logger              logging.Logger
}

func New(wayRepository wayRepository.WayRepository, logger logging.Logger) WayService {
	return &impl{
		wayRepository:       wayRepository,
		bulkInsertWaitGroup: &sync.WaitGroup{},
		logger:              logger,
	}
}

func (i *impl) InsertWay(way way.Way) error {
	return i.wayRepository.InsertWay(way)
}

func (i *impl) InsertWayBulk(w way.Way) error {
	if len(i.bulkInsertBuffer) == bulkInsertBufferSize {
		err := i.CommitBulkInsert()
		if err != nil {
			return err
		}
	}

	i.bulkInsertBuffer = append(i.bulkInsertBuffer, w)
	return nil
}

func (i *impl) CommitBulkInsert() error {
	err := i.wayRepository.InsertWays(i.bulkInsertBuffer)
	if err != nil {
		return err
	}
	i.bulkInsertBuffer = make([]way.Way, 0, bulkInsertBufferSize)
	return nil
}

func (i *impl) CreateIndices() error {
	return i.wayRepository.InitIndices()
}

func (i *impl) UpdateCrossings() error {
	return i.wayRepository.UpdateCrossings()
}

func (i *impl) SelectWayIDsFromNode(nodeID int64) ([]int64, error) {
	return i.wayRepository.SelectWayIDsFromNode(nodeID)
}
