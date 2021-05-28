package main

import (
	config2 "github.com/bearname/videohost/cmd/video-scaler/config"
	"github.com/bearname/videohost/pkg/common/amqp"
	"github.com/bearname/videohost/pkg/video-scaler/app/service"
	"github.com/bearname/videohost/pkg/video-scaler/domain"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//err := util.SetupLogger("video-scale.log")
	//if err != nil {
	//	fmt.Println("unable setup log file" + err.Error())
	//}
	config := config2.ParseConfig()

	messageBroker := amqp.NewRabbitMqService(config.MessageBrokerAddress)
	qualities := []domain.Quality{domain.Q1080p, domain.Q720p, domain.Q480p, domain.Q320p}
	scaleService := service.NewVideoScaleService(messageBroker, config.VideoServerAddress, config.AuthServerAddress)
	handler := service.NewVideoScaleHandler(scaleService, qualities, config.VideoServerAddress, config.AuthServerAddress)
	messageBroker.Consume("events_topic", "events.video-uploaded", handler)
}
