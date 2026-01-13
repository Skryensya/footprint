package telemetry

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

func Open(path string) (*sql.DB, error) {
	var err error

	once.Do(func() {
		db, err = sql.Open("sqlite3", path)
	})

	return db, err
}
