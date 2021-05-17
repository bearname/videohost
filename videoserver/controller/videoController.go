package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bearname/videohost/videoserver/repository"
	"github.com/bearname/videohost/videoserver/service"
	"github.com/bearname/videohost/videoserver/util"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os/exec"
	"path/filepath"
)

type VideoController struct {
	BaseController
	videoRepository repository.VideoRepository
	messageBroker   service.MessageBroker
}

func NewVideoController(videoRepository repository.VideoRepository) *VideoController {
	v := new(VideoController)

	v.videoRepository = videoRepository
	v.messageBroker = service.NewRabbitMqService("guest", "guest", "localhost", 5672)
	if v.messageBroker == nil {
		return nil
	}
	return v
}

func (c VideoController) GetVideo() func(w http.ResponseWriter, r *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		c.BaseController.SetupCors(&writer, request)
		vars := mux.Vars(request)
		var id string
		var ok bool

		if id, ok = vars["ID"]; !ok {
			http.Error(writer, "404 video not found not found", http.StatusNotFound)
			return
		}

		video, err := c.videoRepository.GetVideo(id)
		if err != nil {
			log.Error(err.Error())
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		c.writeResponseData(writer, video)
	}
}

func (c VideoController) GetVideoList() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		c.BaseController.SetupCors(&writer, request)

		var page int
		page, success := c.getIntRouteParameter(writer, request, "page")
		if !success {
			return
		}
		var countVideoOnPage int
		countVideoOnPage, success = c.getIntRouteParameter(writer, request, "countVideoOnPage")
		if !success {
			return
		}

		log.Info("page ", page, " count video ", countVideoOnPage)
		pageCount, ok := c.videoRepository.GetPageCount(countVideoOnPage)
		if !ok {
			http.Error(writer, "Failed get page countVideoOnPage", http.StatusInternalServerError)
			return
		}

		videos, err := c.videoRepository.GetVideoList(page, countVideoOnPage)
		if err != nil {
			log.Error(err.Error())
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		responseData := make(map[string]interface{})
		responseData["pageCount"] = pageCount
		responseData["videos"] = videos

		c.writeResponseData(writer, responseData)
	}
}

func (c VideoController) UploadVideo() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		c.BaseController.SetupCors(&writer, request)
		//request.ParseForm()
		//
		//form := request.Form
		//encode := form.Encode()
		//log.Info(encode)
		//videoName := form.Get("name")
		//if len(videoName) == 0 {
		//	http.Error(writer, "`name` not present", http.StatusBadRequest)
		//	return
		//}
		//
		//description := request.Form.Get("description")
		//if len(description) == 0 {
		//	http.Error(writer, "`description` not present", http.StatusBadRequest)
		//	return
		//}

		fileReader, header, err := request.FormFile("file")
		if err != nil {
			log.Error(err.Error())
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		contentType := header.Header.Get("Content-Type")
		if contentType != util.VideoContentType {
			log.Error("Unexpected content type", contentType)
			http.Error(writer, "Unexpected content type", http.StatusBadRequest)
			return
		}

		id, err := uuid.NewUUID()
		if err != nil {
			log.Error(err.Error())
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		videoId := id.String()

		err = util.CopyFile(fileReader, videoId)
		if err != nil {
			log.Error(err.Error())
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		err = c.videoRepository.NewVideo(
			videoId,
			"videoName",
			"description1",
			filepath.Join(util.ContentDir, videoId, util.VideoFileName),
		)
		if err != nil {
			log.Error(err.Error())
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		output, err := cmdExec(`ffprobe`, `-v`, `error`, `-select_streams`, `v:0`, `-show_entries`, `stream=width,height`, `-of`, `csv=s=x:p=0`, util.VideoFileName)
		log.Info("resolution of file " + util.VideoFileName + " equal " + output)
		c.messageBroker.Publish("events_topic", "events.upload-video", id.String())

		writer.WriteHeader(http.StatusOK)
		c.BaseController.JsonResponse(writer, id)
	}
}

func cmdExec(args ...string) (string, error) {
	baseCmd := args[0]
	cmdArgs := args[1:]

	cmd := exec.Command(baseCmd, cmdArgs...)
	var outputBuffer, errorBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &errorBuffer
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("out:", outputBuffer.String(), "err:", errorBuffer.String())

	return outputBuffer.String(), nil
}

func (c VideoController) getIntRouteParameter(writer http.ResponseWriter, request *http.Request, key string) (int, bool) {
	pageStr, done := c.parseRouteParameter(request, key)
	if !done {
		http.Error(writer, "400 "+key+" parameter not found", http.StatusBadRequest)
		return 0, false
	}

	page, b := c.validate(writer, pageStr)
	if b {
		http.Error(writer, "400 invalid page parameter", http.StatusBadRequest)
		return 0, false
	}
	return page, true
}

func (c VideoController) parseRouteParameter(request *http.Request, key string) (string, bool) {
	query := request.URL.Query()
	keys, ok := query[key]

	if !ok || len(keys) != 1 {
		return "", false
	}

	return keys[0], true
}

func (c VideoController) validate(writer http.ResponseWriter, pageStr string) (int, bool) {
	page, b := util.StrToInt(pageStr)
	if !b || page < 0 {
		http.Error(writer, "400 Invalid page parameter", http.StatusBadRequest)
		return 0, true
	}
	return page, false
}

func (c VideoController) writeResponseData(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if _, err = io.WriteString(w, string(b)); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}
