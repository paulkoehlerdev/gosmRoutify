package kvstorage

import "fmt"

type RamKVStorage[K comparable, V any] map[K]V

func NewRamKVStorage[K comparable, V any](cap int) KVStorage[K, V] {
	return make(RamKVStorage[K, V], cap)
}

func (r RamKVStorage[K, V]) Get(key K) (V, error) {
	value, ok := r[key]
	if !ok {
		return *new(V), fmt.Errorf("error: key not found")
	}
	return value, nil
}

func (r RamKVStorage[K, V]) Set(key K, value V) error {
	r[key] = value
	return nil
}
