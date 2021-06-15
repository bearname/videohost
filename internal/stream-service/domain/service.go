package domain

import "net/http"

type StreamService interface {
	ServeHlsM3u8(w http.ResponseWriter, r *http.Request, videoId string, m3u8Name string)
	ServeHlsTs(w http.ResponseWriter, r *http.Request, segName, videoId string)
}
