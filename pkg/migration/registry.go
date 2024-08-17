package migration

import (
	"go-mailing/internal/app/logging"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Migration struct {
	ID   string
	Up   string
	Down string
}

var Migrations = []Migration{}

// LoadMigrationsFromDir reads all SQL files in a directory and loads them into the Migrations slice.
func LoadMigrationsFromDir(dirPath string) error {
	loadedMigrations := make(map[string]bool) // To track loaded migration IDs and avoid duplicates

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return logging.LogAndReturnError("Error while walking through the directory", err)
		}

		// Process only files with .sql extension
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".sql") {
			// Use the file name without extension as the ID
			id := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))

			// Skip if migration ID is already loaded
			if _, exists := loadedMigrations[id]; exists {
				log.Warnf("Skipping duplicate migration ID: %s", id)
				return nil
			}

			// Load migration from the file
			err := LoadMigrationFromFile(path, id)
			if err != nil {
				return logging.LogAndReturnError("Could not load migrations", err)
			}

			// Mark this migration ID as loaded
			loadedMigrations[id] = true
		}
		return nil
	})

	if err != nil {
		return logging.LogAndReturnError("Could not load migrations", err)
	}

	return nil
}

// LoadMigrationFromFile reads a migration file and adds it to the Migrations slice.
func LoadMigrationFromFile(filePath, id string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return logging.LogAndReturnError("Could not open the file", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return logging.LogAndReturnError("Could not read the file", err)
	}

	// Split the file content by "-- Migration Down"
	parts := strings.Split(string(content), "-- Migration Down")
	if len(parts) != 2 {
		return logging.LogAndReturnError("Invalid migration file format", nil)
	}

	upPart := strings.TrimSpace(strings.Replace(parts[0], "-- Migration Up", "", 1))
	downPart := strings.TrimSpace(parts[1])

	migration := Migration{
		ID:   id,
		Up:   upPart,
		Down: downPart,
	}

	Migrations = append(Migrations, migration)
	return nil
}
