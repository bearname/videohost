package main

import (
	"fmt"
	"github.com/bearname/videohost/pkg/common/amqp"
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/common/util"
	"github.com/bearname/videohost/pkg/video-scaler/app/service"
	"github.com/bearname/videohost/pkg/video-scaler/domain"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/mysql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	err := util.SetupLogger("video-scale.log")
	if err != nil {
		fmt.Println("unable to connect to connector" + err.Error())
	}

	var connector database.Connector
	err = connector.Connect()
	if err != nil {
		log.Error("unable to connect to connector" + err.Error())
		var ok bool

		for i := 0; i < 4; i++ {
			time.Sleep(5 * time.Second)
			err = connector.Connect()
			if err == nil {
				ok = true
				break
			}
		}
		if !ok {
			log.Error("Failed connect to database")
			return
		}
	}

	defer connector.Close()
	videoRepo := mysql.NewMysqlVideoRepository(connector)
	messageBroker := amqp.NewRabbitMqService("guest", "guest", "localhost", 5672)
	qualities := []domain.Quality{domain.Q1440p, domain.Q1080p, domain.Q720p, domain.Q480p, domain.Q320p}
	scalerService := service.NewScalerService(messageBroker, videoRepo)
	handler := service.NewVideoScaleHandler(scalerService, videoRepo, qualities)
	messageBroker.Consume("events_topic", "events.video-uploaded", handler)
}
