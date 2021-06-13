package repository

import (
	"github.com/bearname/videohost/internal/videoserver/domain/dto"
	"github.com/bearname/videohost/internal/videoserver/domain/model"
)

type SubtitleRepository interface {
	Create(subtitle dto.CreateSubtitleRequestDto) (int64, error)
	Find(videoId int) (model.Subtitle, error)
	Update(subtitle model.Subtitle) error
	Delete(subtitleId int) error
}
