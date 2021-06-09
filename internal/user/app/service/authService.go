package service

import (
	"encoding/json"
	"errors"
	commonDto "github.com/bearname/videohost/internal/common/dto"
	"github.com/bearname/videohost/internal/user/app/dto"
	"github.com/bearname/videohost/internal/user/domain/model"
	"github.com/bearname/videohost/internal/user/domain/repository"
	"github.com/bearname/videohost/internal/video-scaler/domain"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type AuthServiceImpl struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) *AuthServiceImpl {
	v := new(AuthServiceImpl)
	v.userRepo = userRepository
	return v
}

func (a *AuthServiceImpl) CreateUser(requestBody io.ReadCloser) (domain.Token, error, int) {
	var newUser dto.SignupUserDto
	err := json.NewDecoder(requestBody).Decode(&newUser)
	if err != nil {
		log.Error(err.Error())
		return domain.Token{}, errors.New("cannot decode request"), http.StatusBadRequest
	}

	if !a.isEmailValid(newUser.Email) {
		return domain.Token{}, errors.New("user email not valid"), http.StatusBadRequest
	}

	userFromDb, err := a.userRepo.FindByUserName(newUser.Username)
	if (err == nil && userFromDb.Username == newUser.Username) || (err != nil && err.Error() != "sql: no rows in result set") {
		log.Error(err.Error())
		return domain.Token{}, errors.New("user already exists"), http.StatusBadRequest
	}
	log.Println(userFromDb.Username == newUser.Username)

	userKey, err := uuid.NewUUID()
	if err != nil {
		log.Error(err.Error())
		return domain.Token{}, errors.New("failed generate id"), http.StatusInternalServerError
	}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	role := model.General
	accessToken, err := CreateToken(userKey.String(), newUser.Username, role)
	if err != nil {
		log.Error(err.Error())
		return domain.Token{}, errors.New("cannot create accessToken"), http.StatusInternalServerError
	}

	refreshToken, err := CreateTokenWithDuration(userKey.String(), newUser.Username, role, time.Hour*24*365*10)
	if err != nil {
		log.Error(err.Error())
		return domain.Token{}, errors.New("cannot create refreshToken"), http.StatusInternalServerError
	}

	err = a.userRepo.CreateUser(userKey.String(), newUser.Username, passwordHash, newUser.Email, newUser.IsSubscribed, role, accessToken, refreshToken)
	if err != nil {
		log.Error(err.Error())
		return domain.Token{}, errors.New("failed save user"), http.StatusInternalServerError
	}

	return domain.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil, http.StatusOK
}

func (a *AuthServiceImpl) Login(requestBody io.ReadCloser) (domain.Token, error, int) {
	var userDto dto.LoginUserDto
	err := json.NewDecoder(requestBody).Decode(&userDto)
	if err != nil {
		log.Error(err.Error())
		return domain.Token{}, errors.New("cannot decode username/password struct"), http.StatusBadRequest
	}
	userFromDb, err := a.userRepo.FindByUserName(userDto.Username)
	if (err == nil && userFromDb.Username != userDto.Username) || err != nil {
		log.Error(err.Error())
		return domain.Token{}, errors.New("user not exist"), http.StatusBadRequest
	}

	err = bcrypt.CompareHashAndPassword(userFromDb.Password, []byte(userDto.Password))
	if err != nil {
		log.Error(err.Error())
		return domain.Token{}, errors.New("wrong password"), http.StatusUnauthorized
	}

	return domain.Token{
		AccessToken:  userFromDb.AccessToken,
		RefreshToken: userFromDb.RefreshToken,
	}, nil, http.StatusOK
}

func (a *AuthServiceImpl) ValidateToken(authorizationHeader string) (commonDto.UserDto, int) {
	tokenString, ok := ParseToken(authorizationHeader)
	if !ok {
		return commonDto.UserDto{}, http.StatusBadRequest
	}

	token, ok := CheckToken(tokenString)
	log.Println("bearerToken " + tokenString)

	if !ok {
		return commonDto.UserDto{}, http.StatusUnauthorized
	}

	username, userId, ok := a.parsePayload(token)
	if !ok {
		return commonDto.UserDto{}, http.StatusUnauthorized
	}

	user, err := a.userRepo.FindById(userId)
	if err != nil {
		return commonDto.UserDto{}, http.StatusUnauthorized
	}

	return commonDto.UserDto{Username: username,
		UserId: userId,
		Ok:     true,
		Role:   user.Role.Values(),
		Token:  tokenString}, http.StatusOK
}

func (a *AuthServiceImpl) RefreshToken(request *http.Request) (domain.Token, error, int) {
	username, ok := context.Get(request, "username").(string)
	if !ok {
		return domain.Token{}, errors.New("cannot check username"), http.StatusInternalServerError
	}
	userKey, ok := context.Get(request, "userId").(string)
	if !ok {
		return domain.Token{}, errors.New("cannot check userId"), http.StatusInternalServerError
	}

	userFromDb, err := a.userRepo.FindByUserName(username)
	if (err == nil && userFromDb.Username != username) || err != nil {
		return domain.Token{}, errors.New("unauthorized, user not exists"), http.StatusUnauthorized
	}
	token, ok := context.Get(request, "token").(string)
	if !ok {
		return domain.Token{}, errors.New("token to preset on context by token checker"), http.StatusInternalServerError
	}

	if userFromDb.RefreshToken != token {
		return domain.Token{}, errors.New("invalid Refresh token"), http.StatusInternalServerError
	}
	user, err := a.userRepo.FindByUserName(username)
	if err != nil {
		return domain.Token{}, errors.New("unknown user"), http.StatusBadRequest
	}

	accessToken, err := CreateToken(userKey, username, user.Role)
	if err != nil {
		return domain.Token{}, errors.New("cannot create accessToken"), http.StatusInternalServerError
	}

	ok = a.userRepo.UpdateAccessToken(username, accessToken)
	if !ok {
		return domain.Token{}, errors.New("failed update accessToken"), http.StatusInternalServerError
	}

	return domain.Token{
		AccessToken:  accessToken,
		RefreshToken: userFromDb.RefreshToken,
	}, nil, http.StatusOK
}

func (a *AuthServiceImpl) parsePayload(token *jwt.Token) (string, string, bool) {
	var username string
	var userId string
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		username, ok = claims["username"].(string)
		if !ok {
			return "unauthorized, username not exist", "", false
		}
		userId, ok = claims["userId"].(string)
		if !ok {
			return "unauthorized, userId not exist", "", false
		}

		_, err := a.userRepo.FindByUserName(username)
		if err != nil {
			return "unauthorized, user not exists", "", false
		}
	}

	return username, userId, true
}

func (a *AuthServiceImpl) isEmailValid(e string) bool {
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
