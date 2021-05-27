package service

import (
	"encoding/json"
	"errors"
	"github.com/bearname/videohost/pkg/common/amqp"
	"github.com/bearname/videohost/pkg/common/caching"
	commonDto "github.com/bearname/videohost/pkg/common/dto"
	"github.com/bearname/videohost/pkg/common/util"
	model2 "github.com/bearname/videohost/pkg/user/domain/model"
	"github.com/bearname/videohost/pkg/video-scaler/domain"
	"github.com/bearname/videohost/pkg/videoserver/app/dto"
	"github.com/bearname/videohost/pkg/videoserver/domain/model"
	"github.com/bearname/videohost/pkg/videoserver/domain/repository"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/ftp"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"mime/multipart"
	"path/filepath"
	"strconv"
)

type VideoService struct {
	videoRepo     repository.VideoRepository
	messageBroker amqp.MessageBroker
	cache         caching.Cache
}

const videoCachePrefix = "video-"

func NewVideoService(videoRepository repository.VideoRepository, messageBroker amqp.MessageBroker, cache caching.Cache) *VideoService {
	s := new(VideoService)

	s.videoRepo = videoRepository
	s.messageBroker = messageBroker
	s.cache = cache
	return s
}

func (s *VideoService) FindVideo(videoId string) (*model.Video, error) {
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

func (s *VideoService) UploadVideo(userId string, title string, description string, fileReader multipart.File, header *multipart.FileHeader) (uuid.UUID, error) {
	contentType := header.Header.Get("Content-Type")
	if contentType != util.VideoContentType {
		return uuid.UUID{}, errors.New("Unexpected content type")
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return uuid.UUID{}, err
	}
	videoId := id.String()
	connection := ftp.NewFtpConnection("localhost:21", "user", "123")
	err = connection.CopyFile(videoId, fileReader)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = s.videoRepo.Create(
		userId,
		videoId,
		title,
		description,
		filepath.Join(util.ContentDir, videoId, util.VideoFileName),
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = s.messageBroker.Publish("events_topic", "events.video-uploaded", videoId)
	if err != nil {
		log.Error("Failed publish event 'video-uploaded")
	}

	find, err := s.videoRepo.Find(videoId)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = s.writeToCache(videoId, find)

	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (s *VideoService) readFromCache(videoId string) (*model.Video, error) {
	cacheStr, err := s.cache.Get(videoCachePrefix + videoId)
	var cacheVideo model.Video
	if err != nil {
		return &cacheVideo, err
	}

	err = json.Unmarshal([]byte(cacheStr), &cacheVideo)
	if err != nil {
		return &cacheVideo, err
	}

	return &cacheVideo, nil
}

func (s *VideoService) writeToCache(videoId string, video *model.Video) error {
	cacheByte, err := json.Marshal(video)
	if err != nil {
		return err
	}

	return s.cache.Set(videoCachePrefix+videoId, string(cacheByte))
}

func (s *VideoService) UpdateTitleAndDescription(userDto commonDto.UserDto, videoId string, videoDto dto.VideoMetadata) error {
	if len(userDto.UserId) == 0 || len(videoId) == 0 || len(videoDto.Title) == 0 || len(videoDto.Description) == 0 {
		return errors.New("Parameter must be length more than 0")
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

func (s *VideoService) AddQuality(videoId string, userDto commonDto.UserDto, quality dto.Quality) error {
	video, err := s.checkOwner(userDto, videoId)
	if err != nil {
		return err
	}

	if !domain.IsSupportedQuality(quality.Value) {
		return errors.New("Unsupported quality")
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

func (s *VideoService) DeleteVideo(userDto commonDto.UserDto, videoId string) error {

	video, err := s.checkOwner(userDto, videoId)
	if err != nil {
		return err
	}

	connection := ftp.NewFtpConnection("localhost:21", "user", "123")
	if connection == nil {
		return errors.New("Failed connect to video store server")
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

func (s *VideoService) checkOwner(userDto commonDto.UserDto, videoId string) (*model.Video, error) {
	if len(videoId) == 0 {
		return &model.Video{}, errors.New("videoId must be length more than 0")
	}

	video, err := s.videoRepo.Find(videoId)
	if err != nil {
		return nil, err
	}

	if userDto.Role == model2.Admin.Values() {
		return video, nil
	}

	if userDto.UserId != video.OwnerId {
		return nil, errors.New("user with id " + userDto.UserId + " not owner of video with id " + videoId)
	}

	return video, err
}
