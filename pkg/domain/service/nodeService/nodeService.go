package nodeService

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/node"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/nodeRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
)

const bulkInsertBufferSize = 1<<20 - 1

type NodeService interface {
	InsertNode(node node.Node) error

	InsertNodeBulk(node node.Node) error
	CommitBulkInsert() error

	CreateIndices() error

	SelectNodeFromID(id int64) (*node.Node, error)
	SelectNodesFromIDs(ids []int64) ([]*node.Node, error)

	LocateOsmID(osmID int64) (lat, lon float64, err error)
}

type impl struct {
	logger           logging.Logger
	nodeRepository   nodeRepository.NodeRepository
	bulkInsertBuffer []node.Node
}

func New(nodeRepository nodeRepository.NodeRepository, logger logging.Logger) NodeService {
	return &impl{
		nodeRepository: nodeRepository,
		logger:         logger,
	}
}

func (i *impl) InsertNode(node node.Node) error {
	return i.nodeRepository.InsertNode(node)
}

func (i *impl) InsertNodeBulk(n node.Node) error {
	if len(i.bulkInsertBuffer) == bulkInsertBufferSize {
		err := i.CommitBulkInsert()
		if err != nil {
			return fmt.Errorf("error while committing bulk insert: %s", err.Error())
		}
	}

	i.bulkInsertBuffer = append(i.bulkInsertBuffer, n)
	return nil
}

func (i *impl) CommitBulkInsert() error {
	err := i.nodeRepository.InsertNodes(i.bulkInsertBuffer)
	if err != nil {
		return fmt.Errorf("error while inserting nodes: %s", err.Error())
	}
	i.bulkInsertBuffer = make([]node.Node, 0, bulkInsertBufferSize)
	return nil
}

func (i *impl) CreateIndices() error {
	return i.nodeRepository.InitIndices()
}

func (i *impl) SelectNodeFromID(id int64) (*node.Node, error) {
	return i.nodeRepository.SelectNodeFromID(id)
}

func (i *impl) SelectNodesFromIDs(ids []int64) ([]*node.Node, error) {
	var out []*node.Node
	for _, id := range ids {
		node, err := i.SelectNodeFromID(id)
		if err != nil {
			return nil, fmt.Errorf("error while selecting node from id: %s", err.Error())
		}

		out = append(out, node)
	}
	return out, nil
}

func (i *impl) LocateOsmID(osmID int64) (lat, lon float64, err error) {
	node, err := i.nodeRepository.SelectNodeFromID(osmID)
	if err == nil {
		return node.Lat, node.Lon, nil
	}

	lat, lon, err = i.nodeRepository.SelectCenterOfWayID(osmID)
	if err != nil {
		return 0, 0, fmt.Errorf("error while locating osmID: %s", err.Error())
	}

	return lat, lon, nil
}
