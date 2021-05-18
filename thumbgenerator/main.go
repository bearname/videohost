package main

import (
	"database/sql"
	"github.com/bearname/videohost/thumbgenerator/worker"
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

	db, err := sql.Open("mysql", `root:123@tcp(127.0.0.1)/video`)
	if err != nil {
		log.Error(err)
	}
	defer db.Close()

	rand.Seed(time.Now().Unix())
	stopChan := make(chan struct{})

	killChan := getKillSignalChan()
	wg := worker.WorkerPool(stopChan, db)

	waitForKillSignal(killChan)
	stopChan <- struct{}{}
	wg.Wait()
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
