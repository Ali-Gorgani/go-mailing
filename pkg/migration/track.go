package migration

import (
	"database/sql"
	"go-mailing/internal/app/logging"
)

func isMigrationApplied(db *sql.DB, id string) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM migrations WHERE id = $1)`
	err := db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		log.WithError(err).Error("Could not check if migration is applied")
	}
	return exists
}

func markMigrationApplied(db *sql.DB, id string) error {
	_, err := db.Exec(`INSERT INTO migrations (id) VALUES ($1)`, id)
	if err != nil {
		return logging.LogAndReturnError("Could not mark migration as applied", err)
	}
	return nil
}

func unmarkMigrationApplied(db *sql.DB, id string) error {
	_, err := db.Exec(`DELETE FROM migrations WHERE id = $1`, id)
	if err != nil {
		return logging.LogAndReturnError("Could not unmark migration", err)
	}
	return nil
}
