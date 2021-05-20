package main

import (
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/thumbgenerator/app/worker"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Info("Started")
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("thumbgenerator.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	var connector database.Connector
	err = connector.Connect()
	if err != nil {
		panic("unable to connect to connector" + err.Error())
	}
	defer connector.Close()

	rand.Seed(time.Now().Unix())
	stopChan := make(chan struct{})

	killChan := getKillSignalChan()
	waitGroup := worker.WorkerPool(stopChan, connector.Database)

	waitForKillSignal(killChan)
	stopChan <- struct{}{}
	waitGroup.Wait()
	log.Info("Stopped")
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(killSignalChan chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM...")
	}
}
