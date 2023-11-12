package graphRepository

import (
	"encoding/base32"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/graphRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/graph"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/osmpbfData"
	"os"
)

type graphRepositoryImpl struct {
	graphFolder string
	logger      logging.Logger
}

func New(graphFolder string, logger logging.Logger) graphRepository.GraphRepository {
	gob.Register(&osmpbfData.Node{})
	gob.Register(&osmpbfData.Way{})

	out := &graphRepositoryImpl{
		logger:      logger,
		graphFolder: graphFolder,
	}

	for _, level := range graph.GetAllTileLevels() {
		out.createLevelFolder(level)
	}

	return out
}

func (g *graphRepositoryImpl) AddWay(way *osmpbfData.Way, tID graph.TileID) graph.GraphID {
	highwayClass, ok := way.Tags["highway"]
	if !ok {
		g.logger.Debug().Msgf("way %d has no highway tag", way.ID)
		highwayClass = "unclassified"
	}
	level := graph.LevelIDFromHighwayClass(highwayClass)
	t := g.getTile(tID, level)

	gID := t.addWay(tID, level, way)

	err := g.updateTile(tID, level, t)
	if err != nil {
		g.logger.Error().Msgf("error while updating tile: %s", err.Error())
	}

	return gID
}

func (g *graphRepositoryImpl) AddIntersection(node *osmpbfData.Node, a graph.GraphID, b graph.GraphID) graph.GraphID {
	tileA := g.getTile(a.TileID(), a.Level())
	tileB := g.getTile(b.TileID(), b.Level())

	if graph.TileIDFromNode(node, a.Level()) != a.TileID() && graph.TileIDFromNode(node, b.Level()) != b.TileID() {
		g.logger.Error().Msgf("node %d is not in tile %d or %d", node.ID, a.TileID(), b.TileID())
	}

	var gID graph.GraphID

	if graph.TileIDFromNode(node, a.Level()) != a.TileID() {
		gID := tileA.addNode(a.TileID(), a.Level(), node)
		tileA.addRelation(a, b, gID)
		tileB.addRelation(b, a, gID)
	} else {
		gID := tileB.addNode(b.TileID(), b.Level(), node)
		tileB.addRelation(b, a, gID)
		tileA.addRelation(a, b, gID)
	}

	err := g.updateTile(a.TileID(), a.Level(), tileA)
	if err != nil {
		g.logger.Error().Msgf("error while updating tile: %s", err.Error())
	}

	err = g.updateTile(b.TileID(), b.Level(), tileB)
	if err != nil {
		g.logger.Error().Msgf("error while updating tile: %s", err.Error())
	}

	return gID
}

func (g *graphRepositoryImpl) buildFilePath(id graph.TileID, level graph.LevelID) string {
	// encode tile id in base 32
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(id))
	b64id := base32.HexEncoding.EncodeToString(bytes)
	return fmt.Sprintf("%s/%d/%s.graph", g.graphFolder, level, b64id)
}

func (g *graphRepositoryImpl) createLevelFolder(level graph.LevelID) {
	path := fmt.Sprintf("%s/%d", g.graphFolder, level)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			g.logger.Error().Msgf("error while creating folder %s: %s", path, err.Error())
		}
	}
}

func (g *graphRepositoryImpl) getTileFromDisk(path string) (*tile, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error while reading file %s: %w", path, err)
	}

	out := new(tile)
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&out); err != nil {
		return nil, fmt.Errorf("error while decoding tile: %w", err)
	}

	return out, nil
}

func (g *graphRepositoryImpl) saveTileToDisk(path string, tile *tile) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error while opening file %s: %w", path, err)
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	if err := enc.Encode(tile); err != nil {
		return fmt.Errorf("error while encoding tile: %w", err)
	}

	return nil
}

func (g *graphRepositoryImpl) updateTile(id graph.TileID, level graph.LevelID, in *tile) error {
	tileFileName := g.buildFilePath(id, level)
	err := g.saveTileToDisk(tileFileName, in)
	if err != nil {
		return fmt.Errorf("error while saving tile to disk: %w", err)
	}
	return nil
}

func (g *graphRepositoryImpl) getTile(id graph.TileID, level graph.LevelID) *tile {
	tileFileName := g.buildFilePath(id, level)
	out, err := g.getTileFromDisk(tileFileName)
	if err != nil {
		out = new(tile)
	}

	return out
}
