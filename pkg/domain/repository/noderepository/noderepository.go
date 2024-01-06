package noderepository

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/nodetype"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinate"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/osmid"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/kvstorage"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"math"
)

const minNodeCapacity = 1 << 16

type NodeRepository interface {
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
	CalcNodeTypeStatistics() map[nodetype.NodeType]int
	Close() error
}

type impl struct {
	logger         logging.Logger
	nodeTypes      map[osmid.OsmID]nodetype.NodeType
	database       kvstorage.KVStorage[int64, []byte]
	nodeTags       kvstorage.Collection[int64, []byte]
	coordinates    kvstorage.Collection[int64, []byte]
	nodesToBeSplit map[osmid.OsmID]struct{}
}

func New(dbPath string, enableDiskStorage bool, logger logging.Logger) (NodeRepository, error) {
	kvStorageConstructor := kvstorage.New[int64, []byte]
	if !enableDiskStorage {
		kvStorageConstructor = kvstorage.NewRam[int64, []byte]
	}

	database, err := kvStorageConstructor(dbPath, kvstorage.DefaultKVStorageOptions())
	if err != nil {
		return nil, fmt.Errorf("error while creating kvstorage: %s", err.Error())
	}

	nodeTagsCollection, err := database.GetCollection("nodeTags")
	if err != nil {
		nodeTagsCollection, err = database.NewCollection("nodeTags")
		if err != nil {
			return nil, fmt.Errorf("error while creating collection: %s", err.Error())
		}
	}

	coordinatesCollection, err := database.GetCollection("coordinates")
	if err != nil {
		coordinatesCollection, err = database.NewCollection("coordinates")
		if err != nil {
			return nil, fmt.Errorf("error while creating collection: %s", err.Error())
		}
	}

	repo := &impl{
		logger:         logger,
		nodeTypes:      make(map[osmid.OsmID]nodetype.NodeType),
		database:       database,
		nodeTags:       nodeTagsCollection,
		coordinates:    coordinatesCollection,
		nodesToBeSplit: make(map[osmid.OsmID]struct{}),
	}

	return repo, nil
}

func (i *impl) Add(osmID osmid.OsmID, nodeType nodetype.NodeType) {
	i.nodeTypes[osmID] = nodeType
}

func (i *impl) AddOrUpdate(osmID osmid.OsmID, newNodeType nodetype.NodeType, updateFunc func(nodeType nodetype.NodeType) nodetype.NodeType) {
	nodeType, ok := i.nodeTypes[osmID]
	if !ok {
		i.Add(osmID, newNodeType)
		return
	}

	if nodeType == nodetype.EMPTYNODE {
		i.Add(osmID, newNodeType)
		return
	}

	i.Add(osmID, updateFunc(nodeType))
}

func (i *impl) GetNodeType(osmID osmid.OsmID) nodetype.NodeType {
	nodeType, ok := i.nodeTypes[osmID]
	if !ok {
		return nodetype.EMPTYNODE
	}
	return nodeType
}

func (i *impl) SetCoordinate(osmID osmid.OsmID, coordinate coordinate.Coordinate) nodetype.NodeType {
	nodeType, ok := i.nodeTypes[osmID]
	if !ok {
		return nodetype.EMPTYNODE
	}

	buf := make([]byte, 16)
	binary.LittleEndian.PutUint64(buf[0:8], math.Float64bits(coordinate.Lat()))
	binary.LittleEndian.PutUint64(buf[8:16], math.Float64bits(coordinate.Lon()))

	err := i.coordinates.Set(int64(osmID), buf)
	if err != nil {
		i.logger.Error().Msgf("error while setting coordinate: %s", err.Error())
	}

	return nodeType
}

func (i *impl) GetCoordinate(osmID osmid.OsmID) (coordinate.Coordinate, error) {
	_, ok := i.nodeTypes[osmID]
	if !ok {
		return coordinate.Coordinate{}, fmt.Errorf("error while getting coordinate: node not found")
	}

	bufptr, err := i.coordinates.Get(int64(osmID))
	if err != nil {
		return coordinate.Coordinate{}, fmt.Errorf("error while getting coordinate: %s", err.Error())
	}

	if bufptr == nil {
		return coordinate.Coordinate{}, fmt.Errorf("error while getting coordinate: coordinate not found")
	}
	buf := *bufptr

	lat := binary.LittleEndian.Uint64(buf[0:8])
	lon := binary.LittleEndian.Uint64(buf[8:16])

	coo := coordinate.New(math.Float64frombits(lat), math.Float64frombits(lon))

	return coo, nil
}

func (i *impl) SetTags(osmID osmid.OsmID, tags map[string]string) error {

	nodeType, ok := i.nodeTypes[osmID]
	if !ok {
		return fmt.Errorf("error while setting tags: node not found")
	}

	if nodeType == nodetype.EMPTYNODE {
		return fmt.Errorf("error while setting tags: node not found")
	}

	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(tags)
	if err != nil {
		return fmt.Errorf("error while parsing tags: %s", err.Error())
	}

	err = i.nodeTags.Set(int64(osmID), buf.Bytes())
	if err != nil {
		i.logger.Error().Msgf("error while setting tags: %s", err.Error())
	}

	return nil
}

func (i *impl) GetTags(osmID osmid.OsmID) (map[string]string, error) {
	nodeType, ok := i.nodeTypes[osmID]
	if !ok {
		return nil, fmt.Errorf("error while getting tags: node not found")
	}

	if nodeType == nodetype.EMPTYNODE {
		return nil, fmt.Errorf("error while getting tags: node not found")
	}

	bufptr, err := i.nodeTags.Get(int64(osmID))
	if err != nil {
		return nil, fmt.Errorf("error while getting tags: %s", err.Error())
	}

	if bufptr == nil {
		return nil, fmt.Errorf("error while getting tags: tagsChan were nil")
	}

	tags := make(map[string]string)

	buf := *bufptr
	err = gob.NewDecoder(bytes.NewReader(buf)).Decode(&tags)
	if err != nil {
		return nil, fmt.Errorf("error while parsing tagsChan: %s", err.Error())
	}

	return tags, nil
}

func (i *impl) SetSplitNode(osmID osmid.OsmID) {
	i.nodesToBeSplit[osmID] = struct{}{}
}

func (i *impl) IsSplitNode(osmID osmid.OsmID) bool {
	_, ok := i.nodesToBeSplit[osmID]
	return ok
}

func (i *impl) UnsetSplitNode(osmID osmid.OsmID) {
	delete(i.nodesToBeSplit, osmID)
}

func (i *impl) CalcNodeTypeStatistics() map[nodetype.NodeType]int {
	statistics := make(map[nodetype.NodeType]int)
	for _, nodeType := range i.nodeTypes {
		if _, ok := statistics[nodeType]; !ok {
			statistics[nodeType] = 0
		}
		statistics[nodeType]++
	}
	return statistics
}

func (i *impl) Close() error {
	err := i.database.Close()
	if err != nil {
		return fmt.Errorf("error while closing database: %s", err.Error())
	}

	return nil
}
