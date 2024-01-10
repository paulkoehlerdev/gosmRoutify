package crossingRepository

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/crossing"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/database"
	"sync"
)

type CrossingRepository interface {
	Init() error
	SelectCrossingsFromWayID(wayID int64) ([]*crossing.Crossing, error)
}

type impl struct {
	db         database.Database
	buf        bytes.Buffer
	bufferLock sync.Mutex
	preparedStatements
}

type preparedStatements struct {
	selectCrossingsFromWayID *sql.Stmt
}

func New(db database.Database) CrossingRepository {
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
	selectCrossingsFromWayID, err := i.db.Prepare(selectCrossingsFromWayID)
	if err != nil {
		return fmt.Errorf("error while preparing statement: %s", err.Error())
	}
	i.preparedStatements.selectCrossingsFromWayID = selectCrossingsFromWayID

	return nil
}

func (i *impl) SelectCrossingsFromWayID(wayID int64) ([]*crossing.Crossing, error) {
	if i.preparedStatements.selectCrossingsFromWayID == nil {
		return nil, fmt.Errorf("statements not prepared: you need to call Init() before you can call SelectCrossingsFromWay()")
	}

	rows, err := i.preparedStatements.selectCrossingsFromWayID.Query(wayID)
	if err != nil {
		return nil, fmt.Errorf("error while selecting crossings from way: %s", err.Error())
	}
	defer rows.Close()

	crossings, err := decodeCrossings(rows)
	if err != nil {
		return nil, fmt.Errorf("error while decoding crossings: %s", err.Error())
	}

	return crossings, nil
}

func decodeCrossings(rows *sql.Rows) ([]*crossing.Crossing, error) {
	var crossings []*crossing.Crossing
	for rows.Next() {
		var crossing crossing.Crossing
		var buf []byte
		err := rows.Scan(&crossing.OsmID, &crossing.Lat, &crossing.Lon, &buf, &crossing.IsCrossing)
		if err != nil {
			return nil, fmt.Errorf("error while scanning crossing id: %s", err.Error())
		}

		buffer := bytes.NewBuffer(buf)
		crossing.Tags, err = decodeTags(buffer)
		if err != nil {
			return nil, fmt.Errorf("error while decoding tags: %s", err.Error())
		}

		crossings = append(crossings, &crossing)
	}

	return crossings, nil
}

func decodeTags(buffer *bytes.Buffer) (map[string]string, error) {
	var tags map[string]string
	err := json.NewDecoder(buffer).Decode(&tags)
	if err != nil {
		return nil, fmt.Errorf("error while decoding tags: %s, buffer: %s", err.Error(), buffer.String())
	}
	return tags, nil
}
