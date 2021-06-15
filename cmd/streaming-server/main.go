package main

import (
	"github.com/bearname/videohost/internal/common/infrarstructure/server"
	"github.com/bearname/videohost/internal/common/util"
	"github.com/bearname/videohost/internal/stream-service/infrastructure/transport/router"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func main() {
	logFile := "stream.log"
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	portStr := util.ParseEnvString("PORT", "8020")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Info("Start with default port: 8020")
		log.WithError(err).Fatal("failed to parse parseConfig")
	}

	server.ExecuteServer("streaming-server", port, router.Router())
}
