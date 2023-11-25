package noderepository2

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/entity/nodetype"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/osmid"
	"io"
	"os"
)

const dataBlockLen = 1 << 6

type NodeData struct {
	Type nodetype.NodeType
}

type treeNode struct {
	Osmid       osmid.OsmID
	LeftOffset  int64
	RightOffset int64
	DataOffset  int64
	DataLength  int64
}

type NodeRepository interface {
	GetData(id osmid.OsmID) (*NodeData, error)
	SetData(id osmid.OsmID, data NodeData) error
}

type impl struct {
	database *os.File
	buffer   *bytes.Buffer
}

func New(databaseFileName string) (NodeRepository, error) {
	database, err := os.OpenFile(databaseFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("error while opening database file: %s", err.Error())
	}

	return &impl{
		database: database,
		buffer:   bytes.NewBuffer(make([]byte, 0, dataBlockLen)),
	}, nil
}

func (i *impl) GetData(id osmid.OsmID) (*NodeData, error) {
	root, err := i.getNode(0)
	if err != nil {
		return nil, fmt.Errorf("error while getting root: %s", err.Error())
	}

	node, err := i.findNode(root, id)
	if err != nil {
		return nil, fmt.Errorf("error while finding node: %s", err.Error())
	}

	data, err := i.getData(node)
	if err != nil {
		return nil, fmt.Errorf("error while getting data: %s", err.Error())
	}

	return data, nil
}

func (i *impl) SetData(id osmid.OsmID, data NodeData) error {
	root, err := i.getNode(0)
	if err != nil {
		_, err = i.writeNode(id, data)
		if err != nil {
			return fmt.Errorf("error while writing root node: %s", err.Error())
		}
		return nil
	}

	err, updated := i.insertNode(root, id, data)
	if err != nil {
		return fmt.Errorf("error while inserting node: %s", err.Error())
	}

	if !updated {
		return nil
	}

	err = i.updateNode(0, root)
	if err != nil {
		return fmt.Errorf("error while updating root: %s", err.Error())
	}

	return nil
}

func (i *impl) findNode(node *treeNode, id osmid.OsmID) (*treeNode, error) {
	if id == node.Osmid {
		return node, nil
	}

	if id < node.Osmid {
		if node.LeftOffset == -1 {
			return nil, fmt.Errorf("node not found")
		}

		node, err := i.getNode(node.LeftOffset)
		if err != nil {
			return nil, fmt.Errorf("error while getting node: %s", err.Error())
		}

		return i.findNode(node, id)
	}

	if node.RightOffset == -1 {
		return nil, fmt.Errorf("node not found")
	}

	node, err := i.getNode(node.RightOffset)
	if err != nil {
		return nil, fmt.Errorf("error while getting node: %s", err.Error())
	}

	return i.findNode(node, id)
}

func (i *impl) insertNode(node *treeNode, id osmid.OsmID, data NodeData) (err error, updated bool) {
	if id == node.Osmid {
		node.DataOffset, node.DataLength, err = i.insertData(data)
		if err != nil {
			return fmt.Errorf("error while inserting data: %s", err.Error()), false
		}
		return nil, true
	}

	if id < node.Osmid {
		if node.LeftOffset == -1 {
			offset, err := i.writeNode(id, data)
			if err != nil {
				return fmt.Errorf("error while writing node: %s", err.Error()), false
			}
			node.LeftOffset = offset
			return nil, true
		}

		childNode, err := i.getNode(node.LeftOffset)
		if err != nil {
			return fmt.Errorf("error while getting node: %s", err.Error()), false
		}

		err, updated := i.insertNode(childNode, id, data)
		if err != nil {
			return fmt.Errorf("error while inserting node: %s", err.Error()), false
		}

		if !updated {
			return nil, false
		}

		err = i.updateNode(node.LeftOffset, childNode)
		if err != nil {
			return fmt.Errorf("error while updating node: %s", err.Error()), false
		}

		return nil, false

	}

	if node.RightOffset == -1 {
		offset, err := i.writeNode(id, data)
		if err != nil {
			return fmt.Errorf("error while writing node: %s", err.Error()), false
		}
		node.RightOffset = offset
		return nil, true
	}

	childNode, err := i.getNode(node.RightOffset)
	if err != nil {
		return fmt.Errorf("error while getting node: %s", err.Error()), false
	}

	err, updated = i.insertNode(childNode, id, data)
	if err != nil {
		return fmt.Errorf("error while inserting node: %s", err.Error()), false
	}

	if !updated {
		return nil, false
	}

	err = i.updateNode(node.RightOffset, childNode)
	if err != nil {
		return fmt.Errorf("error while updating node: %s", err.Error()), false
	}

	return nil, false
}

