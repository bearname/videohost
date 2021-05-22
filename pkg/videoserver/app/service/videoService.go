package service

import (
	"errors"
	"github.com/bearname/videohost/pkg/common/amqp"
	"github.com/bearname/videohost/pkg/common/util"
	"github.com/bearname/videohost/pkg/videoserver/app/dto"
	"github.com/bearname/videohost/pkg/videoserver/domain/repository"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/ftp"
	"github.com/google/uuid"
	"mime/multipart"
	"path/filepath"
)

type VideoService struct {
	videoRepository repository.VideoRepository
	messageBroker   amqp.MessageBroker
}

func NewVideoService(videoRepository repository.VideoRepository, messageBroker amqp.MessageBroker) *VideoService {
	v := new(VideoService)

	v.videoRepository = videoRepository
	v.messageBroker = messageBroker
	return v
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

	err = s.videoRepository.Create(
		userId,
		videoId,
		title,
		description,
		filepath.Join(util.ContentDir, videoId, util.VideoFileName),
	)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (s *VideoService) UpdateTitleAndDescription(userId string, videoId string, videoDto dto.VideoMetadata) error {
	if len(userId) == 0 || len(videoId) == 0 || len(videoDto.Title) == 0 || len(videoDto.Description) == 0 {
		return errors.New("Parameter must be length more than 0")
	}
	video, err := s.videoRepository.Find(videoId)
	if err != nil {
		return err
	}
	if userId != video.OwnerId {
		return errors.New("user with id " + userId + " not owner of video with id " + videoId)
	}

	return s.videoRepository.Update(videoId, videoDto.Title, videoDto.Description)
}

func (s *VideoService) DeleteVideo(userId string, videoId string) error {
	if len(videoId) == 0 {
		return errors.New("videoId must be length more than 0")
	}

	video, err := s.videoRepository.Find(videoId)
	if err != nil {
		return err
	}

	if userId != video.OwnerId {
		return errors.New("user with id " + userId + " not owner of video with id " + videoId)
	}

	connection := ftp.NewFtpConnection("localhost:21", "user", "123")
	if connection == nil {
		return errors.New("Failed connect to video store server")
	}
	err = connection.RemoveDirRecur(videoId)
	err = connection.RemoveDir(videoId)

	if err != nil {
		return s.videoRepository.Save(video)
	}
	return s.videoRepository.Delete(videoId)
}
