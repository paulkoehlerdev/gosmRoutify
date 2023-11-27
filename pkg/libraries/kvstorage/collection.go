package kvstorage

import (
	"fmt"
)

type Collection[K SerializableData, V SerializableData] interface {
	Get(key K) (*V, error)
	Set(key K, value V) error
}

type collectionImpl[K SerializableData, V SerializableData] struct {
	btreeLayer btreeLayer
}

func newCollection[K SerializableData, V SerializableData](btreeLayer btreeLayer) Collection[K, V] {
	return &collectionImpl[K, V]{
		btreeLayer: btreeLayer,
	}
}

func (c collectionImpl[K, V]) Get(key K) (*V, error) {
	buf, err := serializeData(key)
	if err != nil {
		return nil, fmt.Errorf("error while serializing key: %s", err.Error())
	}

	value, err := c.btreeLayer.Get(buf)
	if err != nil {
		return nil, fmt.Errorf("error while getting value: %s", err.Error())
	}

	v := new(V)
	err = deserializeData(value, v)
	if err != nil {
		return nil, fmt.Errorf("error while deserializing value: %s", err.Error())
	}

	return v, nil
}

func (c collectionImpl[K, V]) Set(key K, value V) error {
	keyBuf, err := serializeData(key)
	if err != nil {
		return fmt.Errorf("error while serializing key: %s", err.Error())
	}

	valBuf, err := serializeData(value)
	if err != nil {
		return fmt.Errorf("error while serializing value: %s", err.Error())
	}

	err = c.btreeLayer.Set(keyBuf, valBuf)
	if err != nil {
		return fmt.Errorf("error while setting value: %s", err.Error())
	}

	return nil
}
