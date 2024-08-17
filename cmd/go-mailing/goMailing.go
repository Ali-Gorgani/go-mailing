package goMailing

import (
	"go-mailing/configs"
	"go-mailing/internal/app/database"
	"go-mailing/internal/app/logging"
	"go-mailing/internal/app/routes"
)

func StartServer() error {
	// Load the environment variables
	cfg, err := configs.LoadConfig("configs")
	if err != nil {
		return err
	}

	// Set up the logger
	log := logging.GetLogger()

	// Set up the database
	db, err := database.Open(database.DefultPostgresConfig())
	if err != nil {
		return logging.LogAndReturnError("Could not open database", err)
	}
	defer db.Close()

	// Set up the routes
	e := routes.NewRouter(db)

	// Start server
	log.Infof("Starting server on %s", cfg.Server.Address)
	if err := e.Start(cfg.Server.Address); err != nil {
		return logging.LogAndReturnError("Could not start server", err)
	}
	return nil
}
