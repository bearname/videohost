package main

import (
	"fmt"
	config2 "github.com/bearname/videohost/cmd/user/config"
	"github.com/bearname/videohost/pkg/common/infrarstructure/mysql"
	"github.com/bearname/videohost/pkg/common/infrarstructure/server"
	"github.com/bearname/videohost/pkg/user/infrastructure/transport/router"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := config2.ParseConfig()
	if err != nil {
		log.WithError(err).Fatal("failed to parse config")
	}

	connector := mysql.ConnectorImpl{}
	err = connector.Connect(config.DbUser, config.DbPassword, config.DbAddress, config.DbName)
	if err != nil {
		fmt.Println("unable to connect to connector" + err.Error())
		return
	}

	server.ExecuteServer("userserver", config.Port, router.Router(&connector))
}
