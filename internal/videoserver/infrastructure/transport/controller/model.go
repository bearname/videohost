package controller

import (
	"github.com/bearname/videohost/internal/videoserver/domain"
	"github.com/bearname/videohost/internal/videoserver/domain/model"
)

type CreatePlayListRequest struct {
	Name     string            `json:"name"`
	Privacy  model.PrivacyType `json:"privacy"`
	VideosId []string          `json:"videos"`
}

type PlayListVideoModificationRequest struct {
	Action domain.Action `json:"act"`
	Videos []string      `json:"videos"`
}
