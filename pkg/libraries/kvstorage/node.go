package kvstorage

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type item struct {
	key   []byte
	value []byte
}

const (
	pageNumSizeBytes    = 8
	nodeHeaderSizeBytes = 3
)

func (i *item) String() string {
	return fmt.Sprintf("(key: %s, value: %s)", string(i.key), string(i.value))
}

type node interface {
	serializable
	GetItem(index int) *item
	GetItemsAfter(index int) []*item
	GetItemsBefore(index int) []*item
	AddItemAt(index int, item *item)
	SetItems(items []*item)
	GetChildPageNumber(index int) pageNumber
	GetChildPageNumbersAfter(index int) []pageNumber
	GetChildPageNumbersBefore(index int) []pageNumber
	SetChildPageNumberAt(index int, child pageNumber)
	SetChildPageNumbers(children []pageNumber)
	GetSplitIndex(sizeThreshold float64) int
	FindIndex(key []byte) (index int, found bool)
	IsLeaf() bool
	GetSize() int
	GetPageNumber() pageNumber
	SetPageNumber(page pageNumber)
}

type nodeImpl struct {
	page     pageNumber
	items    []*item
	children []pageNumber
}

func newEmptyNode(page pageNumber) node {
	return &nodeImpl{
		page:     page,
		items:    nil,
		children: nil,
	}
}

func newNode(page pageNumber, items []*item, children []pageNumber) node {
	return &nodeImpl{
		page:     page,
		items:    items,
		children: children,
	}
}

func (n *nodeImpl) GetItem(index int) *item {
	if index < 0 || index >= len(n.items) {
		return nil
	}
	return n.items[index]
}

func (n *nodeImpl) GetItemsAfter(index int) []*item {
	if index < 0 || index >= len(n.items) {
		return nil
	}
	return n.items[index+1:]
}

func (n *nodeImpl) GetItemsBefore(index int) []*item {
	if index < 0 || index >= len(n.items) {
		return nil
	}
	return n.items[:index]
}

func (n *nodeImpl) SetItems(items []*item) {
	n.items = items
}

func (n *nodeImpl) AddItemAt(index int, item *item) {
	n.items = append(n.items, nil)
	copy(n.items[index+1:], n.items[index:])
	n.items[index] = item
}

func (n *nodeImpl) GetChildPageNumber(index int) pageNumber {
	if index < 0 || index >= len(n.children) {
		return -1
	}
	return n.children[index]
}

func (n *nodeImpl) GetChildPageNumbersAfter(index int) []pageNumber {
	if index < 0 || index >= len(n.children) {
		return nil
	}
	return n.children[index+1:]
}

func (n *nodeImpl) GetChildPageNumbersBefore(index int) []pageNumber {
	if index < 0 || index >= len(n.children) {
		return nil
	}
	return n.children[:index+1]
}

func (n *nodeImpl) SetChildPageNumberAt(index int, child pageNumber) {
	n.children = append(n.children, 0)
	copy(n.children[index+1:], n.children[index:])
	n.children[index] = child
}

func (n *nodeImpl) SetChildPageNumbers(children []pageNumber) {
	n.children = children
}

func (n *nodeImpl) FindIndex(key []byte) (index int, found bool) {
	left := 0
	right := len(n.items) - 1

	mid := 0
	for left <= right {
		mid = left + (right-left)/2
		switch bytes.Compare(n.items[mid].key, key) {
		case -1:
			left = mid + 1
		case 0:
			return mid, true
		case 1:
			right = mid - 1
		}
	}

	return right + 1, false
}

func (n *nodeImpl) IsLeaf() bool {
	return len(n.children) == 0
}

func (n *nodeImpl) getItemSize(i int) int {
	item := n.items[i]
	return pageNumSizeBytes + 2 + len(item.key) + 2 + len(item.value)
}

