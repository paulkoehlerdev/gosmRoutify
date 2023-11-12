package graph

import "github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/osmpbfData"

type TileID uint32

func TileIDFromNode(node *osmpbfData.Node, Level LevelID) TileID {
	width := Level.ToTileWidth()
	id := uint32(node.Lat / width)
	id = id << 16
	id = id | uint32(node.Lon/width)
	return TileID(id)
}
