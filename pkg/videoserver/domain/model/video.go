package model

type Video struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Thumbnail   string `json:"thumbnail"`
	Url         string `json:"url"`
	Uploaded    string `json:"uploaded"`
	Quality     string `json:"quality"`
	Views       string `json:"views"`
}
