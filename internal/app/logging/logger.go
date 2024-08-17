package logging

import (
	"fmt"
	"go-mailing/configs"
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	cfg, err := configs.LoadConfig("configs")
	if err != nil {
		logrus.WithError(err).Fatal("Could not load configuration")
	}

	log = logrus.New() // Initialize the global log variable
	if cfg.Server.Environment == "production" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	// log.SetReportCaller(true)
}

// GetLogger returns the initialized logrus logger
func GetLogger() *logrus.Logger {
	return log
}

func LogAndReturnError(context string, err error) error {
	// Log the error with additional context
	logrus.WithFields(logrus.Fields{
		"context": context,
	}).Error(err)

	// Return a formatted error with the same context
	return fmt.Errorf("%s: %w", context, err)
}
