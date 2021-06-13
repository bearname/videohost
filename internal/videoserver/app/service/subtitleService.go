package service

import (
	"encoding/json"
	"github.com/bearname/videohost/internal/common/caching"
	"github.com/bearname/videohost/internal/videoserver/domain/dto"
	"github.com/bearname/videohost/internal/videoserver/domain/model"
	"github.com/bearname/videohost/internal/videoserver/domain/repository"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type SubtitleService struct {
	playlistRepo repository.SubtitleRepository
	cache        caching.Cache
}

const cacheSubtitlePrefix = "sub-"

func NewSubtitleService(playlistRepo repository.SubtitleRepository, cache caching.Cache) *SubtitleService {
	s := new(SubtitleService)
	s.playlistRepo = playlistRepo
	s.cache = cache
	return s
}

func (s *SubtitleService) Create(subtitle dto.CreateSubtitleRequestDto) (int64, error) {
	id, err := s.playlistRepo.Create(subtitle)
	if err != nil {
		return 0, err
	}

	itoa := strconv.Itoa(int(id))
	key := cacheSubtitlePrefix + itoa

	marshal, err := json.Marshal(subtitle)
	if err == nil {
		_ = s.cache.Set(key, string(marshal))
		log.Error("failed save to cache")
	}
	return id, nil
}

func (s *SubtitleService) FindByVideo(videoId int) (model.Subtitle, error) {
	panic("implement subService.FindByVideo")
}

func (s *SubtitleService) Delete(subtitleId int) error {
	panic("implement subService.Delete")
}
