package repository

import (
	"github.com/bearname/videohost/internal/videoserver/domain/dto"
	"github.com/bearname/videohost/internal/videoserver/domain/model"
)

type PlaylistRepository interface {
	Create(playlist dto.CreatePlaylistDto) (int64, error)
	FindPlaylist(playlistId int) (model.Playlist, error)
	AddVideos(playlistId int, userId string, videosId []string) error
	RemoveVideos(playlistId int, userId string, videosId []string) error
	ChangeOrder(playlistId int, videoId []string, order []int) error
	ChangePrivacy(id string, playlistId int, privacyType model.PrivacyType) error
	Delete(ownerId string, playlistId int) error
}
