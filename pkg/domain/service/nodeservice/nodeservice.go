package nodeservice

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/nodetype"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/noderepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinate"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/osmid"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
)

type NodeService interface {
	Add(osmID osmid.OsmID, nodeType nodetype.NodeType)
	AddOrUpdate(osmID osmid.OsmID, nodeType nodetype.NodeType, updateFunc func(nodeType nodetype.NodeType) nodetype.NodeType)
	GetNodeType(osmID osmid.OsmID) nodetype.NodeType
	SetCoordinate(osmID osmid.OsmID, coordinate coordinate.Coordinate) nodetype.NodeType
	GetCoordinate(osmID osmid.OsmID) (coordinate.Coordinate, error)
	SetTags(osmID osmid.OsmID, tags map[string]string) error
	GetTags(osmID osmid.OsmID) (map[string]string, error)
	SetSplitNode(osmID osmid.OsmID)
	IsSplitNode(osmID osmid.OsmID) bool
	UnsetSplitNode(osmID osmid.OsmID)
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

func (i *impl) Add(osmID osmid.OsmID, nodeType nodetype.NodeType) {
	i.nodeRepository.Add(osmID, nodeType)
}

func (i *impl) AddOrUpdate(osmID osmid.OsmID, nodeType nodetype.NodeType, updateFunc func(nodeType nodetype.NodeType) nodetype.NodeType) {
	i.nodeRepository.AddOrUpdate(osmID, nodeType, updateFunc)
}

func (i *impl) GetNodeType(osmID osmid.OsmID) nodetype.NodeType {
	return i.nodeRepository.GetNodeType(osmID)
}

func (i *impl) SetCoordinate(osmID osmid.OsmID, coordinate coordinate.Coordinate) nodetype.NodeType {
	return i.nodeRepository.SetCoordinate(osmID, coordinate)
}

func (i *impl) GetCoordinate(osmID osmid.OsmID) (coordinate.Coordinate, error) {
	return i.nodeRepository.GetCoordinate(osmID)
}

func (i *impl) SetTags(osmID osmid.OsmID, tags map[string]string) error {
	return i.nodeRepository.SetTags(osmID, tags)
}

func (i *impl) GetTags(osmID osmid.OsmID) (map[string]string, error) {
	return i.nodeRepository.GetTags(osmID)
}

func (i *impl) SetSplitNode(osmID osmid.OsmID) {
	i.nodeRepository.SetSplitNode(osmID)
}

func (i *impl) IsSplitNode(osmID osmid.OsmID) bool {
	return i.nodeRepository.IsSplitNode(osmID)
}

func (i *impl) UnsetSplitNode(osmID osmid.OsmID) {
	i.nodeRepository.UnsetSplitNode(osmID)
}

func (i *impl) PrintNodeTypeStatistics() {
	nodeTypeStatistics := i.nodeRepository.CalcNodeTypeStatistics()
	i.logger.Info().Msgf("Node type statistics: %v", nodeTypeStatistics)
}
