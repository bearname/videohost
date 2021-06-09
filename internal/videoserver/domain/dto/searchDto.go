package dto

import "github.com/bearname/videohost/internal/videoserver/domain/model"

type SearchResultDto struct {
	Videos    []model.VideoListItem
	PageCount int
}

type SearchDto struct {
	Page         int
	Count        int
	SearchString string
}
