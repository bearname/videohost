package main

import (
	"fmt"
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/video-scaler/app/service"
	"github.com/bearname/videohost/pkg/video-scaler/domain"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/mysql"
	log "github.com/sirupsen/logrus"
)

func main() {
	var connector database.Connector
	err := connector.Connect()
	if err != nil {
		fmt.Println("unable to connect to connector" + err.Error())
	}

	defer connector.Close()

	videoRepo := mysql.NewMysqlVideoRepository(connector)
	scalerService := service.NewScalerService(nil, videoRepo)
	qualities := []domain.Quality{domain.Q1080p, domain.Q720p, domain.Q480p, domain.Q320p}
	videos, err := videoRepo.FindVideosByPage(0, 100)
	for _, video := range videos {
		fmt.Println(video)
		inputVideoPath := "C:\\Users\\mikha\\go\\src\\videohost\\bin\\videoserver\\content\\" + video.Id + "\\index.mp4"
		if video.Quality == "" {
			ok := scalerService.PrepareToStream(video.Id, inputVideoPath, qualities)
			log.Info(getResultMessage(ok))
		}
	}
}

func getResultMessage(ok bool) string {
	message := "Add video ok "
	if ok {
		message += "success"
	} else {
		message += "failed"
	}
	return message
}
