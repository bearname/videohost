package domain

import (
	"github.com/bearname/videohost/pkg/user/app/dto"
	"net/http"
)

type UserService interface {
	Find(usernameOrId string) (dto.FindUserDto, error)
	UpdatePassword(request *http.Request) error
}
