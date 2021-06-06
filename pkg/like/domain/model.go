package domain

type Like struct {
	IdVideo string `json:"id_video"`
	OwnerId string `json:"owner_id"`
	Liked   string `json:"liked"`
	IsLike  bool   `json:"is_like"`
}

type VideoLikes struct {
	VideoId      string `json:"video_id"`
	CountLike    int    `json:"like"`
	CountDislike int    `json:"dislike"`
}