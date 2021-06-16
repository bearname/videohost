package service

import (
	commonDto "github.com/bearname/videohost/internal/common/dto"
	"github.com/bearname/videohost/internal/common/util"
	"github.com/bearname/videohost/internal/user/app/dto"
	"github.com/bearname/videohost/internal/user/domain"
	"github.com/bearname/videohost/internal/user/domain/model"
	"github.com/bearname/videohost/internal/user/domain/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthServiceImpl struct {
	userRepo repository.UserRepo
}

func NewAuthService(userRepository repository.UserRepo) *AuthServiceImpl {
	v := new(AuthServiceImpl)
	v.userRepo = userRepository
	return v
}

func (a *AuthServiceImpl) CreateUser(newUser dto.SignupUserDto) (util.Token, error) {

	userFromDb, err := a.userRepo.FindByUserName(newUser.Username)
	if (err == nil && userFromDb.Username == newUser.Username) || (err != nil && err.Error() != "sql: no rows in result set") {
		log.Error(err.Error())
		return util.Token{}, domain.ErrDuplicateUser
	}
	log.Println(userFromDb.Username == newUser.Username)

	userKey, err := uuid.NewUUID()
	if err != nil {
		log.Error(err.Error())
		return util.Token{}, domain.ErrFailedCreateUserID
	}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	role := model.General
	accessToken, err := CreateToken(userKey.String(), newUser.Username, role)
	if err != nil {
		log.Error(err.Error())
		return util.Token{}, domain.ErrFailedCreateAccessToken
	}

	refreshToken, err := CreateTokenWithDuration(userKey.String(), newUser.Username, role, time.Hour*24*365*10)
	if err != nil {
		log.Error(err.Error())
		return util.Token{}, domain.ErrFailedUpdateAccessToken
	}

	err = a.userRepo.CreateUser(userKey.String(), newUser.Username, passwordHash, newUser.Email, newUser.IsSubscribed, role, accessToken, refreshToken)
	if err != nil {
		log.Error(err.Error())
		return util.Token{}, domain.ErrFailedSaveUser
	}

	return util.Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (a *AuthServiceImpl) Login(userDto dto.LoginUserDto) (util.Token, error) {

	userFromDb, err := a.userRepo.FindByUserName(userDto.Username)
	if (err == nil && userFromDb.Username != userDto.Username) || err != nil {
		log.Error(err.Error())
		return util.Token{}, domain.ErrUserNotExist
	}

	err = bcrypt.CompareHashAndPassword(userFromDb.Password, []byte(userDto.Password))
	if err != nil {
		log.Error(err.Error())
		return util.Token{}, domain.ErrWrongPassword
	}

	role := model.General
	accessToken, err := CreateToken(userFromDb.Key, userFromDb.Username, role)
	if err != nil {
		log.Error(err.Error())
		return util.Token{}, domain.ErrFailedCreateAccessToken
	}

	refreshToken, err := CreateTokenWithDuration(userFromDb.Key, userFromDb.Username, role, time.Hour*24*365*10)
	if err != nil {
		log.Error(err.Error())
		return util.Token{}, domain.ErrFailedUpdateAccessToken
	}

	return util.Token{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (a *AuthServiceImpl) ValidateToken(authorizationHeader string) (commonDto.UserDto, error) {
	tokenString, ok := ParseToken(authorizationHeader)
	if !ok {
		return commonDto.UserDto{}, domain.ErrInvalidAuthorizationHeader
	}

	token, ok := CheckToken(tokenString)
	log.Println("bearerToken " + tokenString)

	if !ok {
		return commonDto.UserDto{}, domain.ErrInvalidAccessToken
	}

	username, userId, ok := a.parsePayload(token)
	if !ok {
		return commonDto.UserDto{}, domain.ErrInvalidAccessToken
	}

	user, err := a.userRepo.FindById(userId)
	if err != nil {
		return commonDto.UserDto{}, domain.ErrUserNotExist
	}

	return commonDto.UserDto{Username: username,
		UserId: userId,
		Ok:     true,
		Role:   user.Role.Values(),
		Token:  tokenString}, nil
}

func (a *AuthServiceImpl) RefreshToken(refreshTokenDto dto.RefreshTokenDto) (util.Token, error) {
	username := refreshTokenDto.Username
	userFromDb, err := a.userRepo.FindByUserName(username)
	if (err == nil && userFromDb.Username != username) || err != nil {
		return util.Token{}, domain.ErrUserNotExist
	}

	if userFromDb.RefreshToken != refreshTokenDto.Token {
		return util.Token{}, domain.ErrInvalidRefreshToken
	}

	accessToken, err := CreateToken(refreshTokenDto.UserId, username, userFromDb.Role)
	if err != nil {
		return util.Token{}, domain.ErrFailedCreateAccessToken
	}

	ok := a.userRepo.UpdateAccessToken(username, accessToken)
	if !ok {
		return util.Token{}, domain.ErrFailedUpdateAccessToken
	}

	return util.Token{AccessToken: accessToken, RefreshToken: userFromDb.RefreshToken}, nil
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
