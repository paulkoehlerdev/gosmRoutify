package wayRepository

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/way"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/database"
	"sync"
)

type WayRepository interface {
	Init() error
	InsertWay(node way.Way) error

	InsertWays(ways []way.Way) error

	SelectWayIDsFromNode(nodeID int64) ([]int64, error)
}

type impl struct {
	db         database.Database
	buf        bytes.Buffer
	bufferLock sync.Mutex
	preparedStatements
}

type preparedStatements struct {
	insertWay               *sql.Stmt
	insertWayToNodeRelation *sql.Stmt

	selectWayIDsFromNodeID *sql.Stmt
	selectWaysFromNodeID   *sql.Stmt
}

func New(db database.Database) WayRepository {
	return &impl{
		db: db,
	}
}

func (i *impl) Init() error {
	_, err := i.db.Exec(dataModel)
	if err != nil {
		return fmt.Errorf("error while creating data model: %s", err.Error())
	}

	err = i.prepareStatements()
	if err != nil {
		return fmt.Errorf("error while preparing statements: %s", err.Error())
	}

	return nil
}

func (i *impl) prepareStatements() error {
	insertWay, err := i.db.Prepare(insertWay)
	if err != nil {
		return fmt.Errorf("error while preparing insert way statement: %s", err.Error())
	}

	insertWayToNodeRelation, err := i.db.Prepare(insertWayToNodeRelation)
	if err != nil {
		return fmt.Errorf("error while preparing insert way to node relation statement: %s", err.Error())
	}

	selectWayIDsFromNodeID, err := i.db.Prepare(selectWayIDsFromNodeID)
	if err != nil {
		return fmt.Errorf("error while preparing select wayids ids from node statement: %s", err.Error())
	}

	selectWaysFromNodeID, err := i.db.Prepare(selectWaysFromNodeID)
	if err != nil {
		return fmt.Errorf("error while preparing select way ids from node statement: %s", err.Error())
	}

	i.preparedStatements.insertWay = insertWay
	i.preparedStatements.insertWayToNodeRelation = insertWayToNodeRelation

	i.preparedStatements.selectWayIDsFromNodeID = selectWayIDsFromNodeID
	i.preparedStatements.selectWaysFromNodeID = selectWaysFromNodeID

	return nil
}

func (i *impl) encodeTags(tags map[string]string) ([]byte, error) {
	i.bufferLock.Lock()
	defer i.bufferLock.Unlock()

	i.buf.Reset()
	defer i.buf.Reset()

	err := json.NewEncoder(&i.buf).Encode(tags)
	if err != nil {
		return nil, fmt.Errorf("error while encoding tags: %s", err.Error())
	}

	return i.buf.Bytes(), nil
}

func decodeTags(buf bytes.Buffer) (map[string]string, error) {
	var tags map[string]string
	err := json.NewDecoder(&buf).Decode(&tags)
	if err != nil {
		return nil, fmt.Errorf("error while decoding tags: %s", err.Error())
	}
	return tags, nil
}

func (i *impl) InsertWay(w way.Way) error {
	return i.InsertWays([]way.Way{w})
}

func (i *impl) InsertWays(ways []way.Way) error {
	if i.preparedStatements.insertWay == nil {
		return fmt.Errorf("statements not prepared: you need to call Init() before you can call InsertNode()")
	}

	if i.preparedStatements.insertWayToNodeRelation == nil {
		return fmt.Errorf("statements not prepared: you need to call Init() before you can call InsertNode()")
	}

	tx, err := i.db.Begin()
	if err != nil {
		return fmt.Errorf("error while starting transaction: %s", err.Error())
	}

	insertWay := tx.Stmt(i.preparedStatements.insertWay)
	insertWayToNodeRelation := tx.Stmt(i.preparedStatements.insertWayToNodeRelation)

	for _, way := range ways {
		tags, err := i.encodeTags(way.Tags)
		if err != nil {
			return fmt.Errorf("error while encoding tags: %s", err.Error())
		}

		_, err = insertWay.Exec(way.OsmID, tags)
		if err != nil {
			return fmt.Errorf("error while inserting way: %s", err.Error())
		}

		for _, nodeId := range way.Nodes {
			_, err = insertWayToNodeRelation.Exec(nodeId, way.OsmID)
			if err != nil {
				return fmt.Errorf("error while inserting way to node relation: %s", err.Error())
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error while committing transaction: %s", err.Error())
	}

	return nil
}

func (i *impl) SelectWayIDsFromNode(nodeID int64) ([]int64, error) {
	rows, err := i.preparedStatements.selectWayIDsFromNodeID.Query(nodeID)
	if err != nil {
		return nil, fmt.Errorf("error while querying ways from node: %s", err.Error())
	}

	var ways []int64
	for rows.Next() {
		var wayID int64
		err = rows.Scan(&wayID)
		if err != nil {
			return nil, fmt.Errorf("error while scanning way id: %s", err.Error())
		}

		ways = append(ways, wayID)
	}

	return ways, nil
}

func (i *impl) SelectWaysFromNode(nodeID int64) ([]way.Way, error) {
	rows, err := i.preparedStatements.selectWaysFromNodeID.Query(nodeID)
	if err != nil {
		return nil, fmt.Errorf("error while querying ways from node: %s", err.Error())
	}
	defer rows.Close()

	var ways []way.Way
	for rows.Next() {
		var way way.Way
		var buf bytes.Buffer
		err := rows.Scan(&way.OsmID, &buf)
		if err != nil {
			return nil, fmt.Errorf("error while scanning way id: %s", err.Error())
		}

		way.Tags, err = decodeTags(buf)
		if err != nil {
			return nil, fmt.Errorf("error while decoding tags: %s", err.Error())
		}

		ways = append(ways, way)
	}

	return ways, nil
}
