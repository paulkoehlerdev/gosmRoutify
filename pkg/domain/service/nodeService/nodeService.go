package nodeService

import (
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/node"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/nodeRepository"
)

type NodeService interface {
	InsertNode(node node.Node) error
}

type impl struct {
	nodeRepository nodeRepository.NodeRepository
}

func New(nodeRepository nodeRepository.NodeRepository) NodeService {
	return &impl{
		nodeRepository: nodeRepository,
	}
}

func (i *impl) InsertNode(node node.Node) error {
	return i.nodeRepository.InsertNode(node)
}