func (i *impl) writeNode(id osmid.OsmID, data NodeData) (offset int64, err error) {
	offset, err = i.database.Seek(0, io.SeekEnd)
	if err != nil {
		return -1, fmt.Errorf("error while seeking file end: %s", err.Error())
	}

	root := treeNode{
		Osmid:       id,
		LeftOffset:  -1,
		RightOffset: -1,
		DataOffset:  -1,
		DataLength:  -1,
	}

	err = binary.Write(i.database, binary.LittleEndian, root)
	if err != nil {
		return -1, fmt.Errorf("error while writing tree root: %s", err.Error())
	}

	dataOffset, dataLength, err := i.insertData(data)
	if err != nil {
		return -1, fmt.Errorf("error while inserting data: %s", err.Error())
	}

	root.DataOffset = dataOffset
	root.DataLength = dataLength

	err = i.updateNode(offset, &root)
	if err != nil {
		return -1, fmt.Errorf("error while updating tree node: %s", err.Error())
	}

	return
}

func (i *impl) updateNode(offset int64, node *treeNode) error {
	_, err := i.database.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error while seeking to node at %d: %s", offset, err.Error())
	}

	err = binary.Write(i.database, binary.LittleEndian, node)
	if err != nil {
		return fmt.Errorf("error while writing node at %d: %s", offset, err.Error())
	}

	return nil
}

func (i *impl) getData(node *treeNode) (*NodeData, error) {
	_, err := i.database.Seek(node.DataOffset, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("error while seeking to data: %s", err.Error())
	}

	i.buffer.Reset()

	_, err = io.CopyN(i.buffer, i.database, node.DataLength)
	if err != nil {
		return nil, fmt.Errorf("error while reading data: %s", err.Error())
	}

	var data NodeData
	err = gob.NewDecoder(i.buffer).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("error while decoding data: %s", err.Error())
	}

	return &data, nil
}

func (i *impl) insertData(data NodeData) (offset int64, length int64, err error) {
	offset, err = i.database.Seek(0, io.SeekEnd)
	if err != nil {
		return -1, 0, fmt.Errorf("error while seeking to end of database: %s", err.Error())
	}

	i.buffer.Reset()
	err = gob.NewEncoder(i.buffer).Encode(data)
	if err != nil {
		return -1, 0, fmt.Errorf("error while encoding data: %s", err.Error())
	}

	if i.buffer.Len() > dataBlockLen {
		return -1, 0, fmt.Errorf("data too large")
	}

	_, err = i.database.Write(bytes.Repeat([]byte{0}, dataBlockLen))
	if err != nil {
		return -1, 0, fmt.Errorf("error while writing data block: %s", err.Error())
	}

	_, err = i.database.Seek(offset, io.SeekStart)
	if err != nil {
		return -1, 0, fmt.Errorf("error while returning to start of data block: %s", err.Error())
	}

	n, err := io.CopyN(i.database, i.buffer, int64(i.buffer.Len()))
	length = n

	return offset, length, nil
}

func (i *impl) getNode(offset int64) (*treeNode, error) {
	_, err := i.database.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("error while seeking to tree root at %d: %s", offset, err.Error())
	}

	var root treeNode
	err = binary.Read(i.database, binary.LittleEndian, &root)
	if err != nil {
		return nil, fmt.Errorf("error while reading tree node at %d: %s", offset, err.Error())
	}

	return &root, nil
}
