package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bearname/videohost/pkg/notifier/domain"
	"io/ioutil"
	"net/http"
)

type SendInBlueMailSender struct {
	apiKey string
}

func NewSendInBlueMailSender(apiKey string) *SendInBlueMailSender {
	s := new(SendInBlueMailSender)
	s.apiKey = apiKey
	return s
}

func (s *SendInBlueMailSender) Send(from domain.User, to domain.User, subject string, content string) error {
	message := domain.MailMessage{
		Sender:      from,
		To:          []domain.User{to},
		Subject:     subject,
		HtmlContent: content,
	}
	body, err := json.Marshal(message)
	fmt.Println(string(body))
	if err != nil {
		fmt.Println(err.Error())
	}
	request, err := http.NewRequest("POST", "https://api.sendinblue.com/v3/smtp/email", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err.Error())
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("api-key", s.apiKey)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	responseBody, err := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(responseBody))
	return err
}
