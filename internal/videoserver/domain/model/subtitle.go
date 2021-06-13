package model

import "github.com/bearname/videohost/internal/common/domain"

type Subtitle struct {
	domain.IdInt
	VideoId string         `json:"video_id"`
	Parts   []SubtitlePart `json:"parts"`
}

type SubtitlePart struct {
	domain.IdInt
	SubtitleId int    `json:"subtitle_id"`
	Start      int    `json:"start"`
	End        int    `json:"end"`
	Text       string `json:"text"`
}
