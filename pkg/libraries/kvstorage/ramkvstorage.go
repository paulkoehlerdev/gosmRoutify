package kvstorage

import "fmt"

type ramimpl[K SerializableData, V SerializableData] ramCollectionManager[K, V]

type ramCollectionManager[K SerializableData, V SerializableData] map[string]Collection[K, V]

func NewRam[K SerializableData, V SerializableData](_ string, _ *KVStorageOptions) (KVStorage[K, V], error) {
	return make(ramimpl[K, V]), nil
}

func (r ramimpl[K, V]) GetCollection(name string) (Collection[K, V], error) {
	if coll, ok := r[name]; !ok {
		return nil, fmt.Errorf("collection not found")
	} else {
		return coll, nil
	}
}

func (r ramimpl[K, V]) NewCollection(name string) (Collection[K, V], error) {
	if _, ok := r[name]; !ok {
		r[name] = make(ramCollection[K, V])
	}
	return r[name], nil
}

func (r ramimpl[K, V]) Close() error {
	return nil
}

type ramCollection[K SerializableData, V SerializableData] map[string][]byte

func (r ramCollection[K, V]) Get(key K) (*V, error) {
	k, err := serializeData(key)
	if err != nil {
		return nil, fmt.Errorf("error while serializing key: %s", err.Error())
	}

	if v, ok := r[string(k)]; ok {
		out := new(V)
		err = deserializeData(v, out)
		if err != nil {
			return nil, fmt.Errorf("error while deserializing value: %s", err.Error())
		}
		return out, nil
	} else {
		return nil, fmt.Errorf("key not found")
	}
}

func (r ramCollection[K, V]) Set(key K, value V) error {
	k, err := serializeData(key)
	if err != nil {
		return fmt.Errorf("error while serializing key: %s", err.Error())
	}

	v, err := serializeData(value)
	if err != nil {
		return fmt.Errorf("error while serializing value: %s", err.Error())
	}

	r[string(k)] = v
	return nil
}
