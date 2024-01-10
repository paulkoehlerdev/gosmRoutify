
PRAGMA journal_mode = WAL;
PRAGMA busy_timeout = 5000;
PRAGMA foreign_keys = ON;

-- tags will be a json encoded map[string]string

CREATE TABLE IF NOT EXISTS node (
    osm_id INTEGER PRIMARY KEY,
    lat REAL,
    lon REAL,
    tags BLOB -- JSON
) STRICT;

CREATE TABLE IF NOT EXISTS wayToNodeRelation (
    node_id INTEGER FOREIGN KEY REFERENCES node(osm_id),
    way_id INTEGER FOREIGN KEY REFERENCES way(osm_id),
) STRICT;

CREATE TABLE IF NOT EXISTS way (
    osm_id INTEGER PRIMARY KEY,
    tags BLOB -- JSON
) STRICT;

CREATE TABLE IF NOT EXISTS relationToNodeRelation (
    node_id INTEGER FOREIGN KEY REFERENCES node(osm_id),
    relation_id INTEGER FOREIGN KEY REFERENCES relation(osm_id),
) STRICT;

CREATE TABLE IF NOT EXISTS relationToWayRelation (
    way_id INTEGER FOREIGN KEY REFERENCES way(osm_id),
    relation_id INTEGER FOREIGN KEY REFERENCES relation(osm_id),
) STRICT;


CREATE TABLE IF NOT EXISTS relationToRelationRelation (
    way_id INTEGER FOREIGN KEY REFERENCES way(osm_id),
    relation_id INTEGER FOREIGN KEY REFERENCES relation(osm_id),
) STRICT;

CREATE TABLE IF NOT EXISTS relation (
    osm_id INTEGER PRIMARY KEY,
    tags BLOB -- JSON
) STRICT;

CREATE VIRTUAL TABLE IF NOT EXISTS addresses USING FS5(
    lat REAL,
    lon REAL,
    address TEXT,
);