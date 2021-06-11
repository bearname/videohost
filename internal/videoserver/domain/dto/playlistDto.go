package dto

import (
	"github.com/bearname/videohost/internal/videoserver/domain/model"
	"time"
)

type CreatePlaylistDto struct {
	Name    string
	OwnerId string
	Privacy model.PrivacyType
	Videos  []string
}

type PlaylistItem struct {
	Id      int               `json:"id"`
	Name    string            `json:"name"`
	OwnerId string            `json:"ownerId"`
	Created time.Time         `json:"create"`
	Privacy model.PrivacyType `json:"privacy"`
}
