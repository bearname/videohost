package repository

import "github.com/bearname/videohost/videoserver/model"

type VideoRepository interface {
	GetVideo(id string) (*model.Video, error)
	GetVideoList(page int, count int) ([]model.VideoListItem, error)
	NewVideo(id string, fileName string, description string, url string) error
	GetPageCount(countVideoOnPage int) (int, bool)
}
