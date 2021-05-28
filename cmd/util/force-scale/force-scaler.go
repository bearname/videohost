package main

import (
	"fmt"
	mysqlConnector "github.com/bearname/videohost/pkg/common/infrarstructure/mysql"
	"github.com/bearname/videohost/pkg/video-scaler/app/service"
	"github.com/bearname/videohost/pkg/video-scaler/domain"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/mysql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func main() {
	var connector mysqlConnector.ConnectorImpl
	err := connector.Connect()
	if err != nil {
		fmt.Println("unable to connect to mysqlConnector" + err.Error())
	}

	defer connector.Close()

	videoRepo := mysql.NewMysqlVideoRepository(&connector)
	scalerService := service.NewVideoScaleService(nil, videoRepo, "", "")
	qualities := []domain.Quality{domain.Q1080p, domain.Q720p, domain.Q480p, domain.Q320p}
	videos, err := videoRepo.FindVideosByPage(0, 100)
	for _, video := range videos {
		fmt.Println(video)
		inputVideoPath := "C:\\Users\\mikha\\go\\src\\videohost\\bin\\videoserver\\content\\" + video.Id + "\\index.mp4"
		if video.Quality == "" {
			ok := scalerService.PrepareToStream(video.Id, inputVideoPath, qualities, "")
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
