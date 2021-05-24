package intrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/bearname/videohost/pkg/notifier/app/service"
	"github.com/bearname/videohost/pkg/notifier/domain"
	"github.com/bearname/videohost/pkg/user/domain/model"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type EmailSendConsumer struct {
	mailSender        service.MailSender
	adminAccessToken  string
	adminRefreshToken string
}

func NewEmailSendConsumer(sender service.MailSender) *EmailSendConsumer {
	e := new(EmailSendConsumer)
	e.mailSender = sender
	e.adminAccessToken = ""
	e.adminAccessToken = ""
	return e
}

func (c *EmailSendConsumer) Handle(message string) error {
	log.Info(message)
	split := strings.Split(message, ",")
	if len(split) != 3 {
		return errors.New("Invalid message. format <videoId> <quality> <ownerId>")
	}
	videoId := split[0]
	quality := split[1]
	ownerId := split[2]
	client := &http.Client{}
	if c.adminAccessToken == "" {
		bodyStr, err := json.Marshal(struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{"admin", "admin"})
		req, err := http.NewRequest("POST", "http://localhost:8001/api/v1/auth/login", bytes.NewBuffer(bodyStr))
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Error("Failed get id of owner of the video ")
			return err
		}
		defer resp.Body.Close()

		v := struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		}{}
		respBody, err := io.ReadAll(resp.Body)
		err = json.Unmarshal(respBody, &v)
		if err != nil {
			return err
		}
		c.adminRefreshToken = v.RefreshToken
		c.adminAccessToken = v.AccessToken
	}
	req, err := http.NewRequest("GET", "http://localhost:8001"+"/api/v1/users/"+ownerId, nil)
	req.Header.Add("Authorization", "Bearer "+c.adminAccessToken)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error("Failed get id of owner of the video ")
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	var user model.User
	log.Info(respBody)
	err = json.Unmarshal(respBody, &user)
	if err != nil {
		return err
	}
	if !user.IsSubscribed {
		log.Info("User not subscribed")
		return nil
	}

	userEmail := user.Email
	body := "Video with id " + videoId + " available at " + quality + "p" + "by the following link http://localhost:8080/videos/" + videoId

	return c.mailSender.Send(domain.User{Name: "Sender Alex", Email: "senderalex@example.com"},
		domain.User{Name: strings.Split(userEmail, "@")[0], Email: userEmail},
		"video processing",
		body)
}
