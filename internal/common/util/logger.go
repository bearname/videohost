package util

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func SetupLogger(logFile string) error {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}
	log.Info("Started")
	return err
}
