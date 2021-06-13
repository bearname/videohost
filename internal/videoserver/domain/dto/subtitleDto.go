package dto

type CreateSubtitleRequestDto struct {
	VideoId string            `json:"videoId"`
	Parts   []SubtitlePartDto `json:"items"`
}

type SubtitlePartDto struct {
	SubtitleId int    `json:"subtitleId"`
	Start      int    `json:"start"`
	End        int    `json:"end"`
	Text       string `json:"text"`
}
