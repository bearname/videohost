package controller

import (
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/controller"
	"github.com/bearname/videohost/internal/stream-service/app"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type StreamController struct {
	controller.BaseController
	streamService app.StreamService
}

func NewStreamController(streamService app.StreamService) *StreamController {
	v := new(StreamController)
	v.streamService = streamService
	return v
}

func (c *StreamController) StreamHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	vars := mux.Vars(request)
	var ok bool
	variable, err := c.BaseController.ParseMuxVariable(request, []string{"videoId"})
	if err != nil {
		c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "401 id not present")
		return
	}

	videoId := variable[0]

	var qualityString string
	qualityString, hasQuality := vars["quality"]
	segName, ok := vars["segName"]
	log.Info("videoId: " + videoId + " segName " + segName)
	if !ok && hasQuality {
		log.Info("serveHls" + qualityString + "M3u8")
		m3u8Name := "index-" + qualityString + ".m3u8"

		c.streamService.ServeHlsM3u8(writer, request, videoId, m3u8Name)
	} else if !ok {
		log.Info("serveHlsM3u8")
		c.streamService.ServeHlsM3u8(writer, request, videoId, "index.m3u8")
	} else {
		log.Info("serveHlsTs")
		c.streamService.ServeHlsTs(writer, request, segName, videoId)
	}
}
