package main

import (
	"github.com/bearname/videohost/cmd/notifier/config"
	"github.com/bearname/videohost/internal/common/infrarstructure/amqp"
	"github.com/bearname/videohost/internal/notifier/app/service"
	"github.com/bearname/videohost/internal/notifier/intrastructure"
)

func main() {
	conf := config.ParseConfig()
	rabbitMqService := amqp.NewRabbitMqService(conf.MessageBrokerAddress)

	sender := service.NewSendInBlueMailSender(conf.SendInBlueApiKey, conf.SendBlueAddress)
	rabbitMqService.Consume("events_topic", "events.video-scaled", intrastructure.NewEmailSendConsumer(sender, conf.AuthServerAddress))
}
