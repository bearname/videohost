package main

import (
	"fmt"
	"github.com/bearname/videohost/cmd/thumbgenerator/config"
	"github.com/bearname/videohost/pkg/common/infrarstructure/caching"
	"github.com/bearname/videohost/pkg/common/infrarstructure/mysql"
	"github.com/bearname/videohost/pkg/common/infrarstructure/server"
	"github.com/bearname/videohost/pkg/common/model"
	"github.com/bearname/videohost/pkg/thumbgenerator/app/worker"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"time"
)

func main() {
	logFile := "thumbgenerator.log"
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}
	parseConfig := config.ParseConfig()

	log.Info("Started")
	var connector mysql.ConnectorImpl
	err = connector.Connect(parseConfig.DbUser, parseConfig.DbPassword, parseConfig.DbAddress, parseConfig.DbName)
	if err != nil {
		fmt.Println("unable to connect to connector" + err.Error())
	}
	defer connector.Close()

	rand.Seed(time.Now().Unix())
	stopChan := make(chan struct{})

	killChan := server.GetKillSignalChan()
	cache := caching.NewRedisCache(model.NewDsn(parseConfig.RedisAddress, "", parseConfig.RedisPassword, ""))
	waitGroup := worker.PoolOfWorker(stopChan, &connector, cache)

	server.WaitForKillSignal(killChan)
	stopChan <- struct{}{}
	waitGroup.Wait()

	log.Info("Stopped")
}
