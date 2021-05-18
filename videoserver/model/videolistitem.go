package model

type VideoListItem struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Duration  int    `json:"duration"`
	Thumbnail string `json:"thumbnail"`
	Uploaded  string `json:"uploaded"`
	Views     string `json:"views"`
}
