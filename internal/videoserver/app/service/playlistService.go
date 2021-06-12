package service

import (
	"github.com/bearname/videohost/internal/common/caching"
	"github.com/bearname/videohost/internal/videoserver/domain"
	"github.com/bearname/videohost/internal/videoserver/domain/dto"
	"github.com/bearname/videohost/internal/videoserver/domain/model"
	"github.com/bearname/videohost/internal/videoserver/domain/repository"
)

type PlayListService struct {
	playlistRepo repository.PlaylistRepository
	cache        caching.Cache
}

const cachePrefix = "playlist-"

func NewPlayListService(playlistRepo repository.PlaylistRepository, cache caching.Cache) *PlayListService {
	s := new(PlayListService)

	s.playlistRepo = playlistRepo
	s.cache = cache
	return s
}

func (s *PlayListService) CreatePlaylist(playlist dto.CreatePlaylistDto) (int64, error) {
	id, err := s.playlistRepo.Create(playlist)
	if err != nil {
		return 0, err
	}
	//TODO add caching
	return id, nil
}

func (s *PlayListService) FindPlaylist(playlistId int) (model.Playlist, error) {
	return s.playlistRepo.FindPlaylist(playlistId)
}

func (s *PlayListService) FindUserPlaylists(userId string, privacyType []model.PrivacyType) ([]dto.PlaylistItem, error) {
	return s.playlistRepo.FindPlaylists(userId, privacyType)
}

func (s *PlayListService) ModifyVideosOnPlaylist(playlistId int, userId string, videosId []string, action domain.Action) error {
	switch action {
	case domain.AddVideoAction:
		return s.playlistRepo.AddVideos(playlistId, userId, videosId)
	case domain.RemoveVideoAction:
		return s.playlistRepo.RemoveVideos(playlistId, userId, videosId)
	default:
		return domain.ErrUnknownModificationAction
	}
}

func (s *PlayListService) ChangeOrder(playlistId int, videoId []string, order []int) error {
	panic("implement ChangeOrder")
}

func (s *PlayListService) ChangePrivacy(id string, playlistId int, privacyType model.PrivacyType) error {
	return s.playlistRepo.ChangePrivacy(id, playlistId, privacyType)
}

func (s *PlayListService) Delete(ownerId string, playlistId int) error {
	return s.playlistRepo.Delete(ownerId, playlistId)
}
