package main

import (
	"fmt"
	"github.com/bearname/videohost/cmd/videoserver/config"
	"github.com/bearname/videohost/pkg/common/infrarstructure/mysql"
	"github.com/bearname/videohost/pkg/common/infrarstructure/server"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/transport/router"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func main() {
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

	handler := router.Router(&connector, parseConfig.MessageBrokerAddress, parseConfig.AuthServerAddress, parseConfig.RedisAddress, parseConfig.RedisPassword)
	if handler == nil {
		return
	}

	server.ExecuteServer("videoserver", parseConfig.Port, handler)
}