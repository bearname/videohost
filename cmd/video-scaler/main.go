package main

import (
	"fmt"
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/video-scaler/domain"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/mysql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	topic := "events_topic"
	err = channel.ExchangeDeclare(
		topic,   // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	queue, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	routingKey := "events.upload-video"
	log.Printf("Binding queue %s to exchange %s with routing key %s", queue.Name, topic, routingKey)

	err = channel.QueueBind(
		queue.Name, // queue name
		routingKey, // routing key
		topic,      // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	messages, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	var connector database.Connector
	err = connector.Connect()
	if err != nil {
		panic("unable to connect to connector" + err.Error())
	}
	defer connector.Close()

	videoRepository := mysql.NewMysqlVideoRepository(connector)
	go func() {
		for data := range messages {
			videoId := data.Body
			log.Printf("'%s'", videoId)
			video, err2 := videoRepository.Find(string(videoId))
			if err2 != nil {
				log.Error(err2)
			}
			log.Info("Uploaded video " + video.Id + " " + video.Uploaded)
			//TODO save to file server
			inputVideoPath := "..\\cmd\\videoserver\\" + video.Url
			videos := scaleVideos(inputVideoPath)
			if videos {
				quality := videoRepository.AddVideoQuality(video.Id, "720, 480, 320, 144")
				message := getResultMessage(quality)
				log.Info(message)
			}
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

func getResultMessage(quality bool) string {
	message := "Add video quality "
	if quality {
		message += "success"
	} else {
		message += "failed"
	}
	return message
}

func IsSupportedQuality(quality int) bool {
	ints := []int{2160, 1440, 1080, 720, 480, 320, 144}
	for _, a := range ints {
		if a == quality {
			return true
		}
	}

	return false
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func scaleVideos(inputVideoPath string) bool {
	const extension = ".mp4"

	output, err := exec.Command("./resolution.bat", inputVideoPath).Output()
	if err != nil {
		log.Error()
		return false
	}

	split := strings.Split(string(output), "\n")
	height, ok := getDimension(split, "height")
	if !ok {
		log.Error("Failed get resolution")
		return false
	}

	if !IsSupportedQuality(height) {
		log.Error("Not supported quality")
		return false
	}

	scaleVideosToQuality(inputVideoPath, extension)

	return true
}

func scaleVideosToQuality(inputVideoPath string, extension string) {
	qualities := []domain.Quality{domain.Q720p, domain.Q480p, domain.Q320p, domain.Q144p}

	for _, quality := range qualities {
		scaleVideoToQuality(inputVideoPath, extension, quality)
	}
}

func scaleVideoToQuality(inputVideoPath string, extension string, quality domain.Quality) {
	outputVideoPath := getOutputVideoPath(inputVideoPath, extension, quality)
	log.Info("scale video to " + quality.String())
	video := scaleVideo(inputVideoPath, outputVideoPath, quality)
	if !video {
		log.Error("Failed convert")
	}
	outputHls := outputVideoPath[0 : strings.LastIndex(outputVideoPath, "\\")+1]
	inputVideoPath = strings.ReplaceAll(inputVideoPath, "\\", "\\")
	outputHls = strings.ReplaceAll(outputHls, "\\", "\\")

	_, err := prepareToStream(outputVideoPath, outputHls)
	if err != nil {
		log.Error("Failed prepare to stream file " + inputVideoPath + " in quality " + quality.String() + "p")
	}
}

func scaleVideo(inputVideoPath string, outputVideoPath string, quality domain.Quality) bool {
	err := exec.Command("scale.bat", inputVideoPath, quality.String(), outputVideoPath).Run()
	if err != nil {
		log.Error(err.Error())
		return false
	}
	return true
}

func getDimension(split []string, key string) (int, bool) {
	value := strings.Split(split[1], "=")
	if value[0] != key {
		return 0, false
	}
	fmt.Println("'" + value[1] + "'")

	s := value[1][0 : len(value[1])-1]
	atoi, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}

	return atoi, true
}

func getOutputVideoPath(videoPath string, extension string, quality domain.Quality) string {
	return videoPath[0:len(videoPath)-len(extension)] + "-" + quality.String() + "p" + extension
}

func prepareToStream(videoPath string, output string) ([]byte, error) {
	return exec.Command(`ffmpeg`, `-i`, videoPath, `-profile:v`, `baseline`, `-level`, `3.0`, `-s`, `640x360`,
		`-start_number`, `0`, `-hls_time`, `10`, `-hls_list_size`, `0`, `-f`, `hls`, output+`index.m3u8`).Output()
}
