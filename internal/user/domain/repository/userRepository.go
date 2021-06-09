package repository

import (
	"github.com/bearname/videohost/internal/user/domain/model"
)

type UserRepo interface {
	CreateUser(key string, username string, password []byte, email string, isSubscribed bool, role model.Role, accessToken string, refreshToken string) error
	FindById(userId string) (model.User, error)
	FindByUserName(username string) (model.User, error)
	UpdatePassword(username string, password []byte) bool
	UpdateAccessToken(username string, token string) bool
	UpdateRefreshToken(username string, token string) bool
	GetCountVideos(userId string) (int, bool)
}
