package repository

import (
	"github.com/bearname/videohost/videoserver/model"
)

type UserRepository interface {
	FindByUserName(username string) (*model.User, error)
	CreateUser(key string, username string, password []byte, role model.Role, accessToken string, refreshToken string) error
	UpdatePassword(username string, password []byte) bool
	UpdateAccessToken(username string, token string) bool
	UpdateRefreshToken(username string, token string) bool
	GetCountVideos(userId string) (int, bool)
}
