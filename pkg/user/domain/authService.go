package domain

import (
	commonDto "github.com/bearname/videohost/pkg/common/dto"
	"github.com/bearname/videohost/pkg/common/util"
	"io"
	"net/http"
)

type AuthService interface {
	CreateUser(requestBody io.ReadCloser) (util.Token, error, int)
	Login(requestBody io.ReadCloser) (util.Token, error, int)
	ValidateToken(authorizationHeader string) (commonDto.UserDto, int)
	RefreshToken(request *http.Request) (util.Token, error, int)
}
