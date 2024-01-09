package wayRepository

const (
	dataModel = `
CREATE TABLE IF NOT EXISTS wayToNodeRelation (
    node_id INTEGER NOT NULL,
    way_id INTEGER NOT NULL
) STRICT;


CREATE TABLE IF NOT EXISTS way (
    osm_id INTEGER PRIMARY KEY UNIQUE NOT NULL,
    tags BLOB -- JSON
) STRICT;

CREATE INDEX IF NOT EXISTS idx_way_highway_tag ON way((tags->>'$.highway'));

CREATE INDEX IF NOT EXISTS wayToNodeRelation_node_id_idx ON wayToNodeRelation (node_id);
CREATE INDEX IF NOT EXISTS wayToNodeRelation_way_id_idx ON wayToNodeRelation (way_id);

CREATE INDEX IF NOT EXISTS way_osm_id_idx ON way (osm_id);
`

	insertWay = `
INSERT INTO way (osm_id, tags) VALUES (?, ?)
	ON CONFLICT (osm_id) DO UPDATE SET tags = excluded.tags;
`

	insertWayToNodeRelation = `
INSERT INTO wayToNodeRelation (node_id, way_id) VALUES (?, ?);
`

	selectWayIDsFromNodeID = `
SELECT way_id FROM wayToNodeRelation WHERE node_id = ?;
`

	selectWaysFromNodeID = `
SELECT osm_id, tags FROM way JOIN wayToNodeRelation AS relation ON way.osm_id = relation.way_id WHERE relation.node_id = ?;
`

	selectWaysFromTwoNodeIDs = `
SELECT way.osm_id, way.tags
FROM way
 JOIN wayToNodeRelation AS rel1 ON rel1.way_id = way.osm_id
 JOIN wayToNodeRelation AS rel2 ON rel2.way_id = way.osm_id
WHERE rel1.node_id = ?
  AND rel2.node_id = ?;
`
)
