package wayService

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/way"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/wayRepository"
)

type WayService interface {
	InsertWay(way way.Way) error

	SelectWayIDsFromNode(nodeID int64) ([]int64, error)
}

type impl struct {
	wayRepository wayRepository.WayRepository
}

func New(wayRepository wayRepository.WayRepository) WayService {
	return &impl{
		wayRepository: wayRepository,
	}
}

func (i *impl) InsertWay(way way.Way) error {
	return i.wayRepository.InsertWay(way)
}

func (i *impl) SelectWayIDsFromNode(nodeID int64) ([]int64, error) {
	return i.wayRepository.SelectWayIDsFromNode(nodeID)
}
