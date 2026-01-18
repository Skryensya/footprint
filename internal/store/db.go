package store

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Skryensya/footprint/internal/store/migrations"
)

var (
	db        *sql.DB
	once      sync.Once
	openError error
)

// Open opens the database and runs any pending migrations.
// Uses singleton pattern - subsequent calls return the same connection.
func Open(path string) (*sql.DB, error) {
	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite3", path)
		if err != nil {
			openError = fmt.Errorf("open database: %w", err)
			return
		}

		// Verify connection
		if err = db.Ping(); err != nil {
			openError = fmt.Errorf("ping database: %w", err)
			return
		}

		// Run migrations
		if err = migrations.Run(db); err != nil {
			openError = fmt.Errorf("run migrations: %w", err)
			return
		}
	})

	return db, openError
}

// OpenFresh opens a new database connection without singleton.
// Used for testing with in-memory databases.
func OpenFresh(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err = migrations.Run(db); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return db, nil
}

// ResetSingleton resets the singleton state. Only for testing.
func ResetSingleton() {
	once = sync.Once{}
	db = nil
	openError = nil
}
