package domain

import (
	"github.com/bearname/videohost/pkg/videoserver/app/dto"
	"github.com/bearname/videohost/pkg/videoserver/domain/model"
)

type VideoRepository interface {
	Create(userId string, videoId string, title string, description string, url string) error
	Save(video *model.Video) error
	Find(videoId string) (*model.Video, error)
	FindVideosByPage(page int, count int) ([]model.VideoListItem, error)
	FindUserVideos(userId string, dto dto.SearchDto) ([]model.VideoListItem, error)
	Update(videoId string, title string, description string) error
	Delete(videoId string) error
	GetPageCount(countVideoOnPage int) (int, bool)
	AddVideoQuality(videoId string, quality string) error
	IncrementViews(videoId string) bool
	SearchVideo(searchString string, page int, count int) ([]model.VideoListItem, error)
}
