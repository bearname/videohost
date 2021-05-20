package repository

import (
	model2 "github.com/bearname/videohost/pkg/videoserver/domain/model"
)

type VideoRepository interface {
	GetVideo(id string) (*model2.Video, error)
	GetVideoList(page int, count int) ([]model2.VideoListItem, error)
	NewVideo(userId string, videoId string, title string, description string, url string) error
	GetPageCount(countVideoOnPage int) (int, bool)
	AddVideoQuality(id string, quality string) bool
	SearchVideo(searchString string, page int, count int) ([]model2.VideoListItem, error)
	IncrementViews(id string) bool
	FindUserVideos(userId string, page int, count int) ([]model2.VideoListItem, error)
}
