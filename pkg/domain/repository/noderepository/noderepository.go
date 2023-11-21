package noderepository

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/nodetype"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinate"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinatelist"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/kvstorage"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
)

const minNodeCapacity = 1 << 16

type NodeRepository interface {
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
	CalcNodeTypeStatistics() map[nodetype.NodeType]int
}

type impl struct {
	logger         logging.Logger
	nodes          map[int64]nodeIndex
	towerNodes     coordinatelist.CoordinateList
	pillarNodes    coordinatelist.CoordinateList
	nodeTags       kvstorage.KVStorage[nodeIndex, map[string]string]
	nodesToBeSplit map[int64]struct{}
}

func New(logger logging.Logger) NodeRepository {
	return &impl{
		logger:         logger,
		nodes:          make(map[int64]nodeIndex),
		towerNodes:     coordinatelist.NewCoordinateList(minNodeCapacity),
		pillarNodes:    coordinatelist.NewCoordinateList(minNodeCapacity),
		nodeTags:       kvstorage.NewRamKVStorage[nodeIndex, map[string]string](minNodeCapacity),
		nodesToBeSplit: make(map[int64]struct{}),
	}
}

func (i *impl) Add(osmID int64, nodeType nodetype.NodeType) {
	index, err := i.indexFromOsmID(osmID)
	if err == nil {
		index, err = newNodeIndex(nodeType, 0)
		if err != nil {
			panic(fmt.Errorf("error while creating nodeIndex: %s", err.Error()))
		}
	} else {
		index, err = newNodeIndex(nodeType, index.GetNodeID())
		if err != nil {
			panic(fmt.Errorf("error while creating nodeIndex: %s", err.Error()))
		}
	}

	i.nodes[osmID] = index
}

func (i *impl) indexFromOsmID(osmID int64) (nodeIndex, error) {
	nodeIndex, ok := i.nodes[osmID]
	if !ok {
		return 0, fmt.Errorf("node with osmID %d not found", osmID)
	}
	return nodeIndex, nil
}

func (i *impl) AddOrUpdate(osmID int64, newNodeType nodetype.NodeType, updateFunc func(nodeType nodetype.NodeType) nodetype.NodeType) {
	nodeType := nodetype.EMPTYNODE
	index, err := i.indexFromOsmID(osmID)
	if err == nil {
		nodeType = index.GetNodeType()
	}

	if nodeType == nodetype.EMPTYNODE {
		i.Add(osmID, newNodeType)
		return
	}

	i.Add(osmID, updateFunc(nodeType))
}

func (i *impl) GetNodeType(osmID int64) nodetype.NodeType {
	index, err := i.indexFromOsmID(osmID)
	if err != nil {
		return nodetype.EMPTYNODE
	}

	return index.GetNodeType()
}

func (i *impl) SetCoordinate(osmID int64, coordinate coordinate.Coordinate) nodetype.NodeType {
	nodeType := nodetype.EMPTYNODE
	index, err := i.indexFromOsmID(osmID)
	if err == nil {
		nodeType = index.GetNodeType()
	}

	if nodeType.IsTowerNode() {
		index, err = newNodeIndex(nodeType, i.addTowerNode(coordinate))
		if err != nil {
			panic(fmt.Errorf("error while creating nodeIndex: %s", err.Error()))
		}
	} else if nodeType.IsPillarNode() {
		index, err = newNodeIndex(nodeType, i.addPillarNode(coordinate))
		if err != nil {
			panic(fmt.Errorf("error while creating nodeIndex: %s", err.Error()))
		}
	}

	i.nodes[osmID] = index
	return index.GetNodeType()
}

func (i *impl) GetCoordinate(osmID int64) (coordinate.Coordinate, error) {
	index, err := i.indexFromOsmID(osmID)
	if err != nil {
		return coordinate.Coordinate{}, fmt.Errorf("error while getting coordinate: %s", err.Error())
	}

	nodeType := index.GetNodeType()

	if nodeType.IsTowerNode() {
		return i.towerNodes.Get(index.GetNodeID()), nil
	}

	if nodeType.IsPillarNode() {
		return i.pillarNodes.Get(index.GetNodeID()), nil
	}

	return coordinate.Coordinate{}, fmt.Errorf("error while getting coordinate: unknown node type %d", index.GetNodeType())
}

func (i *impl) addTowerNode(coordinate coordinate.Coordinate) uint64 {
	id := i.towerNodes.Len()
	i.towerNodes.Append(coordinate)
	return id
}

func (i *impl) addPillarNode(coordinate coordinate.Coordinate) uint64 {
	id := i.pillarNodes.Len()
	i.pillarNodes.Append(coordinate)
	return id
}

func (i *impl) SetTags(osmID int64, tags map[string]string) error {
	index, err := i.indexFromOsmID(osmID)
	if err != nil {
		return fmt.Errorf("error while setting tags: %s", err.Error())
	}

	if err := i.nodeTags.Set(index, tags); err != nil {
		return fmt.Errorf("error while setting tags: %s", err.Error())
	}

	return nil
}

func (i *impl) GetTags(osmID int64) (map[string]string, error) {
	index, err := i.indexFromOsmID(osmID)
	if err != nil {
		return nil, fmt.Errorf("error while getting tags: %s", err.Error())
	}

	tags, err := i.nodeTags.Get(index)
	if err != nil {
		return nil, fmt.Errorf("error while getting tags: %s", err.Error())
	}

	return tags, nil
}

func (i *impl) SetSplitNode(osmID int64) {
	i.nodesToBeSplit[osmID] = struct{}{}
}

func (i *impl) IsSplitNode(osmID int64) bool {
	_, ok := i.nodesToBeSplit[osmID]
	return ok
}

func (i *impl) UnsetSplitNode(osmID int64) {
	delete(i.nodesToBeSplit, osmID)
}

func (i *impl) CalcNodeTypeStatistics() map[nodetype.NodeType]int {
	statistics := make(map[nodetype.NodeType]int)
	for _, nodeIndex := range i.nodes {
		if _, ok := statistics[nodeIndex.GetNodeType()]; !ok {
			statistics[nodeIndex.GetNodeType()] = 0
		}
		statistics[nodeIndex.GetNodeType()]++
	}
	return statistics
}
