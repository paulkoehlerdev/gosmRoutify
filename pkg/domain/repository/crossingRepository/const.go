package crossingRepository

const (
	dataModel = `
`

	selectCrossingsFromWayID = `
SELECT osm_id, lat, lon, tags, is_crossing FROM node 
	JOIN wayToNodeRelation AS relation ON node.osm_id = relation.node_id 
	WHERE relation.way_id = ? 
	ORDER BY relation.position ASC;
`
)
