package nodeservice

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/nodetype"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/noderepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinate"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
)

type NodeService interface {
	Add(osmID int64, nodeType nodetype.NodeType)
	AddOrUpdate(osmID int64, nodeType nodetype.NodeType, updateFunc func(nodeType nodetype.NodeType) nodetype.NodeType)
	GetNodeType(osmID int64) nodetype.NodeType
	SetCoordinate(osmID int64, coordinate coordinate.Coordinate) nodetype.NodeType
	GetCoordinate(osmID int64) (coordinate.Coordinate, error)
	SetTags(osmID int64, tags map[string]string) error
	GetTags(osmID int64) (map[string]string, error)
	SetSplitNode(osmID int64)
	IsSplitNode(osmID int64) bool
	UnsetSplitNode(osmID int64)
	PrintNodeTypeStatistics()
}

type impl struct {
	nodeRepository noderepository.NodeRepository
	logger         logging.Logger
}

func New(nodeRepository noderepository.NodeRepository, logger logging.Logger) NodeService {
	return &impl{
		nodeRepository: nodeRepository,
		logger:         logger,
	}
}

func (i *impl) Add(osmID int64, nodeType nodetype.NodeType) {
	i.nodeRepository.Add(osmID, nodeType)
}

func (i *impl) AddOrUpdate(osmID int64, nodeType nodetype.NodeType, updateFunc func(nodeType nodetype.NodeType) nodetype.NodeType) {
	i.nodeRepository.AddOrUpdate(osmID, nodeType, updateFunc)
}

func (i *impl) GetNodeType(osmID int64) nodetype.NodeType {
	return i.nodeRepository.GetNodeType(osmID)
}

func (i *impl) SetCoordinate(osmID int64, coordinate coordinate.Coordinate) nodetype.NodeType {
	return i.nodeRepository.SetCoordinate(osmID, coordinate)
}

func (i *impl) GetCoordinate(osmID int64) (coordinate.Coordinate, error) {
	return i.nodeRepository.GetCoordinate(osmID)
}

func (i *impl) SetTags(osmID int64, tags map[string]string) error {
	return i.nodeRepository.SetTags(osmID, tags)
}

func (i *impl) GetTags(osmID int64) (map[string]string, error) {
	return i.nodeRepository.GetTags(osmID)
}

func (i *impl) SetSplitNode(osmID int64) {
	i.nodeRepository.SetSplitNode(osmID)
}

func (i *impl) IsSplitNode(osmID int64) bool {
	return i.nodeRepository.IsSplitNode(osmID)
}

func (i *impl) UnsetSplitNode(osmID int64) {
	i.nodeRepository.UnsetSplitNode(osmID)
}

func (i *impl) PrintNodeTypeStatistics() {
	nodeTypeStatistics := i.nodeRepository.CalcNodeTypeStatistics()
	i.logger.Info().Msgf("Node type statistics: %v", nodeTypeStatistics)
}
