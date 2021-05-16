package repository

import "github.com/bearname/videohost/videoserver/model"

type VideoRepository interface {
	GetVideo(id string) (*model.Video, error)
	GetVideoList() ([]model.VideoListItem, error)
	UploadVideo(id string, fileName string, url string) error
}
