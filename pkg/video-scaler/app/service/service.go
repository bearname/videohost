package service

import (
	"fmt"
	"github.com/bearname/videohost/pkg/common/amqp"
	"github.com/bearname/videohost/pkg/video-scaler/domain"
	"github.com/bearname/videohost/pkg/videoserver/domain/repository"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strconv"
	"strings"
)

type ScalerService struct {
	videoRepo     repository.VideoRepository
	messageBroker *amqp.RabbitMqService
}

func NewScalerService(service *amqp.RabbitMqService, videoRepository repository.VideoRepository) *ScalerService {
	s := new(ScalerService)
	s.videoRepo = videoRepository
	s.messageBroker = service
	return s
}

func (s *ScalerService) PrepareToStream(videoId string, inputVideoPath string, qualities []domain.Quality) bool {
	const extension = ".mp4"

	log.Info(inputVideoPath)
	output, err := exec.Command("C:\\Users\\mikha\\go\\src\\videohost\\bin\\video-scaler\\resolution.bat", inputVideoPath).Output()
	if err != nil {
		log.Error(err.Error())
		return false
	}

	split := strings.Split(string(output), "\n")
	height, ok := s.getDimension(split, "height")
	if !ok {
		log.Error("Failed get resolution")
		return false
	}

	if !domain.IsSupportedQuality(height) {
		log.Error("Not supported quality")
		return false
	}

	s.prepareToStreamByQualities(qualities, videoId, inputVideoPath, extension)

	return true
}

func (s *ScalerService) prepareToStreamByQualities(qualities []domain.Quality, videoId string, inputVideoPath string, extension string) {
	for _, quality := range qualities {
		err := s.scaleVideoToQuality(inputVideoPath, extension, quality)
		if err != nil {
			log.Error("Failed prepare to stream file " + inputVideoPath + " in quality " + quality.String() + "p")
		} else {
			body := videoId + "," + quality.String() + "," + "mihail12russ@gmail.com"
			fmt.Println(body)
			s.messageBroker.Publish("events_topic", "events.video-scaled", body)
			log.Info("Success prepare to stream file " + inputVideoPath + " in quality " + quality.String() + "p")
			ok := s.videoRepo.AddVideoQuality(videoId, quality.String())
			log.Info(s.getResultMessage(ok))
		}
	}
}

func (s *ScalerService) getResultMessage(quality bool) string {
	message := "Add video quality "
	if quality {
		message += "success"
	} else {
		message += "failed"
	}
	return message
}
func (s *ScalerService) scaleVideoToQuality(inputVideoPath string, extension string, quality domain.Quality) error {
	outputVideoPath := s.getOutputVideoPath(inputVideoPath, extension, quality)
	log.Info("prepare video to stream on quality " + quality.String() + "p")
	outputHls := outputVideoPath[0:strings.LastIndex(outputVideoPath, "\\")+1] + "index-" + quality.String() + `.m3u8`
	inputVideoPath = strings.ReplaceAll(inputVideoPath, "\\", "\\")
	outputHls = strings.ReplaceAll(outputHls, "\\", "\\")

	err := s.prepareToStream(inputVideoPath, outputHls, quality)
	if err != nil {
		return err
	}

	return nil
}

func (s *ScalerService) getDimension(split []string, key string) (int, bool) {
	value := strings.Split(split[1], "=")
	if value[0] != key {
		return 0, false
	}
	fmt.Println("'" + value[1] + "'")

	number := value[1][0 : len(value[1])-1]
	atoi, err := strconv.Atoi(number)
	if err != nil {
		return 0, false
	}

	return atoi, true
}

func (s *ScalerService) prepareToStream(videoPath string, output string, quality domain.Quality) error {
	resolution := domain.QualityToResolution(quality)
	fmt.Println(resolution)
	return exec.Command(`ffmpeg`, `-i`, videoPath, `-profile:v`, `baseline`, `-level`, `3.0`, `-s`, resolution.String(),
		`-start_number`, `0`, `-hls_time`, `10`, `-hls_list_size`, `0`, `-f`, `hls`, output).Run()
}

func (s *ScalerService) getOutputVideoPath(videoPath string, extension string, quality domain.Quality) string {
	return videoPath[0:len(videoPath)-len(extension)] + "-" + quality.String() + "p" + extension
}
