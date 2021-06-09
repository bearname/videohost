package service

import (
	"fmt"
	"net/http"
)

type StreamServiceImpl struct {
}

func NewStreamServiceImpl() StreamServiceImpl {
	return *new(StreamServiceImpl)
}

func (s *StreamServiceImpl) ServeHlsM3u8(w http.ResponseWriter, r *http.Request, videoId string, m3u8Name string) {
	mediaBase := s.getMediaBase(videoId)

	mediaFile := fmt.Sprintf("%s\\%s", mediaBase, m3u8Name)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "application/x-mpegURL")
}

func (s *StreamServiceImpl) ServeHlsTs(w http.ResponseWriter, r *http.Request, segName, videoId string) {
	mediaBase := s.getMediaBase(videoId)

	mediaFile := fmt.Sprintf("%s\\%s", mediaBase, segName)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "video/MP2T")
}

func (_ *StreamServiceImpl) getMediaBase(id string) string {
	return "content\\" + id
}
