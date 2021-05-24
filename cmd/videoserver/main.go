package main

import (
	"fmt"
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/common/infrarstructure/server"
	"github.com/bearname/videohost/pkg/common/util"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/transport/router"
	_ "github.com/go-sql-driver/mysql"
	"os"
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

	var connector database.Connector
	err := connector.Connect()
	if err != nil {
		fmt.Println("unable to connect to connector" + err.Error())
		return
	}
	defer connector.Close()

	server.ExecuteServer("videoserver", port, router.Router(connector))
}
