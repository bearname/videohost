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
	config, err := config.ParseConfig()
	if err != nil {
		log.WithError(err).Fatal("failed to parse config")
	}

	connector := mysql.ConnectorImpl{}
	err = connector.Connect(config.DbUser, config.DbPassword, config.DbAddress, config.DbName)
	if err != nil {
		fmt.Println("unable to connect to connector" + err.Error())
		return
	}
	defer connector.Close()

	handler := router.Router(&connector, config.MessageBrokerAddress, config.AuthServerAddress, config.RedisAddress, config.RedisPassword)
	if handler == nil {
		return
	}

	server.ExecuteServer("videoserver", config.Port, handler)
}
