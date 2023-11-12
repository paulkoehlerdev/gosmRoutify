package temporaryOSMDataRepository

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/temporaryOSMDataRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/osmpbfData"
	"io"
	"os"
)

type temporaryOSMDataRepositoryImpl struct {
	filePath      string
	logger        logging.Logger
	wayFilePaths  []string
	nodeFilePaths []string

	currentTemporaryWayFile temporaryWayFile
}

const temporaryWayFileMaxSize = 10000

type temporaryWayFile []*osmpbfData.Way

type temporaryNodeFile map[int64]*osmpbfData.Node

func New(filePath string, logger logging.Logger) temporaryOSMDataRepository.TemporaryOSMDataRepository {
	gob.Register(&osmpbfData.Node{})
	gob.Register(&osmpbfData.Way{})

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.MkdirAll(filePath, 0755)
		if err != nil {
			logger.Error().Msgf("error while creating folder %s: %s", filePath, err.Error())
		}
	}

	return &temporaryOSMDataRepositoryImpl{
		filePath: filePath,
		logger:   logger,
	}
}

func (i *temporaryOSMDataRepositoryImpl) AddOSMData(data any) error {
	switch data.(type) {
	case *osmpbfData.Node:
		i.addNode(data.(*osmpbfData.Node))
		return nil
	case *osmpbfData.Way:
		i.addWay(data.(*osmpbfData.Way))
		return nil
	default:
		return fmt.Errorf("unknown data type %T", data)
	}
}

func (i *temporaryOSMDataRepositoryImpl) addNode(node *osmpbfData.Node) {
	path := i.buildTemporaryNodeFilePath(node.ID)
	i.nodeFilePaths = append(i.nodeFilePaths, path)

	nodeFile, err := i.loadTemporaryNodeFile(path)
	if err != nil {
		nodeFile = make(temporaryNodeFile)
	}

	nodeFile[node.ID] = node

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		i.logger.Error().Msgf("error while opening file %s: %s", path, err.Error())
		return
	}

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(nodeFile); err != nil {
		i.logger.Error().Msgf("error while encoding file %s: %s", path, err.Error())
		return
	}
}

func (i *temporaryOSMDataRepositoryImpl) loadTemporaryNodeFile(path string) (temporaryNodeFile, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error while reading file %s: %w", path, err)
	}

	out := make(temporaryNodeFile)
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&out); err != nil {
		return nil, fmt.Errorf("error while decoding tile: %w", err)
	}

	return out, nil
}

func (i *temporaryOSMDataRepositoryImpl) loadNextTemporaryWayFile() error {
	path := i.wayFilePaths[0]
	i.wayFilePaths = i.wayFilePaths[1:]

	if len(i.wayFilePaths) == 0 {
		return io.EOF
	}

	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("error while reading file %s: %w", path, err)
	}

	var out temporaryWayFile
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&out); err != nil {
		return fmt.Errorf("error while decoding tile: %w", err)
	}

	i.currentTemporaryWayFile = out
	return nil
}

func (i *temporaryOSMDataRepositoryImpl) buildTemporaryNodeFilePath(ID int64) string {
	index := ID >> 16
	return fmt.Sprintf("%s/node_%d.tmp", i.filePath, index)
}

func (i *temporaryOSMDataRepositoryImpl) addWay(way *osmpbfData.Way) {
	i.currentTemporaryWayFile = append(i.currentTemporaryWayFile, way)

	if len(i.currentTemporaryWayFile) > temporaryWayFileMaxSize {
		i.saveTemporaryWayFile()
	}
}

func (i *temporaryOSMDataRepositoryImpl) saveTemporaryWayFile() {
	path := i.buildTemporaryWayFilePath(len(i.wayFilePaths))
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		i.logger.Error().Msgf("error while opening file %s: %s", path, err.Error())
		return
	}

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(i.currentTemporaryWayFile); err != nil {
		i.logger.Error().Msgf("error while encoding file %s: %s", path, err.Error())
		return
	}

	i.wayFilePaths = append(i.wayFilePaths, path)
	i.currentTemporaryWayFile = nil
}

func (i *temporaryOSMDataRepositoryImpl) buildTemporaryWayFilePath(index int) string {
	return fmt.Sprintf("%s/way_%d.tmp", i.filePath, index)
}

func (i *temporaryOSMDataRepositoryImpl) Next() (*osmpbfData.Way, error) {
	if len(i.currentTemporaryWayFile) == 0 {
		err := i.loadNextTemporaryWayFile()
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("error while loading next temporary way file: %w", err)
		}
	}

	if len(i.currentTemporaryWayFile) == 0 {
		return nil, io.EOF
	}

	out := i.currentTemporaryWayFile[0]
	i.currentTemporaryWayFile = i.currentTemporaryWayFile[1:]
	return out, nil
}

func (i *temporaryOSMDataRepositoryImpl) FindNode(osmID int64) *osmpbfData.Node {
	path := i.buildTemporaryNodeFilePath(osmID)
	nodeFile, err := i.loadTemporaryNodeFile(path)
	if err != nil {
		return nil
	}

	node, ok := nodeFile[osmID]
	if !ok {
		return nil
	}

	return node
}

func (i *temporaryOSMDataRepositoryImpl) Cleanup() {
	err := os.RemoveAll(fmt.Sprintf("%s/*", i.filePath))
	if err != nil {
		i.logger.Error().Msgf("error while removing folder %s: %s", i.filePath, err.Error())
	}
}
