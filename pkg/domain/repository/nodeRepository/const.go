package nodeRepository

const (
	dataModel = `
CREATE TABLE IF NOT EXISTS node (
    osm_id INTEGER PRIMARY KEY UNIQUE NOT NULL,
    lat REAL NOT NULL,
    lon REAL NOT NULL,
    tags BLOB -- JSON
) STRICT;

CREATE INDEX IF NOT EXISTS node_osm_id_idx ON node (osm_id);
`

	insertNode = `
INSERT INTO node (osm_id, lat, lon, tags) VALUES (?, ?, ?, ?)
	ON CONFLICT (osm_id) DO UPDATE SET lat = excluded.lat, lon = excluded.lon, tags = excluded.tags;
`

	selectNodeFromID = `
SELECT osm_id, lat, lon, tags FROM node WHERE osm_id = ?;
`

	selectNodeIDsFromWayID = `
SELECT node_id FROM wayToNodeRelation
	WHERE way_id = ? 
	ORDER BY position ASC;
`

	selectNodesFromWayID = `
SELECT osm_id, lat, lon, tags FROM node 
	JOIN wayToNodeRelation AS relation ON node.osm_id = relation.node_id 
	WHERE relation.way_id = ? 
	ORDER BY relation.position ASC;
`

	selectNearNodes = `
SELECT node.osm_id, node.lat, node.lon, node.tags FROM node
  	WHERE node.lat BETWEEN ? AND ? AND node.lon BETWEEN ? AND ?
`
)