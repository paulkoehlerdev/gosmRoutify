package tilerepository

import (
	"encoding/gob"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/graphtile"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/osmid"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/kdtree"
	"os"
)

type tileReader interface {
	Read(name string) (*graphtile.GraphTile, error)
}

type tileReaderImpl struct {
	tilePath string
}

func newTileReader(tilePath string) tileReader {
	gob.Register(kdtree.New[osmid.OsmID]())
	return &tileReaderImpl{
		tilePath: tilePath,
	}
}

func (g *tileReaderImpl) Read(name string) (*graphtile.GraphTile, error) {
	file, err := os.OpenFile(
		fmt.Sprintf("%s/%s%s", g.tilePath, name, tileFileEnding),
		os.O_RDONLY,
		0644,
	)
	if err != nil {
		return nil, fmt.Errorf("could not read tile %s: %w", name, err)
	}

	out := graphtile.GraphTile{}
	err = gob.NewDecoder(file).Decode(&out)
	if err != nil {
		return nil, fmt.Errorf("could not decode tile %s: %w", name, err)
	}

	return &out, nil
}
