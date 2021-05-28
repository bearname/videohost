package util

import (
	"bytes"
	"encoding/json"
	"github.com/bearname/videohost/pkg/common/dto"
	"github.com/bearname/videohost/pkg/video-scaler/domain"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func ValidateToken(authorization string, authServerUrl string) (dto.UserDto, bool) {
	body, err := GetRequest(authServerUrl+"/api/v1/auth/token/validate", authorization)
	if err != nil {
		return dto.UserDto{}, false
	}

	var userDto dto.UserDto
	err = json.Unmarshal(body, &userDto)
	if err != nil {
		log.Error(err.Error())
		return userDto, false
	}
	return userDto, true
}

func InitAccessToken(client *http.Client, authServerAddress string) (*domain.Token, bool) {
	token, err := GetAdminAccessToken(client, authServerAddress)
	if err != nil {
		log.Error(err)
		return nil, false
	}
	return token, true
}

func GetAdminAccessToken(client *http.Client, authServerUrl string) (*domain.Token, error) {
	bodyStr, err := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{"admin", "admin"})

	req, err := http.NewRequest("POST", authServerUrl+"/api/v1/auth/login", bytes.NewBuffer(bodyStr))
	resp, err := client.Do(req)

	if err != nil {
		log.Error("failed get id of owner of the video ")
		return &domain.Token{}, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	token, err := unmarshalToken(err, respBody)
	if err != nil {
		return &domain.Token{}, err
	}
	_, ok := ValidateToken("Bearer "+token.AccessToken, authServerUrl)
	if !ok {
		return updateToken(respBody, err, token, authServerUrl)
	}

	return token, nil
}

func updateToken(respBody []byte, err error, token *domain.Token, authServerUrl string) (*domain.Token, error) {
	respBody, err = GetRequest(authServerUrl+"/api/v1/auth/token", "Bearer "+token.RefreshToken)
	if err != nil {
		return nil, err
	}
	return unmarshalToken(err, respBody)
}

func unmarshalToken(err error, respBody []byte) (*domain.Token, error) {
	var token domain.Token
	err = json.Unmarshal(respBody, &token)
	if err != nil {
		return &token, err
	}
	return &token, nil
}
