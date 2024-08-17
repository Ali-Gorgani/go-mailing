package migration

import (
	"database/sql"
	"fmt"
)

func isMigrationApplied(db *sql.DB, id string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM migrations WHERE id = $1)`
	err := db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("is migration applied: %w", err)
	}
	return exists, nil
}

func markMigrationApplied(db *sql.DB, id string) error {
	_, err := db.Exec(`INSERT INTO migrations (id) VALUES ($1)`, id)
	if err != nil {
		return fmt.Errorf("mark migration applied: %w", err)
	}
	return nil
}

func unmarkMigrationApplied(db *sql.DB, id string) error {
	_, err := db.Exec(`DELETE FROM migrations WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("unmark migration applied: %w", err)
	}
	return nil
}
