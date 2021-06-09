package domain

type MailSender interface {
	Send(from User, to User, subject string, body string) error
}
