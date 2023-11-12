package graph

type LevelID uint8

const (
	Level0 LevelID = iota
	Level1 LevelID = iota
	Level2 LevelID = iota
)

func (t LevelID) ToTileWidth() float64 {
	switch t {
	case Level0:
		return 4
	case Level1:
		return 1
	case Level2:
		return 0.25
	default:
		return 0
	}
}

func LevelIDFromHighwayClass(highwayClass string) LevelID {
	switch highwayClass {
	case "motorway", "trunk", "primary", "motorway_junction", "motorway_link", "trunk_link", "primary_link":
		return Level0
	case "secondary", "tertiary", "secondary_link", "tertiary_link":
		return Level1
	default:
		return Level2
	}
}

func GetAllTileLevels() []LevelID {
	return []LevelID{Level0, Level1, Level2}
}
