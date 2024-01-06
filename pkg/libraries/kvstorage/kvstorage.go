package kvstorage

import (
	"fmt"
	"os"
)

const (
	DefaultMinFillPercent = 0.5
	DefaultMaxFillPercent = 0.90
)

type KVStorageOptions struct {
	PageSize       int
	MinFillPercent float64
	MaxFillPercent float64
}

func DefaultKVStorageOptions() *KVStorageOptions {
	return &KVStorageOptions{
		PageSize:       os.Getpagesize(),
		MinFillPercent: DefaultMinFillPercent,
		MaxFillPercent: DefaultMaxFillPercent,
	}
}

type KVStorage[K SerializableData, V SerializableData] interface {
	GetCollection(name string) (Collection[K, V], error)
	NewCollection(name string) (Collection[K, V], error)
	Close() error
}

type impl[K SerializableData, V SerializableData] struct {
	collectionManager collectionManager[K, V]
	dal               dataAccessLayer
}

func New[K SerializableData, V SerializableData](path string, options *KVStorageOptions) (KVStorage[K, V], error) {
	dal, err := newDataAccessLayer(path, options.PageSize)
	if err != nil {
		return nil, fmt.Errorf("error while creating data access layer: %s", err.Error())
	}

	collectionManager, err := newCollectionManager[K, V](dal, options.MinFillPercent, options.MaxFillPercent)
	if err != nil {
		return nil, fmt.Errorf("error while creating collection manager: %s", err.Error())
	}

	return &impl[K, V]{
		collectionManager: collectionManager,
		dal:               dal,
	}, nil
}

func (i *impl[K, V]) GetCollection(name string) (Collection[K, V], error) {
	return i.collectionManager.GetCollection(name)
}

func (i *impl[K, V]) NewCollection(name string) (Collection[K, V], error) {
	return i.collectionManager.NewCollection(name)
}

func (i *impl[K, V]) Close() error {
	err := i.collectionManager.Close()
	if err != nil {
		return fmt.Errorf("error while closing collection manager: %s", err.Error())
	}

	err = i.dal.Close()
	if err != nil {
		return fmt.Errorf("error while closing data access layer: %s", err.Error())
	}
	return nil
}
