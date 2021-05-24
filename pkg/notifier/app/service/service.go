package service

import (
	"github.com/bearname/videohost/pkg/notifier/domain"
)

type MailSender interface {
	Send(from domain.User, to domain.User, subject string, body string) error
}
