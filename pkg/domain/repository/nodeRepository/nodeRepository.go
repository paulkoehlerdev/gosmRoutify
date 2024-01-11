package nodeRepository

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/node"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/database"
	"sync"
)

type NodeRepository interface {
	Init() error
	InsertNode(node node.Node) error
	InsertNodes(nodes []node.Node) error

	SelectNodeFromID(id int64) (*node.Node, error)

	SelectNodeIDsFromWayID(wayID int64) ([]int64, error)
	SelectNodesFromWayID(wayID int64) ([]*node.Node, error)

	SelectNearNodesApprox(lat float64, lon float64, radius float64) ([]*node.Node, error)
}

type impl struct {
	db         database.Database
	buf        bytes.Buffer
	bufferLock sync.Mutex
	preparedStatements
}

type preparedStatements struct {
	insertNode *sql.Stmt

	selectNodeFromID *sql.Stmt

	selectNodesFromWayID   *sql.Stmt
	selectNodeIDsFromWayID *sql.Stmt

	selectNearNodes *sql.Stmt
}

func New(db database.Database) NodeRepository {
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
	insertNode, err := i.db.Prepare(insertNode)
	if err != nil {
		return fmt.Errorf("error while preparing insert node statement: %s", err.Error())
	}

	selectNodeIDsFromWayID, err := i.db.Prepare(selectNodeIDsFromWayID)
	if err != nil {
		return fmt.Errorf("error while preparing insert way to node relation statement: %s", err.Error())
	}

	selectNodesFromWayID, err := i.db.Prepare(selectNodesFromWayID)
	if err != nil {
		return fmt.Errorf("error while preparing select nodes from way statement: %s", err.Error())
	}

	selectNearNodes, err := i.db.Prepare(selectNearNodes)
	if err != nil {
		return fmt.Errorf("error while preparing select nodes from way statement: %s", err.Error())
	}

	selectNodeFromID, err := i.db.Prepare(selectNodeFromID)
	if err != nil {
		return fmt.Errorf("error while preparing select from nodeid statement: %s", err.Error())
	}

	i.preparedStatements.insertNode = insertNode

	i.preparedStatements.selectNodeFromID = selectNodeFromID

	i.preparedStatements.selectNodeIDsFromWayID = selectNodeIDsFromWayID
	i.preparedStatements.selectNodesFromWayID = selectNodesFromWayID

	i.preparedStatements.selectNearNodes = selectNearNodes

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

func decodeTags(buf *bytes.Buffer) (map[string]string, error) {
	var tags map[string]string
	err := json.NewDecoder(buf).Decode(&tags)
	if err != nil {
		return nil, fmt.Errorf("error while decoding tags: %s", err.Error())
	}
	return tags, nil
}

func (i *impl) InsertNode(node node.Node) error {
	if i.preparedStatements.insertNode == nil {
		return fmt.Errorf("statements not prepared: you need to call Init() before you can call InsertNode()")
	}

	tags, err := i.encodeTags(node.Tags)
	if err != nil {
		return fmt.Errorf("error while encoding tags: %s", err.Error())
	}

	_, err = i.preparedStatements.insertNode.Exec(node.OsmID, node.Lat, node.Lon, tags)
	if err != nil {
		return fmt.Errorf("error while inserting node: %s", err.Error())
	}

	return nil
}

func (i *impl) InsertNodes(nodes []node.Node) error {
	if i.preparedStatements.insertNode == nil {
		return fmt.Errorf("statements not prepared: you need to call Init() before you can call InsertNode()")
	}

	tx, err := i.db.Begin()
	if err != nil {
		return fmt.Errorf("error while starting transaction: %s", err.Error())
	}

	insertNode := tx.Stmt(i.preparedStatements.insertNode)

	for _, node := range nodes {
		tags, err := i.encodeTags(node.Tags)
		if err != nil {
			return fmt.Errorf("error while encoding tags: %s", err.Error())
		}

		_, err = insertNode.Exec(node.OsmID, node.Lat, node.Lon, tags)
		if err != nil {
			return fmt.Errorf("error while inserting node: %s", err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error while committing transaction: %s", err.Error())
	}

	return nil
}

func (i *impl) SelectNodeFromID(id int64) (*node.Node, error) {
	if i.preparedStatements.selectNodeFromID == nil {
		return nil, fmt.Errorf("statements not prepared: you need to call Init() before you can call SelectNodeFromID()")
	}

	rows, err := i.preparedStatements.selectNodeFromID.Query(id)
	if err != nil {
		return nil, fmt.Errorf("error while querying nodes from way: %s", err.Error())
	}

	nodes, err := decodeNodes(rows)
	if err != nil {
		return nil, fmt.Errorf("error while decoding nodes: %s", err.Error())
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no node found")
	}

	return nodes[0], nil
}

func (i *impl) SelectNodeIDsFromWayID(wayID int64) ([]int64, error) {
	rows, err := i.preparedStatements.selectNodeIDsFromWayID.Query(wayID)
	if err != nil {
		return nil, fmt.Errorf("error while querying nodes from way: %s", err.Error())
	}

	var nodes []int64
	for rows.Next() {
		var nodeID int64
		err = rows.Scan(&nodeID)
		if err != nil {
			return nil, fmt.Errorf("error while scanning node id: %s", err.Error())
		}

		nodes = append(nodes, nodeID)
	}

	return nodes, nil
}

func decodeNodes(rows *sql.Rows) ([]*node.Node, error) {
	var nodes []*node.Node
	for rows.Next() {
		var node node.Node
		var buf []byte
		err := rows.Scan(&node.OsmID, &node.Lat, &node.Lon, &buf)
		if err != nil {
			return nil, fmt.Errorf("error while scanning node id: %s", err.Error())
		}

		buffer := bytes.NewBuffer(buf)
		node.Tags, err = decodeTags(buffer)
		if err != nil {
			return nil, fmt.Errorf("error while decoding tags: %s", err.Error())
		}

		nodes = append(nodes, &node)
	}

	return nodes, nil
}

func (i *impl) SelectNodesFromWayID(wayID int64) ([]*node.Node, error) {
	if i.preparedStatements.selectNodesFromWayID == nil {
		return nil, fmt.Errorf("statements not prepared: you need to call Init() before you can call SelectNodeIDsFromWayID()")
	}

	rows, err := i.preparedStatements.selectNodesFromWayID.Query(wayID)
	if err != nil {
		return nil, fmt.Errorf("error while selecting nodes from way: %s", err.Error())
	}
	defer rows.Close()

	nodes, err := decodeNodes(rows)
	if err != nil {
		return nil, fmt.Errorf("error while decoding nodes: %s", err.Error())
	}

	return nodes, nil
}

func (i *impl) SelectNearNodesApprox(lat float64, lon float64, radius float64) ([]*node.Node, error) {
	if i.preparedStatements.selectNearNodes == nil {
		return nil, fmt.Errorf("statements not prepared: you need to call Init() before you can call SelectNearNodesApprox()")
	}

	latMin := lat - radius
	latMax := lat + radius
	lonMin := lon - radius
	lonMax := lon + radius

	rows, err := i.preparedStatements.selectNearNodes.Query(
		latMin,
		latMax,
		lonMin,
		lonMax,
	)

	if err != nil {
		return nil, fmt.Errorf("error while querying nodes from way: %s", err.Error())
	}
	defer rows.Close()

	nodes, err := decodeNodes(rows)
	if err != nil {
		return nil, fmt.Errorf("error while decoding nodes: %s", err.Error())
	}

	return nodes, nil
}
