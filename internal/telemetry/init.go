package telemetry

import (
	"database/sql"
	_ "embed"
)

//go:embed schema.sql
var schema string

func Init(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}
