package main

import (
	"github.com/bearname/videohost/pkg/common/amqp"
	"github.com/bearname/videohost/pkg/notifier/app/service"
	"github.com/bearname/videohost/pkg/notifier/intrastructure"
	log "github.com/sirupsen/logrus"
	"net/smtp"
)

func main() {
	rabbitMqService := amqp.NewRabbitMqService("guest", "guest", "localhost", 5672)

	sender := service.NewSendInBlueMailSender("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")
	rabbitMqService.Consume("events_topic", "events.video-scaled", intrastructure.NewEmailSendConsumer(sender))
}

func sendEmail(from string, pass string, to string, body string) error {
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Video processing\n\n" +
		body

	err := smtp.SendMail("localhost:587",
		smtp.PlainAuth("", from, pass, "localhost"),
		from, []string{to}, []byte(msg))

	if err != nil {
		return err
	}

	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
