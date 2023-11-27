package kvstorage

import (
	"fmt"
	"os"
)

type dataAccessLayer interface {
	OverwritePage(pageNum pageNumber, buf []byte) error
	WritePage(buf []byte) (pageNumber, error)
	ReadPage(num pageNumber) (*page, error)
	AllocateEmptyPage() *page
	GetFirstPageNumber() pageNumber
	GetPageSize() int
	Close() error
}

const (
	metaPageNumber  pageNumber = iota
	freeListNumber  pageNumber = iota
	firstPageNumber pageNumber = iota
)

type pageNumber int64

type page struct {
	num pageNumber
	buf []byte
}

type dal struct {
	file     *os.File
	pageSize int

	meta     meta
	freeList freeList
}

func newDataAccessLayer(path string, pageSize int) (dataAccessLayer, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("error while opening database file: %s", err.Error())
	}

	out := &dal{
		file:     file,
		pageSize: pageSize,
	}

	fstat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("error while getting file size: %s", err.Error())
	}

	if fstat.Size() == 0 {
		err = out.create()
		if err != nil {
			return nil, fmt.Errorf("error while creating database: %s", err.Error())
		}
	} else {
		err = out.init()
		if err != nil {
			return nil, fmt.Errorf("error while initializing database: %s", err.Error())
		}
	}

	return out, nil
}

func (d *dal) OverwritePage(pageNum pageNumber, buf []byte) error {
	if pageNum < firstPageNumber {
		return fmt.Errorf("error while overwriting page: page number %d is too small", pageNum)
	}
	return d.writePage(pageNum, buf)
}

func (d *dal) WritePage(buf []byte) (pageNumber, error) {
	pageNum := d.freeList.Allocate()

	err := d.writePage(pageNum, buf)
	if err != nil {
		d.freeList.Release(pageNum)
		return 0, fmt.Errorf("error while writing page: %s", err.Error())
	}

	return pageNum, nil
}

func (d *dal) ReadPage(pageNum pageNumber) (*page, error) {
	if pageNum < firstPageNumber {
		return nil, fmt.Errorf("error while reading page: page number %d is too small", pageNum)
	}
	return d.readPage(pageNum)
}

func (d *dal) GetFirstPageNumber() pageNumber {
	return firstPageNumber
}

func (d *dal) GetPageSize() int {
	return d.pageSize
}

func (d *dal) AllocateEmptyPage() *page {
	return d.allocateEmptyPage()
}

func (d *dal) Close() error {
	if d.file == nil {
		return nil
	}

	err := d.writeMeta()
	if err != nil {
		return fmt.Errorf("error while writing meta page: %s", err.Error())
	}

	err = d.writeFreeList()
	if err != nil {
		return fmt.Errorf("error while writing free list page: %s", err.Error())
	}

	err = d.file.Close()
	if err != nil {
		return fmt.Errorf("error while closing database file: %s", err.Error())
	}

	d.file = nil
	return nil
}

func (d *dal) create() error {
	err := d.createMeta()
	if err != nil {
		return fmt.Errorf("error while creating meta page: %s", err.Error())
	}

	err = d.createFreeList()
	if err != nil {
		return fmt.Errorf("error while creating free list page: %s", err.Error())
	}

	d.freeList.CloseUntil(firstPageNumber)

	return nil
}

func (d *dal) createFreeList() error {
	page := d.allocateEmptyPage()

	fl := newFreeList()
	err := fl.Serialize(page.buf)
	if err != nil {
		return fmt.Errorf("error while serializing free list page: %s", err.Error())
	}

	err = d.writePage(freeListNumber, page.buf)
	if err != nil {
		return fmt.Errorf("error while writing free list page: %s", err.Error())
	}

	d.freeList = fl
	return nil
}

func (d *dal) writeFreeList() error {
	page := d.allocateEmptyPage()

	err := d.freeList.Serialize(page.buf)
	if err != nil {
		return fmt.Errorf("error while serializing free list page: %s", err.Error())
	}

	err = d.writePage(freeListNumber, page.buf)
	if err != nil {
		return fmt.Errorf("error while writing free list page: %s", err.Error())
	}

	return nil
}

func (d *dal) createMeta() error {
	page := d.allocateEmptyPage()

	m := newMeta()
	err := m.Serialize(page.buf)
	if err != nil {
		return fmt.Errorf("error while serializing meta page: %s", err.Error())
	}

	err = d.writePage(metaPageNumber, page.buf)
	if err != nil {
		return fmt.Errorf("error while writing meta page: %s", err.Error())
	}

	d.meta = m
	return nil
}

func (d *dal) writeMeta() error {
	page := d.allocateEmptyPage()

	err := d.meta.Serialize(page.buf)
	if err != nil {
		return fmt.Errorf("error while serializing meta page: %s", err.Error())
	}

	err = d.writePage(metaPageNumber, page.buf)
	if err != nil {
		return fmt.Errorf("error while writing meta page: %s", err.Error())
	}

	return nil
}

func (d *dal) init() error {
	metaPage, err := d.readPage(metaPageNumber)
	if err != nil {
		return fmt.Errorf("error while reading meta page: %s", err.Error())
	}

	d.meta = newMeta()
	err = d.meta.Deserialize(metaPage.buf)
	if err != nil {
		return fmt.Errorf("error while deserializing meta page: %s", err.Error())
	}

	freeListPage, err := d.readPage(freeListNumber)
	if err != nil {
		return fmt.Errorf("error while reading free list page: %s", err.Error())
	}

	d.freeList = newFreeList()
	err = d.freeList.Deserialize(freeListPage.buf)
	if err != nil {
		return fmt.Errorf("error while deserializing free list page: %s", err.Error())
	}

	return nil
}

func (d *dal) writePage(num pageNumber, buf []byte) error {
	if d.freeList != nil {
		d.freeList.CloseUntil(num)
	}

	if len(buf) != d.pageSize {
		return fmt.Errorf("error while writing page: data has wrong length %d != %d", len(buf), d.pageSize)
	}

	offset := int64(num) * int64(d.pageSize)

	_, err := d.file.WriteAt(buf, offset)
	if err != nil {
		return fmt.Errorf("error while writing page: %s", err.Error())
	}
	return nil
}

func (d *dal) readPage(num pageNumber) (*page, error) {
	offset := int64(num) * int64(d.pageSize)

	p := d.allocateEmptyPage()

	n, err := d.file.ReadAt(p.buf, offset)
	if err != nil {
		return nil, fmt.Errorf("error while reading page: %s", err.Error())
	}

	if n != d.pageSize {
		return nil, fmt.Errorf("error while reading page: read wrong number of bytes %d != %d", n, d.pageSize)
	}

	return p, nil
}

func (d *dal) allocateEmptyPage() *page {
	return &page{
		buf: make([]byte, d.pageSize),
	}
}
