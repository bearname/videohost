package repository

import (
	model2 "github.com/bearname/videohost/pkg/videoserver/domain/model"
)

type UserRepository interface {
	FindByUserName(username string) (*model2.User, error)
	CreateUser(key string, username string, password []byte, role model2.Role, accessToken string, refreshToken string) error
	UpdatePassword(username string, password []byte) bool
	UpdateAccessToken(username string, token string) bool
	UpdateRefreshToken(username string, token string) bool
	GetCountVideos(userId string) (int, bool)
}
