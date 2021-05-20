package util

import (
	"encoding/json"
	"io"
	"net/http"
)

func ValidateToken(authorization string, authServerUrl string) (string, bool) {
	client := &http.Client{}
	validateAccessTokenRequest, err := http.NewRequest("GET", "/api/v1/auth/token/validate", nil)
	validateAccessTokenRequest.Header.Add("Authorization", authorization)
	response, err := client.Do(validateAccessTokenRequest)
	if err != nil {
		return "", false
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	s := struct {
		Username string `json:"username"`
		UserId   string `json:"user_id"`
		ok       bool
	}{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return "", false
	}
	return s.UserId, true
}
