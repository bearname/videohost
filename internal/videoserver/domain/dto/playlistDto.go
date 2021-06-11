package dto

import "github.com/bearname/videohost/internal/videoserver/domain/model"

type CreatePlaylistDto struct {
	Name    string
	OwnerId string
	Privacy model.PrivacyType
	Videos  []string
}
