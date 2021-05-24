package service

import (
	"github.com/bearname/videohost/pkg/video-scaler/domain"
	"github.com/bearname/videohost/pkg/videoserver/domain/repository"
	log "github.com/sirupsen/logrus"
)

type VideoScaleHandler struct {
	videoScaleService *ScalerService
	videoRepo         repository.VideoRepository
	qualities         []domain.Quality
}

func NewVideoScaleHandler(service *ScalerService, videoRepository repository.VideoRepository, qualities []domain.Quality) *VideoScaleHandler {
	v := new(VideoScaleHandler)
	v.videoScaleService = service
	v.videoRepo = videoRepository
	v.qualities = qualities
	return v
}

func (h *VideoScaleHandler) Handle(message string) error {
	videoId := message
	log.Printf("'%s'", videoId)
	video, err := h.videoRepo.Find(videoId)
	if err != nil {
		return err
	}
	log.Info("Uploaded video " + video.Id + " " + video.Uploaded)
	inputVideoPath := "C:\\Users\\mikha\\go\\src\\videohost\\bin\\videoserver\\" + video.Url
	ok := h.videoScaleService.PrepareToStream(video.Id, inputVideoPath, h.qualities, video.OwnerId)
	log.Info(h.getResultMessage(ok))
	return nil
}

func (h *VideoScaleHandler) getResultMessage(ok bool) string {
	message := "Add video ok "
	if ok {
		message += "success"
	} else {
		message += "failed"
	}
	return message
}
