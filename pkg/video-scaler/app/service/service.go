package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bearname/videohost/pkg/common/amqp"
	"github.com/bearname/videohost/pkg/common/util"
	"github.com/bearname/videohost/pkg/video-scaler/domain"
	"github.com/bearname/videohost/pkg/videoserver/domain/repository"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type ScalerService struct {
	videoRepo     repository.VideoRepository
	messageBroker *amqp.RabbitMqService
	token         *util.Token
}

func NewScalerService(service *amqp.RabbitMqService, videoRepository repository.VideoRepository) *ScalerService {
	s := new(ScalerService)
	s.videoRepo = videoRepository
	s.messageBroker = service
	s.token = util.NewToken("", "")

	return s
}

func (s *ScalerService) PrepareToStream(videoId string, inputVideoPath string, allNeededQualities []domain.Quality, ownerId string) bool {
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
	video, err := s.videoRepo.Find(videoId)
	if err != nil {
		return false
	}

	filename := inputVideoPath[0:len(inputVideoPath)-len("index.mp4")] + "index.m3u8"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Error(err)
		return false
	}
	_, err = file.WriteString("#EXTM3U\n")
	if err != nil {
		log.Error(err)
		return false
	}

	availableVideoQualities := strings.Split(video.Quality, ",")
	for _, quality := range allNeededQualities {
		contains := util.Contains(availableVideoQualities, quality.String())
		if !contains {
			s.prepareToStreamByQuality(videoId, inputVideoPath, extension, quality, ownerId)
		}
	}

	return true
}

func (s *ScalerService) prepareToStreamByQuality(videoId string, inputVideoPath string, extension string, quality domain.Quality, ownerId string) {
	err := s.scaleVideoToQuality(inputVideoPath, extension, quality)
	if err != nil {
		log.Error("Failed prepare to stream file " + inputVideoPath + " in quality " + quality.String() + "p")
	} else {
		body := videoId + "," + quality.String() + "," + ownerId
		fmt.Println(body)
		log.Info("Success prepare to stream file " + inputVideoPath + " in quality " + quality.String() + "p")
		ok := s.addVideoQuality(videoId, quality)
		log.Info(s.getResultMessage(ok))
		err = s.messageBroker.Publish("events_topic", "events.video-scaled", body)
		if err != nil {
			log.Error("Failed publish event 'video-scaled")
		}
	}
}

func (s *ScalerService) addVideoQuality(videoId string, quality domain.Quality) bool {
	//ok := s.videoRepo.AddVideoQuality(videoId, quality.String())
	buf := struct {
		Quality int `json:"quality"`
	}{Quality: quality.Values()}

	marshal, err := json.Marshal(buf)
	if err != nil {
		return false
	}

	request, err := http.NewRequest("PUT", "http://localhost:8000/api/v1/videos/"+videoId+"/add-quality", bytes.NewBuffer(marshal))
	client := &http.Client{}
	if s.token.AccessToken == "" {
		token, err := util.GetAdminAccessToken(client, "http://localhost:8001")
		if err != nil {
			log.Error(err)
			return false
		}
		s.token = token
	}

	request.Header.Add("Authorization", "Bearer "+s.token.AccessToken)
	response, err := client.Do(request)
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		token, err := util.GetAdminAccessToken(client, "http://localhost:8001")
		if err != nil {
			log.Error(err)
			return false
		}
		s.token = token
	}

	if err != nil || response.StatusCode != http.StatusOK {
		log.Error("Failed get id of owner of the video ")
		return false
	}
	log.Info(s.getResultMessage(true))
	return true
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
	root := outputVideoPath[0 : strings.LastIndex(outputVideoPath, "\\")+1]
	outputHls := root + "index-" + quality.String() + `.m3u8`
	inputVideoPath = strings.ReplaceAll(inputVideoPath, "\\", "\\")
	outputHls = strings.ReplaceAll(outputHls, "\\", "\\")

	err := s.prepareToStream(inputVideoPath, outputHls, quality)
	if err != nil {
		return err
	}

	resolution := domain.QualityToResolution(quality)
	data := "#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=6221600,CODECS=\"mp4a.40.2,avc1.640028\",RESOLUTION=" + resolution.String() + ",NAME=\"" + quality.String() + "\"\n" +
		"/media/a7e608d9-bc76-11eb-afc7-e4e74940035b/" + quality.String() + "/stream/\n"

	file, err := os.OpenFile(root+"index.m3u8", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	_, err = file.WriteString(data)
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
