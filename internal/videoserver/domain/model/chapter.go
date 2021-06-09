package model

type Chapter struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Start   int    `json:"start"`
	End     int    `json:"end"`
	VideoId string `json:"videoId"`
}
