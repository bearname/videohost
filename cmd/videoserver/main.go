package main

import (
	"fmt"
	"github.com/bearname/videohost/cmd/videoserver/config"
	"github.com/bearname/videohost/internal/common/infrarstructure/mysql"
	"github.com/bearname/videohost/internal/common/infrarstructure/server"
	"github.com/bearname/videohost/internal/videoserver/infrastructure/transport/router"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func main() {
	//logFile := "video.log"
	//log.SetFormatter(&log.JSONFormatter{})
	//file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	//if err == nil {
	//	log.SetOutput(file)
	//	defer file.Close()
	//}
	parseConfig, err := config.ParseConfig()
	if err != nil {
		log.WithError(err).Fatal("failed to parse parseConfig")
	}

	connector := mysql.ConnectorImpl{}

	err = connector.Connect(parseConfig.DbUser, parseConfig.DbPassword, parseConfig.DbAddress, parseConfig.DbName)
	if err != nil {
		fmt.Println("unable to connect to connector" + err.Error())
		return
	}
	defer connector.Close()
	connector.SetMaxOpenConns(10)
	connector.SetConnMaxIdleTime(100)

	handler := router.Router(&connector, parseConfig.MessageBrokerAddress, parseConfig.AuthServerAddress, parseConfig.RedisAddress, parseConfig.RedisPassword)
	if handler == nil {
		return
	}

	server.ExecuteServer("videoserver", parseConfig.Port, handler)
}
