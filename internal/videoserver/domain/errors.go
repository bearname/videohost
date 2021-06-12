package domain

import "errors"

var (
	ErrVideoNotFound             = errors.New("video not found")
	ErrFailedDeleteLike          = errors.New("failed delete like")
	ErrFailedAddLike             = errors.New("failed add like")
	ErrInternal                  = errors.New("internal error")
	ErrAlreadyLike               = errors.New("already like")
	ErrAlreadyDisLike            = errors.New("already dislike")
	ErrPlaylistNotFound          = errors.New("not found")
	ErrPlaylistDuplicate         = errors.New("playlist duplicate video in playlist")
	ErrUnknownModificationAction = errors.New("unknown playlist modification action")
)
