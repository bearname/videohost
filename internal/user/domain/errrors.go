package domain

import "errors"

var (
	ErrUserNotExist               = errors.New("unauthorized, user not exists")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	ErrInvalidAccessToken         = errors.New("invalid access token")
	ErrInvalidRefreshToken        = errors.New("invalid refresh token")
	ErrFailedCreateAccessToken    = errors.New("failed create accessToken")
	ErrFailedUpdateAccessToken    = errors.New("failed update accessToken")
	ErrFailedSaveUser             = errors.New("failed save user")
	ErrFailedCreateUserID         = errors.New("failed generate id")
	ErrDuplicateUser              = errors.New("duplicate user")
	ErrWrongPassword              = errors.New("wrong password")
)