func (n *nodeImpl) GetSize() int {
	size := nodeHeaderSizeBytes
	for i := range n.items {
		size += n.getItemSize(i)
	}
	size += pageNumSizeBytes
	return size
}

func (n *nodeImpl) GetSplitIndex(sizeThreshold float64) int {
	index := 0
	size := nodeHeaderSizeBytes
	for i := range n.items {
		size += n.getItemSize(i)
		if float64(size) > sizeThreshold && index < len(n.items)-1 {
			return index + 1
		}
		index++
	}

	return -1
}

func (n *nodeImpl) GetPageNumber() pageNumber {
	return n.page
}

func (n *nodeImpl) SetPageNumber(page pageNumber) {
	n.page = page
}

func (n *nodeImpl) Serialize(buf []byte) error {
	leftIndex := 0
	rightIndex := len(buf)

	var isLeaf uint8 = 0
	if n.IsLeaf() {
		isLeaf = 1
	}

	binary.LittleEndian.PutUint16(buf[leftIndex:], uint16(len(n.items)))
	leftIndex += 2

	buf[leftIndex] = isLeaf
	leftIndex += 1

	if leftIndex != nodeHeaderSizeBytes {
		panic("wrong node header size. This is an implementation error")
	}

	for index, item := range n.items {
		if isLeaf != 1 {
			child := n.children[index]
			binary.LittleEndian.PutUint64(buf[leftIndex:], uint64(child))
			leftIndex += pageNumSizeBytes
		}

		keyLen := len(item.key)
		valueLen := len(item.value)

		rightIndex -= keyLen + valueLen + 4
		offset := rightIndex
		binary.LittleEndian.PutUint16(buf[leftIndex:], uint16(offset))
		leftIndex += 2

		binary.LittleEndian.PutUint16(buf[offset:], uint16(keyLen))
		offset += 2

		copy(buf[offset:], item.key)
		offset += keyLen

		binary.LittleEndian.PutUint16(buf[offset:], uint16(valueLen))
		offset += 2

		copy(buf[offset:], item.value)
		offset += valueLen

		if leftIndex >= rightIndex {
			return fmt.Errorf("buffer is too small for node")
		}
	}

	if isLeaf != 1 {
		child := n.children[len(n.children)-1]
		binary.LittleEndian.PutUint64(buf[leftIndex:], uint64(child))
		leftIndex += 8
	}

	if leftIndex >= rightIndex {
		return fmt.Errorf("buffer is too small for node")
	}
	return nil
}

func (n *nodeImpl) Deserialize(buf []byte) error {
	index := 0

	itemsLen := binary.LittleEndian.Uint16(buf[index:])
	index += 2

	isLeaf := buf[index]
	index += 1

	if index != nodeHeaderSizeBytes {
		panic("wrong node header size. This is an implementation error")
	}

	for i := uint16(0); i < itemsLen; i++ {
		if isLeaf != 1 {
			child := int64(binary.LittleEndian.Uint64(buf[index:]))
			index += pageNumSizeBytes

			n.children = append(n.children, pageNumber(child))
		}

		offset := binary.LittleEndian.Uint16(buf[index:])
		index += 2

		keyLen := binary.LittleEndian.Uint16(buf[offset:])
		offset += 2

		key := make([]byte, keyLen)
		copy(key, buf[offset:offset+keyLen])
		offset += keyLen

		valueLen := binary.LittleEndian.Uint16(buf[offset:])
		offset += 2

		value := make([]byte, valueLen)
		copy(value, buf[offset:offset+valueLen])
		offset += valueLen

		n.items = append(n.items, &item{
			key:   key,
			value: value,
		})
	}

	if isLeaf != 1 {
		child := binary.LittleEndian.Uint64(buf[index:])
		index += 8
		n.children = append(n.children, pageNumber(child))
	}

	return nil
}

func (n *nodeImpl) String() string {
	return fmt.Sprintf("(page: %d, items: %v, children: %v)", n.page, n.items, n.children)
}
