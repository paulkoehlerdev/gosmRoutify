package graph

import "fmt"

type ObjectID uint32

type GraphID uint64

func NewGraphID(lID LevelID, tID TileID, objectID ObjectID) GraphID {
	if lID > (1<<2)-1 {
		panic(fmt.Sprintf("tileLevel %d is too big (max: %d)", lID, 1<<2))
	}

	if tID > (1<<32)-1 {
		panic(fmt.Sprintf("TileID %d is too big (max: %d)", tID, 1<<32))
	}

	if objectID > (1<<29)-1 {
		panic(fmt.Sprintf("ObjectID %d is too big (max: %d)", objectID, 1<<29))
	}

	return GraphID(lID)<<61 | GraphID(tID)<<29 | GraphID(objectID)
}

func (g GraphID) Level() LevelID {
	return LevelID((g >> 61) & 0x7)
}

func (g GraphID) TileID() TileID {
	return TileID((g >> 29) & 0x7FFFFFFF)
}

func (g GraphID) ObjectID() ObjectID {
	return ObjectID(g & 0x1FFFFFFF)
}
