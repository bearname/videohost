package service

import (
	"encoding/json"
	"errors"
	"github.com/bearname/videohost/pkg/common/amqp"
	"github.com/bearname/videohost/pkg/common/caching"
	commonDto "github.com/bearname/videohost/pkg/common/dto"
	"github.com/bearname/videohost/pkg/common/util"
	userModel "github.com/bearname/videohost/pkg/user/domain/model"
	scaleModel "github.com/bearname/videohost/pkg/video-scaler/domain"
	"github.com/bearname/videohost/pkg/videoserver/app/dto"
	"github.com/bearname/videohost/pkg/videoserver/domain"
	"github.com/bearname/videohost/pkg/videoserver/domain/model"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/ftp"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strconv"
)

type VideoServiceImpl struct {
	videoRepo     domain.VideoRepository
	messageBroker amqp.MessageBroker
	cache         caching.Cache
}

const videoCachePrefix = "video-"

func NewVideoService(videoRepository domain.VideoRepository, messageBroker amqp.MessageBroker, cache caching.Cache) *VideoServiceImpl {
	s := new(VideoServiceImpl)

	s.videoRepo = videoRepository
	s.messageBroker = messageBroker
	s.cache = cache
	return s
}

func (s *VideoServiceImpl) FindVideo(videoId string) (*model.Video, error) {
	video, err := s.readFromCache(videoId)
	if err == nil {
		return video, nil
	}

	video, err = s.videoRepo.Find(videoId)
	if err != nil {
		return &model.Video{}, err
	}

	return video, nil
}

func (s *VideoServiceImpl) UploadVideo(userId string, videoDto *dto.UploadVideoDto) (uuid.UUID, error) {

	contentType := videoDto.FileHeader.Header.Get("Content-Type")
	if contentType != util.VideoContentType {
		return uuid.UUID{}, errors.New("unexpected content type")
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return uuid.UUID{}, err
	}
	videoId := id.String()
	connection := ftp.NewFtpConnection("localhost:21", "user", "123")
	err = connection.CopyFile(videoId, videoDto.MultipartFile)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = s.videoRepo.Create(
		userId,
		videoId,
		videoDto.Title,
		videoDto.Description,
		filepath.Join(util.ContentDir, videoId, util.VideoFileName),
		videoDto.Chapters,
	)
	if err != nil {
		err = connection.RemoveDirRecur(videoId)
		err = connection.RemoveDir(videoId)

		return uuid.UUID{}, err
	}

	find, err := s.videoRepo.Find(videoId)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = s.writeToCache(videoId, find)

	err = s.messageBroker.Publish("events_topic", "events.video-uploaded", videoId)
	if err != nil {
		log.Error("Failed publish event 'video-uploaded")
	}
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (s *VideoServiceImpl) readFromCache(videoId string) (*model.Video, error) {
	cacheStr, err := s.cache.Get(videoCachePrefix + videoId)
	if err != nil {
		return nil, err
	}
	var cacheVideo model.Video

	err = json.Unmarshal([]byte(cacheStr), &cacheVideo)
	if err != nil {
		return &cacheVideo, err
	}

	return &cacheVideo, nil
}

func (s *VideoServiceImpl) writeToCache(videoId string, video *model.Video) error {
	if !s.cache.IsOk() {
		return caching.ErrCacheUnavailable
	}
	cacheByte, err := json.Marshal(video)
	if err != nil {
		return err
	}

	return s.cache.Set(videoCachePrefix+videoId, string(cacheByte))
}

func (s *VideoServiceImpl) UpdateTitleAndDescription(userDto commonDto.UserDto, videoId string, videoDto dto.VideoMetadata) error {
	if len(userDto.UserId) == 0 || len(videoId) == 0 || len(videoDto.Title) == 0 || len(videoDto.Description) == 0 {
		return errors.New("parameter must be length more than 0")
	}

	video, err := s.checkOwner(userDto, videoId)
	if err != nil {
		return err
	}

	video.Name = videoDto.Title
	video.Description = videoDto.Description

	err = s.videoRepo.Update(videoId, videoDto.Title, videoDto.Description)
	if err == nil {
		err = s.writeToCache(videoId, video)
	}

	return err
}

func (s *VideoServiceImpl) AddQuality(videoId string, userDto commonDto.UserDto, quality model.Quality) error {
	video, err := s.checkOwner(userDto, videoId)
	if err != nil {
		return err
	}

	if !scaleModel.IsSupportedQuality(quality.Value) {
		return errors.New("unsupported quality")
	}
	err = s.videoRepo.AddVideoQuality(videoId, strconv.Itoa(quality.Value))
	if err != nil {
		return err
	}
	if len(video.Quality) != 0 {
		video.Quality += ","
	}

	return s.writeToCache(videoId, video)
}

func (s *VideoServiceImpl) DeleteVideo(userDto commonDto.UserDto, videoId string) error {

	video, err := s.checkOwner(userDto, videoId)
	if err != nil {
		return err
	}

	connection := ftp.NewFtpConnection("localhost:21", "user", "123")
	if connection == nil {
		return errors.New("failed connect to video store server")
	}
	err = connection.RemoveDirRecur(videoId)
	err = connection.RemoveDir(videoId)

	if err != nil {
		return s.videoRepo.Save(video)
	}
	err = s.videoRepo.Delete(videoId)
	if err != nil {
		return s.videoRepo.Save(video)
	}
	_, err = s.cache.Get(videoId)
	if err == nil {
		err := s.cache.Del(videoCachePrefix + videoId)
		if err != nil {
			return err
		}
	}
	return err
}

func (s *VideoServiceImpl) checkOwner(userDto commonDto.UserDto, videoId string) (*model.Video, error) {
	if len(videoId) == 0 {
		return &model.Video{}, errors.New("videoId must be length more than 0")
	}

	video, err := s.videoRepo.Find(videoId)
	if err != nil {
		return nil, err
	}

	if userDto.Role == userModel.Admin.Values() {
		return video, nil
	}

	if userDto.UserId != video.OwnerId {
		return nil, errors.New("user with id " + userDto.UserId + " not owner of video with id " + videoId)
	}

	return video, err
}

func (s *VideoServiceImpl) FindVideoOnPage(searchDto *dto.SearchDto) (dto.SearchResultDto, error) {
	pageCount, ok := s.videoRepo.GetPageCount(searchDto.Count)
	if !ok {
		return dto.SearchResultDto{}, errors.New("failed get page count")
	}

	videos, err := s.videoRepo.FindVideosByPage(searchDto.Page, searchDto.Count)
	if err != nil {
		return dto.SearchResultDto{}, err
	}

	return dto.SearchResultDto{PageCount: pageCount, Videos: videos}, nil
}
