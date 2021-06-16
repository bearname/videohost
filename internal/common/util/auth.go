package util

import (
	"bytes"
	"encoding/json"
	"github.com/bearname/videohost/internal/common/dto"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewToken(accessToken string, refreshToken string) *Token {
	t := new(Token)
	t.AccessToken = accessToken
	t.RefreshToken = refreshToken
	return t
}

func ValidateToken(authorization string, authServerUrl string) (dto.UserDto, bool) {
	if len(authorization) == 0 {
		return dto.UserDto{}, false
	}
	body, err := GetRequest(&http.Client{}, authServerUrl+"/api/v1/auth/token/validate", authorization)
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

func InitAccessToken(client *http.Client, authServerAddress string) (*Token, bool) {
	token, err := GetAdminAccessToken(client, authServerAddress)
	if err != nil {
		log.Error(err)
		return nil, false
	}
	return token, true
}

func GetAdminAccessToken(client *http.Client, authServerUrl string) (*Token, error) {
	bodyStr, err := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{"admin", "admin"})
	if err != nil {
		log.Error(err.Error())
		return &Token{}, err
	}
	req, err := http.NewRequest("POST", authServerUrl+"/api/v1/auth/login", bytes.NewBuffer(bodyStr))
	if err != nil {
		log.Error(err.Error())
		return &Token{}, err
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Error("failed get id of owner of the video ")
		return &Token{}, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed read response")
		return &Token{}, err
	}
	token, err := unmarshalToken(respBody)
	if err != nil {
		return &Token{}, err
	}
	_, ok := ValidateToken("Bearer "+token.AccessToken, authServerUrl)
	if !ok {
		return updateToken(token, authServerUrl)
	}

	return token, nil
}

func updateToken(token *Token, authServerUrl string) (*Token, error) {
	respBody, err := GetRequest(&http.Client{}, authServerUrl+"/api/v1/auth/token", "Bearer "+token.RefreshToken)
	if err != nil {
		return nil, err
	}
	return unmarshalToken(respBody)
}

func unmarshalToken(respBody []byte) (*Token, error) {
	var token Token
	err := json.Unmarshal(respBody, &token)
	if err != nil {
		return &token, err
	}
	return &token, nil
}
