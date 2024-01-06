package kvstorage

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestBtreeLayerImpl_Set_Get(t *testing.T) {
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	var items []*item
	for i := 0; i < testAmountOfItems; i += 2 {
		items = append(items, &item{
			key:   []byte(fmt.Sprintf("Key %s", string(alphabet[i%len(alphabet)]))),
			value: []byte(fmt.Sprintf("Value %s", string(alphabet[len(alphabet)-1-i%len(alphabet)]))),
		})
	}

	err := os.Remove("test.db")
	if err != nil {
		return
	}

	dal, err := newDataAccessLayer("test.db", os.Getpagesize())
	if err != nil {
		t.Errorf("error while creating data access layer: %s", err.Error())
		return
	}
	defer dal.Close()

	btree, err := newBtreeLayer(dal, -1, 0.0125, 0.025)
	if err != nil {
		t.Errorf("error while creating btree layer: %s", err.Error())
		return
	}

	for _, item := range items {
		err := btree.Set(item.key, item.value)
		if err != nil {
			t.Errorf("error while setting item: %s", err.Error())
			return
		}
	}

	btreeRootPage := btree.GetRootPageNumber()

	err = dal.Close()
	if err != nil {
		t.Errorf("error while closing data access layer: %s", err.Error())
		return
	}

	dal, err = newDataAccessLayer("test.db", os.Getpagesize())
	if err != nil {
		t.Errorf("error while creating data access layer: %s", err.Error())
		return
	}
	defer dal.Close()

	btree, err = newBtreeLayer(dal, btreeRootPage, 0.0125, 0.025)

	for _, item := range items {
		value, err := btree.Get(item.key)
		if err != nil {
			t.Errorf("error while getting item '%s': %s", item.key, err.Error())
			return
		}

		if len(value) != len(item.value) {
			t.Errorf("error while getting item '%s': wrong value length %d != %d", item.key, len(value), len(item.value))
			return
		}

		if bytes.Compare(value, item.value) != 0 {
			t.Errorf("error while getting item '%s': wrong value %d != %d", item.key, value, item.value)
			return
		}
	}

	err = dal.Close()
	if err != nil {
		t.Errorf("error while closing data access layer: %s", err.Error())
		return
	}

	err = os.Remove("test.db")
	if err != nil {
		t.Errorf("error while removing test database: %s", err.Error())
		return
	}
}
