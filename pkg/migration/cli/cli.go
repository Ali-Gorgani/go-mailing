package main

import (
	"flag"
	"fmt"
	"go-mailing/internal/app/database"
	"go-mailing/pkg/migration"
)

// Run executes the migration CLI
func main() {
	if err := RunMigrationCLI(); err != nil {
		fmt.Println(err)
	}
}

func RunMigrationCLI() error {
	var action string
	flag.StringVar(&action, "action", "migrate-up", "Action to perform: migrate-up, migrate-up-by-number, migrate-down, migrate-down-by-number, current-version")
	var number int
	flag.IntVar(&number, "number", 1, "Number of migrations to apply for migrate-up-by-number and migrate-down-by-number actions")
	var path string
	flag.StringVar(&path, "path", "internal/app/migrations", "Path to the migrations directory")
	flag.Parse()

	db, err := database.Open(database.DefultPostgresConfig())
	if err != nil {
		return fmt.Errorf("could not open database: %w", err)
	}

	err = migration.LoadMigrationsFromDir(flag.Lookup("path").Value.String())
	if err != nil {
		return fmt.Errorf("could not load migrations: %w", err)
	}

	// Ensure the migrations table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS migrations (id VARCHAR(255) PRIMARY KEY)`)
	if err != nil {
		return fmt.Errorf("could not create migrations table: %w", err)
	}

	switch action {
	case "migrate-up":
		err = migration.MigrateUp(db)
		if err != nil {
			return fmt.Errorf("could not apply migrations: %w", err)
		}
	case "migrate-up-by-number":
		err = migration.MigrateUpByNumber(db, number)
		if err != nil {
			return fmt.Errorf("could not apply migrations: %w", err)
		}
	case "migrate-down":
		err = migration.MigrateDown(db)
		if err != nil {
			return fmt.Errorf("could not rollback migrations: %w", err)
		}
	case "migrate-down-by-number":
		err = migration.MigrateDownByNumber(db, number)
		if err != nil {
			return fmt.Errorf("could not rollback migrations: %w", err)
		}
	case "current-version":
		version, err := migration.CurrentVersion(db)
		if err != nil {
			return fmt.Errorf("could not get current version: %w", err)
		}
		fmt.Printf("Current version: %d\n", version)
	default:
		return fmt.Errorf("invalid action: %s", action)
	}
	return nil
}
