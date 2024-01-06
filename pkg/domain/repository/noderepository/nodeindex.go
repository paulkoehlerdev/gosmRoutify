package noderepository

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/nodetype"
	"math"
)

type nodeIndex uint64

func newNodeIndex(nodeType nodetype.NodeType, nodeID uint64) (nodeIndex, error) {
	if nodeID > math.MaxInt64>>8 {
		return 0, fmt.Errorf("nodeID %d is too large", nodeID)
	}
	return nodeIndex(nodeType)<<56 | nodeIndex(nodeID), nil
}

func (i nodeIndex) GetNodeType() nodetype.NodeType {
	return nodetype.NodeType(i >> 56)
}

func (i nodeIndex) GetNodeID() uint64 {
	return uint64(i) & (math.MaxUint64 >> 8)
}
