package migration

import (
	"database/sql"
	"go-mailing/internal/app/logging"
)

var log = logging.GetLogger()

func MigrateUp(db *sql.DB) error {
	// Track if any migration was applied
	migrationsApplied := false

	for _, m := range Migrations {
		if !isMigrationApplied(db, m.ID) {
			if _, err := db.Exec(m.Up); err != nil {
				return logging.LogAndReturnError("Failed to apply migration", err)
			}
			if err := markMigrationApplied(db, m.ID); err != nil {
				return logging.LogAndReturnError("Failed to mark migration as applied", err)
			}
			// Log once after successful migration
			log.Infof("Successfully applied migration: %s", m.ID)
			migrationsApplied = true
		}
	}

	if !migrationsApplied {
		// If no migrations were applied and no error occurred, it means we are already at the latest version
		log.Info("Database is already at the latest version.")
	}

	return nil
}

func MigrateUpByNumber(db *sql.DB, number int) error {
	var migrationsApplied bool

	for i := 0; i < number; i++ {
		m := Migrations[i]
		if !isMigrationApplied(db, m.ID) {
			if _, err := db.Exec(m.Up); err != nil {
				return logging.LogAndReturnError("Failed to apply migration", err)
			}
			if err := markMigrationApplied(db, m.ID); err != nil {
				return logging.LogAndReturnError("Failed to mark migration as applied", err)
			}
			// Log once after successful migration
			log.Infof("Successfully applied migration: %s", m.ID)
			migrationsApplied = true
		}
	}

	if !migrationsApplied {
		// If no migrations were applied, it means we are already at the latest version
		log.Info("Database is already at the latest version or all specified migrations are already applied.")
	}

	return nil
}

func MigrateDown(db *sql.DB) error {
	for i := len(Migrations) - 1; i >= 0; i-- {
		m := Migrations[i]
		if isMigrationApplied(db, m.ID) {
			if _, err := db.Exec(m.Down); err != nil {
				return logging.LogAndReturnError("Failed to rollback migration", err)
			}
			if err := unmarkMigrationApplied(db, m.ID); err != nil {
				return logging.LogAndReturnError("Failed to unmark migration", err)
			}
			// Log once after successful rollback
			log.Infof("Successfully rolled back migration: %s", m.ID)
		}
	}
	return nil
}

func MigrateDownByNumber(db *sql.DB, number int) error {
	for i := 0; i < number; i++ {
		m := Migrations[len(Migrations)-1-i]
		if isMigrationApplied(db, m.ID) {
			if _, err := db.Exec(m.Down); err != nil {
				return logging.LogAndReturnError("Failed to rollback migration", err)
			}
			if err := unmarkMigrationApplied(db, m.ID); err != nil {
				return logging.LogAndReturnError("Failed to unmark migration", err)
			}
			// Log once after successful rollback
			log.Infof("Successfully rolled back migration: %s", m.ID)
		}
	}
	return nil
}

func CurrentVersion(db *sql.DB) (int, error) {
	var version int
	query := `SELECT COUNT(*) FROM migrations`
	err := db.QueryRow(query).Scan(&version)
	if err != nil {
		return 0, logging.LogAndReturnError("Could not get current version", err)
	}
	return version, nil
}
