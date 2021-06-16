package domain

import (
	commonDto "github.com/bearname/videohost/internal/common/dto"
	"github.com/bearname/videohost/internal/common/util"
	"github.com/bearname/videohost/internal/user/app/dto"
)

type AuthService interface {
	CreateUser(newUserDto dto.SignupUserDto) (util.Token, error)
	Login(loginUserDto dto.LoginUserDto) (util.Token, error)
	ValidateToken(authorizationHeader string) (commonDto.UserDto, error)
	RefreshToken(refreshTokenDto dto.RefreshTokenDto) (util.Token, error)
}
