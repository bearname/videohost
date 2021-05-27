package service

import (
	"github.com/bearname/videohost/pkg/video-scaler/domain"
	"github.com/bearname/videohost/pkg/videoserver/domain/repository"
	log "github.com/sirupsen/logrus"
)

type VideoScaleVisitor struct {
	videoScaleService domain.ScalerService
	videoRepo         repository.VideoRepository
	qualities         []domain.Quality
}

func NewVideoScaleHandler(service domain.ScalerService, videoRepository repository.VideoRepository, qualities []domain.Quality) *VideoScaleVisitor {
	v := new(VideoScaleVisitor)
	v.videoScaleService = service
	v.videoRepo = videoRepository
	v.qualities = qualities
	return v
}

func (h *VideoScaleVisitor) Handle(message string) error {
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

func (h *VideoScaleVisitor) getResultMessage(ok bool) string {
	message := "Add video ok "
	if ok {
		message += "success"
	} else {
		message += "failed"
	}
	return message
}
