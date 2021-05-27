package controller

import (
	"encoding/json"
	commonDto "github.com/bearname/videohost/pkg/common/dto"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/controller"
	"github.com/bearname/videohost/pkg/user/app/dto"
	"github.com/bearname/videohost/pkg/user/app/service"
	"github.com/bearname/videohost/pkg/user/domain/model"
	"github.com/bearname/videohost/pkg/user/domain/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type AuthController struct {
	controller.BaseController
	userRepo repository.UserRepository
}

func NewAuthController(userRepository repository.UserRepository) *AuthController {
	v := new(AuthController)
	v.userRepo = userRepository
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
	var newUser dto.SignupUserDto
	err := json.NewDecoder(request.Body).Decode(&newUser)
	if err != nil {
		http.Error(writer, "Cannot decode request", http.StatusBadRequest)
		return
	}

	if !c.isEmailValid(newUser.Email) {
		c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "User email not valid")
		return
	}

	userFromDb, err := c.userRepo.FindByUserName(newUser.Username)
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
	role := model.General
	accessToken, err := service.CreateToken(userKey.String(), newUser.Username, role)
	if err != nil {
		http.Error(writer, "Cannot create accessToken", http.StatusInternalServerError)
		return
	}

	refreshToken, err := service.CreateTokenWithDuration(userKey.String(), newUser.Username, role, time.Hour*24*365*10)
	if err != nil {
		http.Error(writer, "Cannot create refreshToken", http.StatusInternalServerError)
		return
	}

	err = c.userRepo.CreateUser(userKey.String(), newUser.Username, passwordHash, newUser.Email, newUser.IsSubscribed, role, accessToken, refreshToken)
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
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	var userDto dto.LoginUserDto
	err := json.NewDecoder(request.Body).Decode(&userDto)
	if err != nil {
		http.Error(writer, "cannot decode username/password struct", http.StatusBadRequest)
		return
	}
	userFromDb, err := c.userRepo.FindByUserName(userDto.Username)
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

func (c *AuthController) isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
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

	userFromDb, err := c.userRepo.FindByUserName(username)
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
	user, err := c.userRepo.FindByUserName(username)
	if err != nil {
		http.Error(writer, "Unknown user", http.StatusBadRequest)
		return
	}

	accessToken, err := service.CreateToken(userKey, username, user.Role)
	if err != nil {
		http.Error(writer, "Cannot create accessToken", http.StatusInternalServerError)
		return
	}

	ok = c.userRepo.UpdateAccessToken(username, accessToken)
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
		tokenString, ok := service.ParseToken(header)
		if !ok {
			http.Error(writer, tokenString, http.StatusBadRequest)
			return
		}
		token, ok := service.CheckToken(tokenString)
		log.Println("bearerToken " + tokenString)

		if !ok {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !c.setContext(writer, request, ok, token, tokenString) {
			return
		}

		log.Println("success")

		next(writer, request)
	}
}

func (c *AuthController) setContext(writer http.ResponseWriter, request *http.Request, ok bool, token *jwt.Token, tokenString string) bool {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return false
		}
		userId, ok := claims["userId"].(string)
		if !ok {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return false
		}

		_, err := c.userRepo.FindByUserName(username)
		if err != nil {
			http.Error(writer, "Unauthorized, user not exists", http.StatusUnauthorized)
			return false
		}

		context.Set(request, "username", username)
		context.Set(request, "userId", userId)
		context.Set(request, "token", tokenString)
		return true
	}
	return false
}

func (c *AuthController) CheckToken(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	log.Println("check token handler")
	header := request.Header.Get("Authorization")
	tokenString, ok := service.ParseToken(header)
	if !ok {
		http.Error(writer, tokenString, http.StatusBadRequest)
		return
	}

	token, ok := service.CheckToken(tokenString)
	log.Println("bearerToken " + tokenString)

	if !ok {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username, userId, ok := c.parsePayload(ok, token)
	if !ok {
		http.Error(writer, username, http.StatusUnauthorized)
		return
	}

	user, err := c.userRepo.FindById(userId)
	if err != nil {
		http.Error(writer, username, http.StatusUnauthorized)
		return
	}

	log.Println("success")

	c.BaseController.JsonResponse(writer, commonDto.UserDto{Username: username, UserId: userId, Ok: true, Role: user.Role.Values()})
}

func (c *AuthController) parsePayload(ok bool, token *jwt.Token) (string, string, bool) {
	var username string
	var userId string
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		username, ok = claims["username"].(string)
		if !ok {
			return "Unauthorized, username not exist", "", false
		}
		userId, ok = claims["userId"].(string)
		if !ok {
			return "Unauthorized, userId not exist", "", false
		}

		_, err := c.userRepo.FindByUserName(username)
		if err != nil {
			return "Unauthorized, user not exists", "", false
		}
	}
	return username, userId, true
}
