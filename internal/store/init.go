package store

import "database/sql"

// Init is deprecated. Migrations are now run automatically by Open.
// Kept for backward compatibility with existing code.
func Init(_ *sql.DB) error {
	return nil
}
