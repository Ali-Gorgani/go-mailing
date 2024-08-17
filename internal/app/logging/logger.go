package logging

import (
	"go-mailing/configs"
	"os"

	"github.com/sirupsen/logrus"
)

func LoggingSetup(cfg configs.Config) *logrus.Logger {
	log := logrus.New()
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

	return log
}
