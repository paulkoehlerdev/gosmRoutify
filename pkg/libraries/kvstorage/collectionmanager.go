package kvstorage

import (
	"encoding/binary"
	"fmt"
)

type collectionManager[K SerializableData, V SerializableData] interface {
	GetCollection(name string) (Collection[K, V], error)
	NewCollection(name string) (Collection[K, V], error)
	Close() error
}

type collectionManagerMeta struct {
	rootPageNumber pageNumber
}

func (c *collectionManagerMeta) ToSerializable() serializable {
	return c
}

func (c *collectionManagerMeta) Serialize(buf []byte) error {
	index := 0

	binary.LittleEndian.PutUint64(buf[index:], uint64(c.rootPageNumber))
	index += 8

	return nil
}

func (c *collectionManagerMeta) Deserialize(buf []byte) error {
	index := 0

	c.rootPageNumber = pageNumber(binary.LittleEndian.Uint64(buf[index:]))
	index += 8

	return nil
}

type collectionManagerImpl[K SerializableData, V SerializableData] struct {
	btree          btreeLayer
	dal            dataAccessLayer
	meta           *collectionManagerMeta
	minFillPercent float64
	maxFillPercent float64
}

func newCollectionManager[K SerializableData, V SerializableData](dal dataAccessLayer, minFillPercent float64, maxFillPercent float64) (collectionManager[K, V], error) {
	c := collectionManagerImpl[K, V]{
		dal:            dal,
		minFillPercent: minFillPercent,
		maxFillPercent: maxFillPercent,
	}

	err := c.init()
	if err != nil {
		return nil, fmt.Errorf("error while initializing collection manager: %s", err.Error())
	}

	return &c, nil
}

func (c *collectionManagerImpl[K, V]) GetCollection(name string) (Collection[K, V], error) {
	collectionRootIndex, err := c.getCollectionRootIndex([]byte(name))
	if err != nil {
		return nil, fmt.Errorf("error while getting collection root index: %s", err.Error())
	}

	btree, err := newBtreeLayer(c.dal, collectionRootIndex, c.minFillPercent, c.maxFillPercent)
	if err != nil {
		return nil, fmt.Errorf("error while creating btree layer: %s", err.Error())
	}

	return newCollection[K, V](btree), nil
}

func (c *collectionManagerImpl[K, V]) getCollectionRootIndex(name []byte) (pageNumber, error) {
	rootIndexBytes, err := c.btree.Get(name)
	if err != nil {
		return 0, fmt.Errorf("error while getting collection root index: %s", err.Error())
	}

	if len(rootIndexBytes) != 8 {
		return 0, fmt.Errorf("invalid collection root index")
	}

	return pageNumber(binary.LittleEndian.Uint64(rootIndexBytes)), nil
}

func (c *collectionManagerImpl[K, V]) NewCollection(name string) (Collection[K, V], error) {
	newBtree, err := newBtreeLayer(c.dal, -1, c.minFillPercent, c.maxFillPercent)
	if err != nil {
		return nil, fmt.Errorf("error while creating btree layer: %s", err.Error())
	}

	newCollectionRootIndex := c.getCollectionRootIndexFromPage(newBtree.GetRootPageNumber())
	err = c.btree.Set([]byte(name), newCollectionRootIndex)
	if err != nil {
		return nil, fmt.Errorf("error while setting collection root index: %s", err.Error())
	}

	return newCollection[K, V](newBtree), nil
}

func (c *collectionManagerImpl[K, V]) getCollectionRootIndexFromPage(page pageNumber) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(page))
	return buf
}

func (c *collectionManagerImpl[K, V]) Close() error {
	err := c.writeMeta()
	if err != nil {
		return fmt.Errorf("error while writing meta page: %s", err.Error())
	}

	return nil
}

func (c *collectionManagerImpl[K, V]) init() error {
	page, err := c.dal.ReadPage(c.dal.GetFirstPageNumber())
	if err != nil {
		return c.create()
	}

	meta := &collectionManagerMeta{}
	err = meta.Deserialize(page.buf)
	if err != nil {
		return fmt.Errorf("error while deserializing collection manager meta: %s", err.Error())
	}
	c.meta = meta

	c.btree, err = newBtreeLayer(c.dal, c.meta.rootPageNumber, c.minFillPercent, c.maxFillPercent)
	if err != nil {
		return fmt.Errorf("error while creating btree layer: %s", err.Error())
	}

	return nil
}

func (c *collectionManagerImpl[K, V]) create() error {
	c.meta = &collectionManagerMeta{
		rootPageNumber: -1,
	}

	err := c.writeMeta()
	if err != nil {
		return fmt.Errorf("error while writing meta page: %s", err.Error())
	}

	c.btree, err = newBtreeLayer(c.dal, c.meta.rootPageNumber, c.minFillPercent, c.maxFillPercent)
	if err != nil {
		return fmt.Errorf("error while creating btree layer: %s", err.Error())
	}

	return nil
}

func (c *collectionManagerImpl[K, V]) writeMeta() error {
	page := c.dal.AllocateEmptyPage()

	if c.btree != nil {
		c.meta.rootPageNumber = c.btree.GetRootPageNumber()
	} else {
		c.meta.rootPageNumber = -1
	}

	err := c.meta.Serialize(page.buf)
	if err != nil {
		return fmt.Errorf("error while serializing collection manager meta: %s", err.Error())
	}

	err = c.dal.OverwritePage(c.dal.GetFirstPageNumber(), page.buf)
	if err != nil {
		return fmt.Errorf("error while overwriting collection manager meta: %s", err.Error())
	}

	return nil
}
