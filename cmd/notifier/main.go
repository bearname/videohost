package main

import (
	"github.com/bearname/videohost/pkg/common/amqp"
	"github.com/bearname/videohost/pkg/notifier/app/service"
	"github.com/bearname/videohost/pkg/notifier/intrastructure"
)

func main() {
	rabbitMqService := amqp.NewRabbitMqService("guest", "guest", "localhost", 5672)

	sender := service.NewSendInBlueMailSender("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")
	rabbitMqService.Consume("events_topic", "events.video-scaled", intrastructure.NewEmailSendConsumer(sender))
}
