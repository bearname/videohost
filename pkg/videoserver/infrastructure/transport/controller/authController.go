package controller

import (
	"encoding/json"
	dto2 "github.com/bearname/videohost/pkg/videoserver/app/dto"
	"github.com/bearname/videohost/pkg/videoserver/app/service"
	model2 "github.com/bearname/videohost/pkg/videoserver/domain/model"
	repository2 "github.com/bearname/videohost/pkg/videoserver/domain/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type AuthController struct {
	BaseController
	userRepository repository2.UserRepository
}

func NewAuthController(userRepository repository2.UserRepository) *AuthController {
	v := new(AuthController)

	v.userRepository = userRepository
	return v
}

func (c *AuthController) CreateUser(writer http.ResponseWriter, request *http.Request) {
	writer = *c.BaseController.AllowCorsRequest(&writer)
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	var newUser dto2.UserDto
	err := json.NewDecoder(request.Body).Decode(&newUser)
	if err != nil {
		http.Error(writer, "Cannot decode request", http.StatusBadRequest)
		return
	}
	userFromDb, err := c.userRepository.FindByUserName(newUser.Username)
	if (err == nil && userFromDb.Username == newUser.Username) || (err != nil && err.Error() != "sql: no rows in result set") {
		http.Error(writer, "User already exists", http.StatusBadRequest)
		return
	}
	log.Println(userFromDb.Username == newUser.Username)

	userKey, err := uuid.NewUUID()
	if err != nil {
		http.Error(writer, "Generate id", http.StatusInternalServerError)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	accessToken, err := service.CreateToken(userKey.String(), newUser.Username, model2.General)
	if err != nil {
		http.Error(writer, "Cannot create accessToken", http.StatusInternalServerError)
		return
	}

	refreshToken, err := service.CreateTokenWithDuration(userKey.String(), newUser.Username, model2.General, time.Hour*24*365*10)
	if err != nil {
		http.Error(writer, "Cannot create refreshToken", http.StatusInternalServerError)
		return
	}

	err = c.userRepository.CreateUser(userKey.String(), newUser.Username, passwordHash, model2.General, accessToken, refreshToken)
	if err != nil {
		http.Error(writer, "User"+err.Error(), http.StatusInternalServerError)
		return
	}

	c.JsonResponse(writer, struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{accessToken,
		refreshToken})
}

func (c *AuthController) GetTokenUserPassword(writer http.ResponseWriter, request *http.Request) {
	writer = *c.BaseController.AllowCorsRequest(&writer)
	var userDto dto2.UserDto
	err := json.NewDecoder(request.Body).Decode(&userDto)
	if err != nil {
		http.Error(writer, "cannot decode username/password struct", http.StatusBadRequest)
		return
	}
	userFromDb, err := c.userRepository.FindByUserName(userDto.Username)
	if (err == nil && userFromDb.Username != userDto.Username) || err != nil {
		http.Error(writer, "User not exist", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword(userFromDb.Password, []byte(userDto.Password))
	if err != nil {
		http.Error(writer, "Wrong password", http.StatusUnauthorized)
		return
	}

	c.JsonResponse(writer, struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{userFromDb.AccessToken,
		userFromDb.RefreshToken})
}

func (c *AuthController) GetTokenByToken(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	username, ok := context.Get(request, "username").(string)
	if !ok {
		http.Error(writer, "Cannot check username", http.StatusInternalServerError)
		return
	}
	userKey, ok := context.Get(request, "userId").(string)
	if !ok {
		http.Error(writer, "Cannot check userId", http.StatusInternalServerError)
		return
	}

	userFromDb, err := c.userRepository.FindByUserName(username)
	if (err == nil && userFromDb.Username != username) || err != nil {
		http.Error(writer, "Unauthorized, user not exists", http.StatusUnauthorized)
		return
	}
	token, ok := context.Get(request, "token").(string)
	if !ok {
		http.Error(writer, "InternalServerError", http.StatusInternalServerError)
		log.Error("token to preset on context by token checker")
		return
	}

	if userFromDb.RefreshToken != token {
		http.Error(writer, "Invalid Refresh token", http.StatusBadRequest)
		return
	}

	accessToken, err := service.CreateToken(userKey, username, model2.General)
	if err != nil {
		http.Error(writer, "Cannot create accessToken", http.StatusInternalServerError)
		return
	}

	ok = c.userRepository.UpdateAccessToken(username, accessToken)
	if !ok {
		http.Error(writer, "Failed update accessToken", http.StatusInternalServerError)
		return
	}

	c.JsonResponse(writer, struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{accessToken,
		userFromDb.RefreshToken})
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
		tokenString, ok := c.parseToken(writer, header)
		if !ok {
			return
		}
		token, ok := service.CheckToken(tokenString)
		log.Println("bearerToken " + tokenString)

		if !ok {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			username, ok := claims["username"].(string)
			if !ok {
				http.Error(writer, "Unauthorized", http.StatusUnauthorized)
				return
			}
			userId, ok := claims["userId"].(string)
			if !ok {
				http.Error(writer, "Unauthorized", http.StatusUnauthorized)
				return
			}

			_, err := c.userRepository.FindByUserName(username)
			if err != nil {
				http.Error(writer, "Unauthorized, user not exists", http.StatusUnauthorized)
				return
			}

			context.Set(request, "username", username)
			context.Set(request, "userId", userId)
			context.Set(request, "token", tokenString)
		}

		log.Println("success")

		next(writer, request)
	}
}

func (c *AuthController) parseToken(writer http.ResponseWriter, header string) (string, bool) {
	bearerToken := strings.Split(header, " ")
	if len(bearerToken) != 2 {
		http.Error(writer, "Cannot read token", http.StatusBadRequest)
		return "", false
	}
	if bearerToken[0] != "Bearer" {
		http.Error(writer, "Error in authorization token. it needs to be in form of 'Bearer <token>'", http.StatusBadRequest)
		return "", false
	}

	return bearerToken[1], true
}
