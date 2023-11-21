package kvstorage

type KVStorage[K comparable, V any] interface {
	Get(key K) (V, error)
	Set(key K, value V) error
}
