package kvstorage

import (
	"bytes"
	"fmt"
	"testing"
)

const (
	testAmountOfItems = 20
)

func TestNodeImpl_Serialize_Deserialize(t *testing.T) {
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	var items []*item
	for i := 0; i < testAmountOfItems; i++ {
		items = append(items, &item{
			key:   []byte{alphabet[i%len(alphabet)]},
			value: []byte{alphabet[len(alphabet)-1-i%len(alphabet)]},
		})
	}

	node := nodeImpl{
		page:     0,
		items:    items,
		children: make([]pageNumber, testAmountOfItems+1),
	}

	buf := make([]byte, 512)
	err := node.Serialize(buf)
	if err != nil {
		t.Errorf("error while serializing node: %s", err.Error())
		return
	}

	deserializedNode := nodeImpl{}
	err = deserializedNode.Deserialize(buf)
	if err != nil {
		t.Errorf("error while deserializing node: %s", err.Error())
		return
	}

	fmt.Println(node.items)
	fmt.Println(deserializedNode.items)

	if len(deserializedNode.items) != len(node.items) {
		t.Errorf("error while deserializing node: wrong amount of items %d != %d", len(deserializedNode.items), len(node.items))
		return
	}

	for i := 0; i < len(deserializedNode.items); i++ {
		if len(deserializedNode.items[i].key) != len(node.items[i].key) {
			t.Errorf("error while deserializing node: wrong key length %d != %d", len(deserializedNode.items[i].key), len(node.items[i].key))
			return
		}

		if len(deserializedNode.items[i].value) != len(node.items[i].value) {
			t.Errorf("error while deserializing node: wrong value length %d != %d", len(deserializedNode.items[i].value), len(node.items[i].value))
			return
		}

		if bytes.Compare(deserializedNode.items[i].key, node.items[i].key) != 0 {
			t.Errorf("error while deserializing node: wrong key %s != %s", deserializedNode.items[i].key, node.items[i].key)
			return
		}

		if bytes.Compare(deserializedNode.items[i].value, node.items[i].value) != 0 {
			t.Errorf("error while deserializing node: wrong value %s != %s", deserializedNode.items[i].value, node.items[i].value)
			return
		}
	}
}

func TestNodeImpl_FindIndex(t *testing.T) {
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	var items []*item
	for i := 0; i < testAmountOfItems; i += 2 {
		items = append(items, &item{
			key:   []byte{alphabet[i%len(alphabet)]},
			value: []byte{alphabet[len(alphabet)-1-i%len(alphabet)]},
		})
	}

	node := nodeImpl{
		page:     0,
		items:    items,
		children: make([]pageNumber, testAmountOfItems+1),
	}

	for i := 0; i < testAmountOfItems; i++ {
		index, found := node.FindIndex([]byte{alphabet[i%len(alphabet)]})
		if i%2 == 0 && !found {
			t.Errorf("Expected to find key %s", string([]byte{alphabet[i%len(alphabet)]}))
		}

		if i%2 != 0 && found {
			t.Errorf("Expected not to find key %s", string([]byte{alphabet[i%len(alphabet)]}))
		}

		if (!found && index != i/2+1) || (found && index != i/2) {
			t.Errorf("Expected different index, got %d", index)
		}

		fmt.Printf("index: %d, found: %t\n", index, found)
	}

	index, found := node.FindIndex([]byte{byte(255)})
	if found {
		t.Errorf("Expected not to find key %s", string([]byte{byte(255)}))
	}

	if index != testAmountOfItems/2 {
		t.Errorf("Expected index %d, got %d", testAmountOfItems/2, index)
	}

	index, found = node.FindIndex([]byte{byte(0)})
	if found {
		t.Errorf("Expected not to find key %s", string([]byte{byte(0)}))
	}

	if index != 0 {
		t.Errorf("Expected index %d, got %d", 0, index)
	}
}
