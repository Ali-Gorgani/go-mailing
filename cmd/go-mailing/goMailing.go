package goMailing

import (
	"go-mailing/configs"
	"go-mailing/internal/app/database"
	"go-mailing/internal/app/logging"
	"go-mailing/internal/app/routes"

	"github.com/sirupsen/logrus"
)

func StartServer() error {
	// Load the environment variables
	cfg, err := configs.LoadConfig("configs")
	if err != nil {
		logrus.Fatalf("Could not load the configuration: %v", err)
		return err
	}

	// Set up the logger
	log := logging.LoggingSetup(cfg)

	// Set up the database
	dbConfig := database.PostgresConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Database: cfg.DB.Database,
		SSLMode:  cfg.DB.SSLMode,
	}
	if dbConfig.Host == "" || dbConfig.Port == "" {
		dbConfig = database.DefultPostgresConfig()
	}
	db, err := database.Open(dbConfig)
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
		return err
	}
	defer db.Close()

	// Set up the routes
	e := routes.NewRouter(db, log, cfg)

	// Start server
	log.Infof("Starting server on %s", cfg.Server.Address)
	if err := e.Start(cfg.Server.Address); err != nil {
		log.Fatalf("Could not start server: %v", err)
		return err
	}
	return nil
}
