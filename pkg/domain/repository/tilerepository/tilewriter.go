package tilerepository

import (
	"encoding/gob"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/graphtile"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/osmid"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/kdtree"
	"os"
)

type tileWriter interface {
	Write(name string, tile *graphtile.GraphTile) error
}

type tileWriterImpl struct {
	filePath string
}

func newTileWriter(filePath string) tileWriter {
	gob.Register(kdtree.New[osmid.OsmID]())
	return &tileWriterImpl{
		filePath: filePath,
	}
}

func (g *tileWriterImpl) Write(name string, tile *graphtile.GraphTile) error {
	file, err := os.OpenFile(
		fmt.Sprintf("%s/%s%s", g.filePath, name, tileFileEnding),
		os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("could not write tile %s: %w", name, err)
	}

	err = gob.NewEncoder(file).Encode(tile)
	if err != nil {
		return fmt.Errorf("could not encode tile %s: %w", name, err)
	}

	return nil
}
