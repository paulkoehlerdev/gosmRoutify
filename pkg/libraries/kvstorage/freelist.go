package kvstorage

import (
	"encoding/binary"
	"fmt"
)

type freeList interface {
	serializable
	CloseUntil(number pageNumber)
	Allocate() pageNumber
	Release(number pageNumber)
}

type freeListImpl struct {
	minPage       pageNumber
	maxPage       pageNumber
	releasedPages []pageNumber
}

func newFreeList() freeList {
	return &freeListImpl{
		maxPage:       0,
		minPage:       0,
		releasedPages: nil,
	}
}

func (f *freeListImpl) CloseUntil(number pageNumber) {
	f.minPage = number
}

func (f *freeListImpl) Allocate() pageNumber {
	if f.maxPage < f.minPage {
		f.maxPage = f.minPage
	}

	if len(f.releasedPages) == 0 {
		f.maxPage++
		return f.maxPage
	}

	p := f.releasedPages[0]
	f.releasedPages = f.releasedPages[1:]
	return p
}

func (f *freeListImpl) Release(number pageNumber) {
	if number < f.minPage {
		return
	}

	f.releasedPages = append(f.releasedPages, number)
}

func (f *freeListImpl) Serialize(buf []byte) error {
	index := 0

	binary.BigEndian.PutUint64(buf[index:], uint64(f.minPage))
	index += 8

	binary.BigEndian.PutUint64(buf[index:], uint64(f.maxPage))
	index += 8

	binary.BigEndian.PutUint16(buf[index:], uint16(len(f.releasedPages)))
	index += 2

	for _, p := range f.releasedPages {
		binary.BigEndian.PutUint64(buf[index:], uint64(p))
		index += 8

		if len(buf) < index {
			return fmt.Errorf("buffer is too small for free list")
		}
	}

	return nil
}

func (f *freeListImpl) Deserialize(buf []byte) error {
	index := 0

	f.minPage = pageNumber(binary.BigEndian.Uint64(buf[index:]))
	index += 8

	f.maxPage = pageNumber(binary.BigEndian.Uint64(buf[index:]))
	index += 8

	releasedPagesLen := binary.BigEndian.Uint16(buf[index:])
	index += 2

	f.releasedPages = make([]pageNumber, releasedPagesLen)

	for i := 0; i < int(releasedPagesLen); i++ {
		f.releasedPages[i] = pageNumber(binary.BigEndian.Uint64(buf[index:]))
		index += 8

		if len(buf) < index {
			return fmt.Errorf("buffer is too small for free list")
		}
	}

	return nil
}
