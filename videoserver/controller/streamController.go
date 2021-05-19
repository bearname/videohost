package controller

import (
	"fmt"
	"github.com/bearname/videohost/videoserver/repository"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type StreamController struct {
	BaseController
	videoRepository repository.VideoRepository
}

func NewStreamController(videoRepository repository.VideoRepository) *StreamController {
	v := new(StreamController)
	v.videoRepository = videoRepository
	return v
}

func (c *StreamController) StreamHandler(writer http.ResponseWriter, request *http.Request) {
	writer = *c.BaseController.AllowCorsRequest(&writer)
	vars := mux.Vars(request)

	var id string
	var ok bool

	if id, ok = vars["id"]; !ok {
		http.Error(writer, "404 video not found not found", http.StatusNotFound)
		return
	}

	video, err := c.videoRepository.GetVideo(id)
	if err != nil {
		log.Error(err.Error())
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	mediaBase := c.getMediaBase(video.Id)
	segName, ok := vars["segName"]
	log.Info("id: " + id + " segName " + segName)
	if !ok {
		m3u8Name := "index.m3u8"
		log.Info("serveHlsM3u8")
		c.serveHlsM3u8(writer, request, mediaBase, m3u8Name)
	} else {
		log.Info("serveHlsTs")
		c.serveHlsTs(writer, request, mediaBase, segName)
	}
}

func (_ *StreamController) getMediaBase(id string) string {
	return "content\\" + id
}

func (_ *StreamController) serveHlsM3u8(w http.ResponseWriter, r *http.Request, mediaBase, m3u8Name string) {
	mediaFile := fmt.Sprintf("%s\\%s", mediaBase, m3u8Name)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "application/x-mpegURL")
}

func (_ *StreamController) serveHlsTs(w http.ResponseWriter, r *http.Request, mediaBase, segName string) {
	mediaFile := fmt.Sprintf("%s\\%s", mediaBase, segName)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "video/MP2T")
}
