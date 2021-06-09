package model

import "time"

type Action int

const (
	AddLike Action = iota
	AddDislike
	DeleteLike
	DeleteDisLike
)

func ActionToString(action Action) string {
	switch action {
	case AddLike:
		return "add like"
	case AddDislike:
		return "add dislike"
	case DeleteLike:
		return "delete like"
	case DeleteDisLike:
		return "delete dislike"
	}

	return "unknown action"
}

type Like struct {
	IdVideo string    `json:"id_video"`
	OwnerId string    `json:"owner_id"`
	Liked   time.Time `json:"liked"`
	IsLike  bool      `json:"is_like"`
}

type VideoLikes struct {
	VideoId      string `json:"video_id"`
	CountLike    int    `json:"like"`
	CountDislike int    `json:"dislike"`
}
