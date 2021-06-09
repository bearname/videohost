package main

import (
	"github.com/bearname/videohost/cmd/video-scaler/config"
	"github.com/bearname/videohost/internal/common/infrarstructure/amqp"
	"github.com/bearname/videohost/internal/video-scaler/app/service"
	"github.com/bearname/videohost/internal/video-scaler/domain"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//err := util.SetupLogger("video-scale.log")
	//if err != nil {
	//	fmt.Println("unable setup log file" + err.Error())
	//}
	conf := config.ParseConfig()

	messageBroker := amqp.NewRabbitMqService(conf.MessageBrokerAddress)
	qualities := []domain.Quality{domain.Q1080p, domain.Q720p, domain.Q480p, domain.Q320p}
	scaleService := service.NewVideoScaleService(messageBroker, conf.VideoServerAddress, conf.AuthServerAddress)
	handler := service.NewVideoScaleHandler(scaleService, qualities, conf.VideoServerAddress, conf.AuthServerAddress)
	messageBroker.Consume("events_topic", "events.video-uploaded", handler)
}
