package repository

import "github.com/bearname/videohost/videoserver/model"

type VideoRepository interface {
	GetVideo(id string) (*model.Video, error)
	GetVideoList(page int, count int) ([]model.VideoListItem, error)
	NewVideo(userId string, videoId string, title string, description string, url string) error
	GetPageCount(countVideoOnPage int) (int, bool)
	AddVideoQuality(id string, quality string) bool
	SearchVideo(searchString string, page int, count int) ([]model.VideoListItem, error)
	IncrementViews(id string) bool
	FindUserVideos(userId string, page int, count int) ([]model.VideoListItem, error)
}
