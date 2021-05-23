package main

import (
	"fmt"
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/video-scaler/app/service"
	"github.com/bearname/videohost/pkg/video-scaler/domain"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/mysql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func main() {
	//log.SetFormatter(&log.JSONFormatter{})
	//file, err := os.OpenFile("scaler.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	//if err == nil {
	//	log.SetOutput(file)
	//	defer file.Close()
	//}
	//log.Info("Started")

	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()
	//
	//channel, err := conn.Channel()
	//failOnError(err, "Failed to open a channel")
	//defer channel.Close()
	//
	//topic := "events_topic"
	//err = channel.ExchangeDeclare(
	//	topic,
	//	"topic",
	//	true,
	//	false,
	//	false,
	//	false,
	//	nil,
	//)
	//failOnError(err, "Failed to declare an exchange")
	//
	//queue, err := channel.QueueDeclare(
	//	"",
	//	true,
	//	false,
	//	true,
	//	false,
	//	nil,
	//)
	//failOnError(err, "Failed to declare a queue")
	//
	//routingKey := "events.upload-video"
	//log.Printf("Binding queue %s to exchange %s with routing key %s", queue.Name, topic, routingKey)
	//
	//err = channel.QueueBind(
	//	queue.Name,
	//	routingKey,
	//	topic,
	//	false,
	//	nil)
	//failOnError(err, "Failed to bind a queue")
	//
	//messages, err := channel.Consume(
	//	queue.Name,
	//	"",
	//	true,
	//	false,
	//	false,
	//	false,
	//	nil,
	//)
	//failOnError(err, "Failed to register a consumer")
	//
	//forever := make(chan bool)
	//
	var connector database.Connector
	err := connector.Connect()
	if err != nil {
		panic("unable to connect to connector" + err.Error())
	}

	defer connector.Close()

	videoRepo := mysql.NewMysqlVideoRepository(connector)
	scalerService := service.NewScalerService(videoRepo)
	qualities := []domain.Quality{domain.Q1080p, domain.Q720p, domain.Q480p, domain.Q320p}
	videos, err := videoRepo.FindVideosByPage(0, 100)
	for _, video := range videos {
		fmt.Println(video)
		inputVideoPath := "C:\\Users\\mikha\\go\\src\\videohost\\bin\\videoserver\\content\\" + video.Id + "\\index.mp4"
		if video.Quality == "" {
			ok := scalerService.ScaleVideos(video.Id, inputVideoPath, qualities)
			log.Info(getResultMessage(ok))
		}
	}

	//go func() {
	//	for data := range messages {
	//		videoId := data.Body
	//		log.Printf("'%s'", videoId)
	//		id := string(videoId)
	//		video, err2 := videoRepo.Find(id)
	//		if err2 != nil {
	//			log.Error(err2.Error())
	//		} else {
	//			log.Info("Uploaded video " + video.Id + " " + video.Uploaded)
	//			//TODO save to file server
	//			inputVideoPath := "C:\\Users\\mikha\\go\\src\\videohost\\bin\\videoserver\\" + video.Url
	//			ok := scalerService.ScaleVideos(video.Id, inputVideoPath, qualities)
	//			log.Info(getResultMessage(ok))
	//		}
	//	}
	//}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	//<-forever
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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
