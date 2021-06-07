package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	logFile := "videoprocessor.log"
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	if len(os.Args) < 3 {
		fmt.Println("Usage videoprocessor.exe <video file path> <output thumbnail path> <outputPath>")
	}
	videoPath := os.Args[1]
	duration, err := getVideoDuration(videoPath)
	if err != nil {
		fmt.Println("Failed get Video Duration")
		log.Fatal("Failed get Video Duration" + err.Error())
	}

	fmt.Printf("%f", duration)

	thumbnailPath := os.Args[2]

	err = createVideoThumbnail(videoPath, thumbnailPath, int64(duration)/2)
	if err != nil {
		fmt.Println("Failed createVideoThumbnail")
		log.Fatal("Failed createVideoThumbnail" + err.Error())
	}
}

func getVideoDuration(videoPath string) (float64, error) {
	result, err := exec.Command(`ffprobe`, `-v`, `error`, `-show_entries`, `format=duration`, `-of`,
		`default=noprint_wrappers=1:nokey=1`, videoPath).Output()
	if err != nil {
		return 0.0, err
	}

	return strconv.ParseFloat(strings.Trim(string(result), "\n\r"), 64)
}

func ffmpegTimeFromSeconds(seconds int64) string {
	return time.Unix(seconds, 0).UTC().Format(`15:04:05.000000`)
}

func prepareToStream(videoPath string, output string) ([]byte, error) {
	return exec.Command(`ffmpeg`, `-i`, videoPath, `-profile:v`, `baseline`, `-level`, `3.0`, `-s`, `640x360`,
		`-start_number`, `0`, `-hls_time`, `10`, `-hls_list_size`, `0`, `-f`, `hls`, output+`index.m3u8`).Output()
}

func createVideoThumbnail(videoPath string, thumbnailPath string, thumbnailOffset int64) error {
	return exec.Command(`ffmpeg`, `-i`, videoPath, `-ss`, ffmpegTimeFromSeconds(thumbnailOffset), `-vframes`, `1`, thumbnailPath).Run()
}
