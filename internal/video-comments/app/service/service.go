package service

import (
	"encoding/json"
	"errors"
	"github.com/bearname/videohost/internal/common/caching"
	"github.com/bearname/videohost/internal/common/util"
	"github.com/bearname/videohost/internal/video-comments/domain"
	"github.com/bearname/videohost/internal/videoserver/domain/model"
	"net/http"
)

type CommentService struct {
	videoRepo          domain.CommentRepository
	cache              caching.Cache
	videoServerAddress string
}

func NewCommentService(repo domain.CommentRepository, cache caching.Cache, videoServer string) *CommentService {
	c := new(CommentService)
	c.videoRepo = repo
	c.cache = cache
	c.videoServerAddress = videoServer
	return c
}

func (s *CommentService) Create(commentDto domain.CommentDto) (int64, error) {
	//TODO send needed field on get parameter
	responseBody, err := util.GetRequest(&http.Client{}, s.videoServerAddress+"/api/v1/videos/"+commentDto.VideoId, "")
	if err != nil {
		return 0, err
	}
	var video model.Video
	err = json.Unmarshal(responseBody, &video)
	if err != nil {
		return 0, err
	}
	return s.videoRepo.Create(commentDto.UserId, commentDto.VideoId, commentDto.ParentId, commentDto.Message)
}

func (s *CommentService) FindRootLevel(videoId string, page *domain.Page) (domain.VideoComments, error) {
	return s.videoRepo.FindRootLevel(videoId, page)
}

func (s *CommentService) FindChildren(rootCommentId int, page *domain.Page) ([]domain.Comment, error) {
	return s.videoRepo.FindChildren(rootCommentId, page)
}

func (s *CommentService) FindUserComments(userId string, page *domain.Page) (domain.Comments, error) {
	return s.videoRepo.FindUserComments(userId, page)
}

func (s *CommentService) Edit(commentId int, message string, userId string) error {
	err := s.isUserOwner(commentId, userId)
	if err != nil {
		return err
	}

	return s.videoRepo.Edit(commentId, message)
}

func (s *CommentService) Delete(commentId int, userId string) error {
	err := s.isUserOwner(commentId, userId)
	if err != nil {
		return err
	}
	return s.videoRepo.Delete(commentId)
}

func (s *CommentService) isUserOwner(commentId int, userId string) error {
	commentDb, err := s.videoRepo.FindById(commentId)
	if err != nil {
		return err
	}
	if commentDb.UserId != userId {
		return errors.New("your not owner")
	}
	return nil
}
