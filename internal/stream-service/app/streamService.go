package app

import (
	"fmt"
	"net/http"
)

type StreamService struct {
}

func NewStreamService() StreamService {
	return *new(StreamService)
}

func (s *StreamService) ServeHlsM3u8(w http.ResponseWriter, r *http.Request, videoId string, m3u8Name string) {
	mediaBase := s.getMediaBase(videoId)

	mediaFile := fmt.Sprintf("%s\\%s", mediaBase, m3u8Name)
	w.Header().Set("Content-Type", "application/x-mpegURL")
	//w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, mediaFile)
}

func (s *StreamService) ServeHlsTs(w http.ResponseWriter, r *http.Request, segName, videoId string) {
	mediaBase := s.getMediaBase(videoId)

	mediaFile := fmt.Sprintf("%s\\%s", mediaBase, segName)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "video/MP2T")
}

func (_ *StreamService) getMediaBase(id string) string {
	return "content\\" + id
}
