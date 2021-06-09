package service

import (
	"encoding/json"
	"errors"
	"github.com/bearname/videohost/internal/common/util"
	"github.com/bearname/videohost/internal/video-scaler/domain"
	"github.com/bearname/videohost/internal/videoserver/domain/model"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type VideoScaleVisitor struct {
	videoScaleService  domain.VideoScaleService
	qualities          []domain.Quality
	videoServerAddress string
	authServerAddress  string
}

func NewVideoScaleHandler(service domain.VideoScaleService, qualities []domain.Quality, videoServerAddress string, authServerAddress string) *VideoScaleVisitor {
	v := new(VideoScaleVisitor)
	v.videoScaleService = service
	v.qualities = qualities
	v.videoServerAddress = videoServerAddress
	v.authServerAddress = authServerAddress
	return v
}

func (h *VideoScaleVisitor) Handle(message string) error {
	videoId := message
	log.Printf("'%s'", videoId)
	client := http.Client{}
	token, ok := util.InitAccessToken(&client, h.authServerAddress)
	if !ok {
		return errors.New("invalid failed get access token")
	}
	response, err := util.GetRequest(&http.Client{}, h.videoServerAddress+"/api/v1/videos/"+videoId, token.RefreshToken)
	if err != nil {
		return err
	}
	var video model.Video
	err = json.Unmarshal(response, &video)
	if err != nil {
		return err
	}

	log.Info("Uploaded video " + video.Id + " " + video.Uploaded)
	inputVideoPath := "C:\\Users\\mikha\\go\\src\\videohost\\bin\\videoserver\\" + video.Url
	ok = h.videoScaleService.PrepareToStream(video.Id, inputVideoPath, h.qualities, video.OwnerId)
	log.Info(h.getResultMessage(ok))
	return nil
}

func (h *VideoScaleVisitor) getResultMessage(ok bool) string {
	message := "Add video ok "
	if ok {
		message += "success"
	} else {
		message += "failed"
	}
	return message
}
