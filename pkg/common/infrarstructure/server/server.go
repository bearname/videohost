package server

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func ExecuteServer(appName string, port int, router http.Handler) {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile(appName+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}
	log.Info("Started")

	killSignalChan := GetKillSignalChan()

	serverUrl := ":" + strconv.Itoa(port)
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")

	srv := StartServer(serverUrl, router)
	fmt.Println(serverUrl)
	WaitForKillSignal(killSignalChan)
	err = srv.Shutdown(context.Background())
	if err != nil {
		return
	}
}

func StartServer(serverUrl string, router http.Handler) *http.Server {
	srv := &http.Server{Addr: serverUrl, Handler: router}
	go func() {
		log.Error(srv.ListenAndServe())
	}()

	return srv
}

func GetKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func WaitForKillSignal(killSignalChan <-chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM...")
	}
}