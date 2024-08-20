package migration

import (
	"database/sql"
	"fmt"
)

func MigrateUp(db *sql.DB) error {
	// Track if any migration was applied
	migrationsApplied := false

	for _, m := range Migrations {
		condition, err := isMigrationApplied(db, m.ID)
		if err != nil {
			return fmt.Errorf("migrate up: %w", err)
		}
		if !condition {
			if _, err := db.Exec(m.Up); err != nil {
				return fmt.Errorf("migrate up: %w", err)
			}
			if err := markMigrationApplied(db, m.ID); err != nil {
				return fmt.Errorf("migrate up: %w", err)
			}
			// Log once after successful migration
			fmt.Printf("Successfully applied migration: %s\n", m.ID)
			migrationsApplied = true
		}
	}

	if !migrationsApplied {
		// If no migrations were applied and no error occurred, it means we are already at the latest version
		fmt.Println("Database is already at the latest version or all migrations are already applied.")
	}

	return nil
}

func MigrateUpByNumber(db *sql.DB, number int) error {
	var migrationsApplied bool

	for i := 0; i < number; i++ {
		m := Migrations[i]
		condition, err := isMigrationApplied(db, m.ID)
		if err != nil {
			return fmt.Errorf("migrate up by number: %w", err)
		}
		if !condition {
			if _, err := db.Exec(m.Up); err != nil {
				return fmt.Errorf("migrate up by number: %w", err)
			}
			if err := markMigrationApplied(db, m.ID); err != nil {
				return fmt.Errorf("migrate up by number: %w", err)
			}
			// Log once after successful migration
			fmt.Printf("Successfully applied migration: %s\n", m.ID)
			migrationsApplied = true
		}
	}

	if !migrationsApplied {
		// If no migrations were applied, it means we are already at the latest version
		fmt.Println("Database is already at the latest version or all migrations are already applied.")
	}

	return nil
}

func MigrateDown(db *sql.DB) error {
	for i := len(Migrations) - 1; i >= 0; i-- {
		m := Migrations[i]
		condition, err := isMigrationApplied(db, m.ID)
		if err != nil {
			return fmt.Errorf("migrate down: %w", err)
		}
		if condition {
			if _, err := db.Exec(m.Down); err != nil {
				return fmt.Errorf("migrate down: %w", err)
			}
			if err := unmarkMigrationApplied(db, m.ID); err != nil {
				return fmt.Errorf("migrate down: %w", err)
			}
			// Log once after successful rollback
			fmt.Printf("Successfully rolled back migration: %s\n", m.ID)
		}
	}
	return nil
}

func MigrateDownByNumber(db *sql.DB, number int) error {
	for i := 0; i < number; i++ {
		m := Migrations[len(Migrations)-1-i]
		condition, err := isMigrationApplied(db, m.ID)
		if err != nil {
			return fmt.Errorf("migrate down by number: %w", err)
		}
		if condition {
			if _, err := db.Exec(m.Down); err != nil {
				return fmt.Errorf("migrate down by number: %w", err)
			}
			if err := unmarkMigrationApplied(db, m.ID); err != nil {
				return fmt.Errorf("migrate down by number: %w", err)
			}
			// Log once after successful rollback
			fmt.Printf("Successfully rolled back migration: %s\n", m.ID)
		}
	}
	return nil
}

func CurrentVersion(db *sql.DB) (int, error) {
	var version int
	query := `SELECT COUNT(*) FROM migrations`
	err := db.QueryRow(query).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("current version: %w", err)
	}
	return version, nil
}
