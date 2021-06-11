package model

import (
	"github.com/bearname/videohost/internal/common/db"
	"time"
)

type PrivacyType int

const (
	Public PrivacyType = iota
	Unlisted
	Private
)

func (p PrivacyType) Int() int {
	switch p {
	case Public:
		return 0
	case Unlisted:
		return 1
	case Private:
		return 2
	default:
		return -1
	}
}

type Playlist struct {
	Id          string
	Name        string
	OwnerId     string
	Created     time.Time
	Privacy     PrivacyType
	VideoString string
	Videos      []VideoListItem
}

type PlaylistFilter struct {
	page    db.Page
	OrderBy db.OrderBy
}
