package repository

import (
	"github.com/bearname/videohost/pkg/videoserver/domain/model"
)

type VideoRepository interface {
	Create(userId string, videoId string, title string, description string, url string) error
	Save(video *model.Video) error
	Find(videoId string) (*model.Video, error)
	FindVideosByPage(page int, count int) ([]model.VideoListItem, error)
	FindUserVideos(userId string, page int, count int) ([]model.VideoListItem, error)
	Update(videoId string, title string, description string) error
	Delete(videoId string) error
	GetPageCount(countVideoOnPage int) (int, bool)
	AddVideoQuality(id string, quality string) bool
	IncrementViews(id string) bool
	SearchVideo(searchString string, page int, count int) ([]model.VideoListItem, error)
}
