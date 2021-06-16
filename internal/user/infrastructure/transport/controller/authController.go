package controller

import (
	"encoding/json"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/controller"
	"github.com/bearname/videohost/internal/user/app/dto"
	"github.com/bearname/videohost/internal/user/domain"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"regexp"
	"strings"
)

type AuthController struct {
	controller.BaseController
	authService domain.AuthService
}

func NewAuthController(authService domain.AuthService) *AuthController {
	v := new(AuthController)
	v.authService = authService
	return v
}

func (c *AuthController) CreateUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	var newUserRequest dto.SignupUserDto
	err := json.NewDecoder(request.Body).Decode(&newUserRequest)
	if err != nil {
		log.Error(err.Error())
		http.Error(writer, "cannot signup request", http.StatusBadRequest)
		return
	}

	if !isEmailValid(newUserRequest.Email) {
		log.Error(err.Error())
		http.Error(writer, "user email not valid", http.StatusBadRequest)
		return
	}

	token, err := c.authService.CreateUser(newUserRequest)
	if err != nil {
		c.BaseController.WriteError(writer, err, TranslateError(err))
		return
	}

	c.WriteJsonResponse(writer, token)
}

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(e) {
		return false
	}
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}

func (c *AuthController) Login(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	var userDto dto.LoginUserDto
	err := json.NewDecoder(request.Body).Decode(&userDto)
	if err != nil {
		log.Error(err.Error())
		http.Error(writer, "cannot decode username/password struct", http.StatusBadRequest)
		return
	}

	token, err := c.authService.Login(userDto)
	if err != nil {
		c.BaseController.WriteError(writer, err, TranslateError(err))
		return
	}

	c.WriteJsonResponse(writer, token)
}

func (c *AuthController) RefreshToken(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	username, ok := context.Get(request, "username").(string)
	if !ok {
		context.Clear(request)
		http.Error(writer, "cannot check username", http.StatusBadRequest)
		return
	}
	userKey, ok := context.Get(request, "userId").(string)
	if !ok {
		context.Clear(request)
		http.Error(writer, "cannot check userId", http.StatusBadRequest)
		return
	}
	accessToken, ok := context.Get(request, "accessToken").(string)
	if !ok {
		context.Clear(request)
		http.Error(writer, "accessToken to preset on context by accessToken checker", http.StatusBadRequest)
		return
	}
	context.Clear(request)

	accessTokenResponse, err := c.authService.RefreshToken(dto.RefreshTokenDto{Username: username, UserId: userKey, Token: accessToken})

	if err != nil {
		c.BaseController.WriteError(writer, err, TranslateError(err))
		return
	}

	c.WriteJsonResponse(writer, accessTokenResponse)
}

func (c *AuthController) CheckTokenHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*request).Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		log.Println("check token handler")

		header := request.Header.Get("Authorization")
		userDto, err := c.authService.ValidateToken(header)
		if err != nil {
			c.BaseController.WriteError(writer, err, TranslateError(err))
			return
		}

		context.Set(request, "userId", userDto.UserId)
		context.Set(request, "username", userDto.Username)
		context.Set(request, "token", userDto.Token)

		log.Println("success")

		next(writer, request)
	}
}

func (c *AuthController) ValidateToken(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	log.Println("check token handler")

	header := request.Header.Get("Authorization")
	userDto, err := c.authService.ValidateToken(header)
	if err != nil {
		c.BaseController.WriteError(writer, err, TranslateError(err))
		return
	}

	c.BaseController.WriteJsonResponse(writer, userDto)
}
