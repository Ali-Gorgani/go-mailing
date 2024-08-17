package main

import (
	"flag"
	"go-mailing/internal/app/database"
	"go-mailing/internal/app/logging"
	"go-mailing/pkg/migration"
)

// Run executes the migration CLI
func main() {
	log := logging.GetLogger()

	var action string
	flag.StringVar(&action, "action", "migrate-up", "Action to perform: migrate-up, migrate-up-by-number, migrate-down, migrate-down-by-number, current-version")
	var number int
	flag.IntVar(&number, "number", 1, "Number of migrations to apply for migrate-up-by-number and migrate-down-by-number actions")
	var path string
	flag.StringVar(&path, "path", "internal/app/migrations", "Path to the migrations directory")
	flag.Parse()

	db, err := database.Open(database.DefultPostgresConfig())
	if err != nil {
		log.WithError(err).Fatal("Could not open database")
		return
	}

	err = migration.LoadMigrationsFromDir(flag.Lookup("path").Value.String())
	if err != nil {
		log.WithError(err).Fatal("Could not load migrations")
		return
	}

	// Ensure the migrations table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS migrations (id VARCHAR(255) PRIMARY KEY)`)
	if err != nil {
		log.WithError(err).Fatal("Could not create migrations table")
		return
	}

	switch action {
	case "migrate-up":
		err = migration.MigrateUp(db)
		if err != nil {
			log.WithError(err).Fatal("Could not apply migrations")
			return
		}
	case "migrate-up-by-number":
		err = migration.MigrateUpByNumber(db, number)
		if err != nil {
			log.WithError(err).Fatal("Could not apply migrations")
			return
		}
	case "migrate-down":
		err = migration.MigrateDown(db)
		if err != nil {
			log.WithError(err).Fatal("Could not rollback migrations")
			return
		}
	case "migrate-down-by-number":
		err = migration.MigrateDownByNumber(db, number)
		if err != nil {
			log.WithError(err).Fatal("Could not rollback migration")
			return
		}
	case "current-version":
		version, err := migration.CurrentVersion(db)
		if err != nil {
			log.WithError(err).Fatal("Could not get current version")
			return
		}
		log.Infof("Current version: %d", version)
	default:
		log.WithError(err).Fatal("Invalid action")
		return
	}
}
