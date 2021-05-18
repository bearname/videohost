package main

import (
	"context"
	"fmt"
	"github.com/bearname/videohost/videoserver/repository/mysql"
	"github.com/bearname/videohost/videoserver/router"
	"github.com/bearname/videohost/videoserver/util"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	port := 8000
	if len(os.Args) > 1 {
		toInt, ok := util.StrToInt(os.Args[1])
		if !ok {
			fmt.Println("Invalid port")
			return
		}
		port = toInt
	}

	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("videoserver.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	var connector mysql.Connector
	err = connector.Connect()
	if err != nil {
		panic("unable to connect to database" + err.Error())
	}
	defer connector.Close()

	killSignalChan := getKillSignalChan()

	serverUrl := ":" + strconv.Itoa(port)
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")
	srv := startServer(serverUrl, connector)

	waitForKillSignal(killSignalChan)
	srv.Shutdown(context.Background())
}

func startServer(serverUrl string, connector mysql.Connector) *http.Server {
	router := router.Router(connector)
	srv := &http.Server{Addr: serverUrl, Handler: router}
	go func() {
		log.Error(srv.ListenAndServe())
	}()

	return srv
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan <-chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM...")
	}
}
