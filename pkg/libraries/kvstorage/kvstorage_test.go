package kvstorage_test

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/kvstorage"
	"os"
	"testing"
)

func TestImpl(t *testing.T) {
	_ = os.Remove("test.db")

	kv, err := kvstorage.New[string, int]("test.db", kvstorage.DefaultKVStorageOptions())
	if err != nil {
		t.Errorf("error while creating kvstorage: %s", err.Error())
		return
	}

	_, err = kv.NewCollection("test")
	if err != nil {
		t.Errorf("error while creating collection: %s", err.Error())
		return
	}

	collection, err := kv.GetCollection("test")
	if err != nil {
		t.Errorf("error while getting collection: %s", err.Error())
		return
	}

	err = collection.Set("key", 1)
	if err != nil {
		t.Errorf("error while setting value: %s", err.Error())
		return
	}

	value, err := collection.Get("key")
	if err != nil {
		t.Errorf("error while getting value: %s", err.Error())
		return
	}

	if value == nil {
		t.Errorf("error while getting value: value was nil: %v", value)
		return
	}

	if *value != 1 {
		t.Errorf("error while getting value: %d != %d", *value, 1)
		return
	}

	_, err = collection.Get("key2")
	if err == nil {
		t.Errorf("error while getting value: error was nil")
		return
	}

	err = collection.Set("key2", 1)
	if err != nil {
		t.Errorf("error while setting value: %s", err.Error())
		return
	}

	value, err = collection.Get("key2")
	if err != nil {
		t.Errorf("error while getting value: %s", err.Error())
		return
	}

	if value == nil {
		t.Errorf("error while getting value: value was nil: %v", value)
		return
	}

	if *value != 1 {
		t.Errorf("error while getting value: %d != %d", *value, 1)
		return
	}

	err = kv.Close()
	if err != nil {
		t.Errorf("error while closing kvstorage: %s", err.Error())
		return
	}

	err = os.Remove("test.db")
	if err != nil {
		t.Errorf("error while removing database: %s", err.Error())
		return
	}
}

func BenchmarkCollectionImpl_Set_Get(b *testing.B) {
	_ = os.Remove("test.db")

	kv, err := kvstorage.New[int, int]("test.db", kvstorage.DefaultKVStorageOptions())
	if err != nil {
		b.Errorf("error while creating kvstorage: %s", err.Error())
		return
	}

	collection, err := kv.NewCollection("test")
	if err != nil {
		b.Errorf("error while creating collection: %s", err.Error())
		return
	}

	fmt.Printf("Running with %d items\n", b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = collection.Set(i, b.N-i)
		if err != nil {
			b.Errorf("error while setting value for %d: %s", i, err.Error())
			return
		}
	}

	for i := 0; i < b.N; i++ {
		val, err := collection.Get(i)
		if err != nil {
			b.Errorf("error while getting value for %d/%d: %s", i, b.N, err.Error())
			continue
		}

		if val == nil {
			b.Errorf("error while getting value for %d: value was nil: %v", i, val)
			continue
		}

		if *val != b.N-i {
			b.Errorf("error while getting value for %d: %d != %d", i, val, b.N-i)
			continue
		}
	}

	b.StopTimer()

	err = kv.Close()
	if err != nil {
		b.Errorf("error while closing kvstorage: %s", err.Error())
		return
	}

	err = os.Remove("test.db")
	if err != nil {
		b.Errorf("error while removing database: %s", err.Error())
		return
	}
}
