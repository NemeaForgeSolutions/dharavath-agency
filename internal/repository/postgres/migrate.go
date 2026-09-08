package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
)

//go:embed schema.sql
var embeddedSchemaSQL string

// ExecuteDDL executes the embedded schema SQL (DDL + seed).
// All tables, indexes, and constraints use IF NOT EXISTS, and all seeds use
// ON CONFLICT (id) DO UPDATE, guaranteeing that any missing tables, indexes,
// or curated listings are always intact and up-to-date across all environments
// without overwriting user-created records.
func ExecuteDDL(ctx context.Context, db *sql.DB) error {
	log.Printf("Verifying database schema and curated property listings...")
	if _, err := db.ExecContext(ctx, embeddedSchemaSQL); err != nil {
		return err
	}
	log.Printf("Database schema and curated listings verified successfully.")
	return nil
}
