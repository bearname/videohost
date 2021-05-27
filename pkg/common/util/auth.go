package util

import (
	"bytes"
	"encoding/json"
	"github.com/bearname/videohost/pkg/common/dto"
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
	body, err := getRequest(authorization, authServerUrl+"/api/v1/auth/token/validate")
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

func getRequest(authorization string, url string) ([]byte, error) {
	client := &http.Client{}
	validateAccessTokenRequest, err := http.NewRequest("GET", url, nil)
	validateAccessTokenRequest.Header.Add("Authorization", authorization)
	response, err := client.Do(validateAccessTokenRequest)
	if err != nil || response.StatusCode == http.StatusUnauthorized {
		log.Error(err.Error())
		return nil, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

func GetAdminAccessToken(client *http.Client, authServerUrl string) (*Token, error) {
	bodyStr, err := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{"admin", "admin"})
	req, err := http.NewRequest("POST", authServerUrl+"/api/v1/auth/login", bytes.NewBuffer(bodyStr))
	resp, err := client.Do(req)

	if err != nil {
		log.Error("Failed get id of owner of the video ")
		return &Token{}, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	token, err := unmarshalToken(err, respBody)
	if err != nil {
		return &Token{}, err
	}
	_, ok := ValidateToken("Bearer "+token.AccessToken, authServerUrl)
	if !ok {
		return updateToken(respBody, err, token, authServerUrl)
	}

	return token, nil
}

func updateToken(respBody []byte, err error, token *Token, authServerUrl string) (*Token, error) {
	respBody, err = getRequest("Bearer "+token.RefreshToken, authServerUrl+"/api/v1/auth/token")
	if err != nil {
		return nil, err
	}
	return unmarshalToken(err, respBody)
}

func unmarshalToken(err error, respBody []byte) (*Token, error) {
	var token Token
	err = json.Unmarshal(respBody, &token)
	if err != nil {
		return &token, err
	}
	return &token, nil
}
