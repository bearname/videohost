package intrastructure

import (
	"errors"
	"github.com/bearname/videohost/pkg/notifier/app/service"
	"github.com/bearname/videohost/pkg/notifier/domain"
	log "github.com/sirupsen/logrus"
	"strings"
)

type EmailSendConsumer struct {
	mailSender service.MailSender
}

func NewEmailSendConsumer(sender service.MailSender) *EmailSendConsumer {
	e := new(EmailSendConsumer)
	e.mailSender = sender
	return e
}

func (c *EmailSendConsumer) Handle(message string) error {
	log.Info(message)
	split := strings.Split(message, ",")
	if len(split) != 3 {
		return errors.New("Invalid message. format <videoId> <quality> <user email>")
	}
	videoId := split[0]
	quality := split[1]
	to := split[2]
	body := "Video with id " + videoId + " available at " + quality + "p" + "by the following link http://localhost:8080/videos/" + videoId

	return c.mailSender.Send(domain.User{Name: "Sender Alex", Email: "senderalex@example.com"},
		domain.User{Name: strings.Split(to, "@")[0], Email: to},
		"video processing",
		body)
}
