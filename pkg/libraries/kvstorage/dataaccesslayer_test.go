package kvstorage

import (
	"bytes"
	"os"
	"testing"
)

const testPageSize = 512

func TestDal_WritePage(t *testing.T) {
	_ = os.Remove("test.db")

	dal, err := newDataAccessLayer("test.db", testPageSize)
	if err != nil {
		t.Errorf("error while creating database: %s", err.Error())
	}

	buf := make([]byte, testPageSize)
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i)
	}

	pageNum, err := dal.WritePage(buf)
	if err != nil {
		t.Errorf("error while writing page: %s", err.Error())
		return
	}

	err = dal.Close()
	if err != nil {
		t.Errorf("error while closing database: %s", err.Error())
		return
	}

	dal, err = newDataAccessLayer("test.db", testPageSize)
	if err != nil {
		t.Errorf("error while creating database: %s", err.Error())
		return
	}

	page, err := dal.ReadPage(pageNum)
	if err != nil {
		t.Errorf("error while reading page: %s", err.Error())
		return
	}

	if len(page.buf) != len(buf) {
		t.Errorf("error while reading page: wrong length %d != %d", len(page.buf), len(buf))
		return
	}

	if bytes.Compare(page.buf, buf) != 0 {
		t.Error("error while reading page: wrong data")
		return
	}

	err = dal.Close()
	if err != nil {
		t.Errorf("error while closing database: %s", err.Error())
		return
	}

	_ = os.Remove("test.db")
}
