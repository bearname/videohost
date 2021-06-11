package controller

import (
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/controller"
	"github.com/bearname/videohost/internal/user/domain"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
	"net/http"
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
	token, err, code := c.authService.CreateUser(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), code)
		return
	}

	c.WriteJsonResponse(writer, token)
}

func (c *AuthController) Login(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	token, err, code := c.authService.Login(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), code)
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

	token, err, code := c.authService.RefreshToken(request)

	if err != nil {
		http.Error(writer, err.Error(), code)
		return
	}

	c.WriteJsonResponse(writer, token)
}

func (c *AuthController) CheckTokenHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*request).Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		log.Println("check token handler")

		header := request.Header.Get("Authorization")
		userDto, statusCode := c.authService.ValidateToken(header)
		if statusCode != http.StatusOK {
			http.Error(writer, "Unauthorized or user not exists", statusCode)
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
	userDto, statusCode := c.authService.ValidateToken(header)
	if statusCode != http.StatusOK {
		http.Error(writer, "Unauthorized or user not exists", statusCode)
		return
	}

	c.BaseController.WriteJsonResponse(writer, userDto)
}
