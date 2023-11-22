package tilerepository

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/graphtile"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/coordinate"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"math"
	"os"
)

const tileFileEnding = ".tile"

type TileRepository interface {
	GetTile(coo coordinate.Coordinate) *graphtile.GraphTile
	SetTile(coo coordinate.Coordinate, tile *graphtile.GraphTile)
	Close()
}

type impl struct {
	logger     logging.Logger
	tileSize   float64
	tileFolder string

	tileCache tileCache
}

func New(logger logging.Logger, tileSize float64, tileFolder string, maxCacheSize int) TileRepository {
	err := os.MkdirAll(tileFolder, 0770)
	if err != nil {
		logger.Warn().Msgf("Cannot create graph directiory: %s", err.Error())
	}

	tileCache := newTileCache(
		maxCacheSize,
		newTileReader(tileFolder),
		newTileWriter(tileFolder),
		logger,
	)

	return &impl{
		logger:   logger,
		tileSize: tileSize,

		tileCache: tileCache,
	}
}

func (i *impl) GetTile(coo coordinate.Coordinate) *graphtile.GraphTile {
	tileName := i.TileNameFromCoordinate(coo)
	return i.getTile(tileName)
}

func (i *impl) getTile(tileName string) *graphtile.GraphTile {
	return i.tileCache.Get(tileName)
}

func (i *impl) SetTile(coo coordinate.Coordinate, tile *graphtile.GraphTile) {
	tileName := i.TileNameFromCoordinate(coo)
	i.setTile(tileName, tile)
}

func (i *impl) setTile(tileName string, tile *graphtile.GraphTile) {
	i.tileCache.Set(tileName, tile)
}

func (i *impl) TileNameFromCoordinate(coo coordinate.Coordinate) string {
	x := uint16(math.Floor(coo.Lat() / i.tileSize))
	y := uint16(math.Floor(coo.Lon() / i.tileSize))
	bytesX := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytesX, x)

	bytesY := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytesY, y)

	bytes := append(bytesX, bytesY...)
	return base64.StdEncoding.EncodeToString(bytes)
}

func (i *impl) Close() {
	i.logger.Info().Msg("closing tile repository")
	i.logger.Info().Msg("writing cache to disk")
	i.tileCache.WriteToDisk()
	i.logger.Info().Msg("tile repository closed")
}
