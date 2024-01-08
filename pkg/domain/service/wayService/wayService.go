package wayService

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/way"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/wayRepository"
)

type WayService interface {
	InsertWay(way way.Way) error

	InsertWayBulk(way way.Way) error
	CommitBulkInsert() error

	SelectWayIDsFromNode(nodeID int64) ([]int64, error)
}

const bulkInsertBufferSize = 2 << 9

type impl struct {
	wayRepository    wayRepository.WayRepository
	bulkInsertBuffer []way.Way
}

func New(wayRepository wayRepository.WayRepository) WayService {
	return &impl{
		wayRepository: wayRepository,
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

func (i *impl) SelectWayIDsFromNode(nodeID int64) ([]int64, error) {
	return i.wayRepository.SelectWayIDsFromNode(nodeID)
}
