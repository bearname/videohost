package main

import (
	"fmt"
	"github.com/bearname/videohost/cmd/user/config"
	"github.com/bearname/videohost/internal/common/infrarstructure/mysql"
	"github.com/bearname/videohost/internal/common/infrarstructure/server"
	"github.com/bearname/videohost/internal/user/infrastructure/transport/router"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"runtime"
)

func main() {
	//logFile := "user.log"
	//
	//log.SetFormatter(&log.JSONFormatter{})
	//file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	//if err == nil {
	//	log.SetOutput(file)
	//	defer file.Close()
	//}
	runtime.GOMAXPROCS(runtime.NumCPU())
	conf, err := config.ParseConfig()
	if err != nil {
		log.WithError(err).Fatal("failed to parse conf")
	}

	connector := mysql.ConnectorImpl{}

	err = connector.Connect(conf.DbUser, conf.DbPassword, conf.DbAddress, conf.DbName)
	if err != nil {
		fmt.Println("unable to connect to connector" + err.Error())
		return
	}
	defer connector.Close()
	connector.SetMaxOpenConns(15)
	connector.SetConnMaxIdleTime(100)
	server.ExecuteServer("userserver", conf.Port, router.Router(&connector))
}
