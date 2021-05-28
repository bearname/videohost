package main

import (
	config2 "github.com/bearname/videohost/cmd/notifier/config"
	"github.com/bearname/videohost/pkg/common/amqp"
	"github.com/bearname/videohost/pkg/notifier/app/service"
	"github.com/bearname/videohost/pkg/notifier/intrastructure"
)

func main() {
	config := config2.ParseConfig()
	rabbitMqService := amqp.NewRabbitMqService(config.MessageBrokerAddress)

	sender := service.NewSendInBlueMailSender(config.SendInBlueApiKey, config.SendBlueAddress)
	rabbitMqService.Consume("events_topic", "events.video-scaled", intrastructure.NewEmailSendConsumer(sender, config.AuthServerAddress))
}
