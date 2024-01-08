package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// 7 minutes
const defaultDatabaseSettings = `
PRAGMA journal_mode = MEMORY;
PRAGMA busy_timeout = 5000;
PRAGMA foreign_keys = ON;

PRAGMA synchronous = OFF;
`

type Database interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Begin() (*sql.Tx, error)
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Close() error
}

type impl struct {
	db *sql.DB
}

func New(filename string) (Database, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, fmt.Errorf("error while opening database: %s", err.Error())
	}

	_, err = db.Exec(defaultDatabaseSettings)
	if err != nil {
		return nil, fmt.Errorf("error while setting database settings: %s", err.Error())
	}

	return &impl{
		db: db,
	}, nil
}

func (i *impl) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return i.db.Query(query, args...)
}

func (i *impl) Begin() (*sql.Tx, error) {
	return i.db.Begin()
}

func (i *impl) Prepare(query string) (*sql.Stmt, error) {
	return i.db.Prepare(query)
}

func (i *impl) Exec(query string, args ...interface{}) (sql.Result, error) {
	return i.db.Exec(query, args...)
}

func (i *impl) Close() error {
	return i.db.Close()
}
