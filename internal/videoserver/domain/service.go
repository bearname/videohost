package domain

import (
	"github.com/bearname/videohost/internal/common/db"
	commonDto "github.com/bearname/videohost/internal/common/dto"
	"github.com/bearname/videohost/internal/videoserver/domain/dto"
	"github.com/bearname/videohost/internal/videoserver/domain/model"
	"github.com/google/uuid"
	"net/http"
)

type StreamService interface {
	ServeHlsM3u8(w http.ResponseWriter, r *http.Request, videoId string, m3u8Name string)
	ServeHlsTs(w http.ResponseWriter, r *http.Request, segName, videoId string)
}

type VideoService interface {
	FindVideo(videoId string) (*model.Video, error)
	UploadVideo(userId string, videoDto *dto.UploadVideoDto) (uuid.UUID, error)
	UpdateTitleAndDescription(userDto commonDto.UserDto, videoId string, videoDto dto.VideoMetadata) error
	AddQuality(videoId string, userDto commonDto.UserDto, quality model.Quality) error
	DeleteVideo(userDto commonDto.UserDto, videoId string) error
	FindVideoOnPage(searchDto *dto.SearchDto) (dto.SearchResultDto, error)
	LikeVideo(like model.Like) (model.Action, error)
	FindUserLikedVideo(userId string, page db.Page) ([]model.VideoListItem, error)
}

type Action int

const (
	AddVideoAction Action = iota
	RemoveVideoAction
	ReorderVideoAction
)

func (a *Action) String() string {
	switch *a {
	case AddVideoAction:
		return "add"
	case RemoveVideoAction:
		return "remove"
	case ReorderVideoAction:
		return "remove"
	default:
		return ""
	}
}

type PlayListService interface {
	CreatePlaylist(playlist dto.CreatePlaylistDto) (int64, error)
	FindPlaylist(playlistId int) (model.Playlist, error)
	ModifyVideosOnPlaylist(playlistId int, userId string, videosId []string, action Action) error
	ChangeOrder(playlistId int, videoId []string, order []int) error
	ChangePrivacy(id string, playlistId int, privacyType model.PrivacyType) error
	Delete(ownerId string, playlistId int) error
}
