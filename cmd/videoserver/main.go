package main

import (
	"fmt"
	"github.com/bearname/videohost/pkg/common/infrarstructure/mysql"
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

	var connector mysql.ConnectorImpl
	err := connector.Connect()
	if err != nil {
		fmt.Println("unable to connect to connector" + err.Error())
		return
	}
	defer connector.Close()

	handler := router.Router(&connector)
	if handler == nil {
		return
	}

	server.ExecuteServer("videoserver", port, handler)
}
