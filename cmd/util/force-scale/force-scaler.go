package main

import (
	"fmt"
	mysqlConnector "github.com/bearname/videohost/internal/common/infrarstructure/mysql"
	"github.com/bearname/videohost/internal/video-scaler/app/service"
	"github.com/bearname/videohost/internal/video-scaler/domain"
	"github.com/bearname/videohost/internal/videoserver/infrastructure/mysql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func main() {
	var connector mysqlConnector.ConnectorImpl
	err := connector.Connect("root", "123", "localhost:3306", "video")
	if err != nil {
		fmt.Println("unable to connect to mysqlConnector" + err.Error())
		return
	}

	defer connector.Close()

	videoRepo := mysql.NewMysqlVideoRepository(&connector)
	scalerService := service.NewVideoScaleService(nil, "", "")
	qualities := []domain.Quality{domain.Q1080p, domain.Q720p, domain.Q480p, domain.Q320p}
	videos, err := videoRepo.FindVideosByPage(0, 100)
	if err != nil {
		fmt.Println("unable find FindVideosByPage" + err.Error())
		return
	}

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
