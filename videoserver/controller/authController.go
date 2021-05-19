package controller

import (
	"encoding/json"
	"github.com/bearname/videohost/videoserver/model"
	"github.com/bearname/videohost/videoserver/repository"
	"github.com/bearname/videohost/videoserver/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

type AuthController struct {
	BaseController
	userRepository *repository.UserRepository
}

func NewAuthController(userRepository *repository.UserRepository) *AuthController {
	v := new(AuthController)

	v.userRepository = userRepository
	return v
}

func (c *AuthController) GetTokenUserPassword(writer http.ResponseWriter, request *http.Request) {
	writer = *c.BaseController.AllowCorsRequest(&writer)
	var user model.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, "cannot decode username/password struct", http.StatusBadRequest)
		return
	}
	passwordHash, found := c.userRepository.Users[user.Username]
	if !found {
		http.Error(writer, "Cannot find the username", http.StatusNotFound)
	}
	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(user.Password))
	if err != nil {
		http.Error(writer, "Wrong password", http.StatusUnauthorized)
		return
	}
	token, err := util.CreateToken(user.Username)
	if err != nil {
		http.Error(writer, "Cannot create token", http.StatusInternalServerError)
		return
	}
	c.JsonResponse(writer, struct {
		Token string `json:"token"`
	}{token})
}
func (c *AuthController) CreateUser(writer http.ResponseWriter, request *http.Request) {
	writer = *c.BaseController.AllowCorsRequest(&writer)
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	var user model.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, "Cannot decode request", http.StatusBadRequest)
		return
	}
	if _, found := c.userRepository.Users[user.Username]; found {
		http.Error(writer, "User already exists", http.StatusBadRequest)
		return
	}
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	c.userRepository.Users[user.Username] = password
	token, err := util.CreateToken(user.Username)
	if err != nil {
		http.Error(writer, "Cannot create token", http.StatusInternalServerError)
		return
	}

	c.JsonResponse(writer, struct {
		Token string `json:"token"`
	}{token})
}

func (c *AuthController) CheckTokenHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		//writer = *c.BaseController.AllowCorsRequest(&writer)
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*request).Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		log.Println("check token handler")
		header := request.Header.Get("Authorization")
		bearerToken := strings.Split(header, " ")
		if len(bearerToken) != 2 {
			http.Error(writer, "Cannot read token", http.StatusBadRequest)
			return
		}
		if bearerToken[0] != "Bearer" {
			http.Error(writer, "Error in authorization token. it needs to be in form of 'Bearer <token>'", http.StatusBadRequest)
			return
		}

		token, ok := util.CheckToken(bearerToken[1])
		log.Println("bearerToken " + bearerToken[1])

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
			if _, ok := c.userRepository.Users[username]; !ok {
				http.Error(writer, "Unauthorized, user not exists", http.StatusUnauthorized)
				return
			}
			//Set the username in the request, so I will use it in check after!
			context.Set(request, "username", username)
		}

		log.Println("success")

		next(writer, request)
	}
}

func (c *AuthController) GetTokenByToken(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	//Here I already have the token checked... Just get the username from Request context
	username, ok := context.Get(request, "username").(string)
	if !ok {
		http.Error(writer, "Cannot check username", http.StatusInternalServerError)
		return
	}
	token, err := util.CreateToken(username)
	if err != nil {
		http.Error(writer, "Cannot create token", http.StatusInternalServerError)
		return
	}
	c.JsonResponse(writer, struct{ Token string }{token})
}
