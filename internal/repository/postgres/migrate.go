package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
)

//go:embed schema.sql
var embeddedSchemaSQL string

// ExecuteDDL runs the full schema SQL (DDL + seed) only when the database
// tables have not been created yet, preventing seed data from re-inserting
// on every startup.
func ExecuteDDL(ctx context.Context, db *sql.DB) error {
	exists, err := tablesExist(ctx, db)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Database tables already exist, skipping schema migration.")
		return nil
	}

	log.Printf("Executing database schema migration and seed...")
	if _, err := db.ExecContext(ctx, embeddedSchemaSQL); err != nil {
		return err
	}
	log.Printf("Database schema and seed executed successfully.")
	return nil
}

// tablesExist returns true if the core application tables are already present.
func tablesExist(ctx context.Context, db *sql.DB) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name = 'users'
		)`,
	).Scan(&exists)
	return exists, err
}
