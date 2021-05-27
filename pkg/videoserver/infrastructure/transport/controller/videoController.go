package controller

import (
	"encoding/json"
	"github.com/bearname/videohost/pkg/common/amqp"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/controller"
	"github.com/bearname/videohost/pkg/common/util"
	"github.com/bearname/videohost/pkg/videoserver/app/dto"
	"github.com/bearname/videohost/pkg/videoserver/app/service"
	"github.com/bearname/videohost/pkg/videoserver/domain/model"
	"github.com/bearname/videohost/pkg/videoserver/domain/repository"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	//"path/filepath"
)

type VideoController struct {
	controller.BaseController
	videoRepository repository.VideoRepository
	messageBroker   amqp.MessageBroker
	videoService    service.VideoService
}

func NewVideoController(videoRepository repository.VideoRepository, videoService *service.VideoService) *VideoController {
	v := new(VideoController)

	v.videoRepository = videoRepository
	v.videoService = *videoService
	v.messageBroker = amqp.NewRabbitMqService("guest", "guest", "localhost", 5672)
	if v.messageBroker == nil {
		return nil
	}
	return v
}

func (c VideoController) GetVideo() func(w http.ResponseWriter, r *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*request).Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		vars := mux.Vars(request)
		var videoId string
		var ok bool

		if videoId, ok = vars["videoId"]; !ok {
			c.BaseController.WriteResponse(&writer, http.StatusNotFound, false, "401 videoId not present")
			return
		}
		video, err := c.videoService.FindVideo(videoId)
		if err != nil {
			log.Error(err.Error())
			c.BaseController.WriteResponse(&writer, http.StatusNotFound, false, "401 videoId not present")
			return
		}

		c.BaseController.WriteResponseData(writer, video)
	}
}

func (c VideoController) GetVideos() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*request).Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		var page int
		page, err := c.GetIntRouteParameter(request, "page")
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Failed get page parameter")
			return
		}

		var countVideoOnPage int
		countVideoOnPage, err = c.GetIntRouteParameter(request, "countVideoOnPage")
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Failed get countVideoOnPage parameter")
			return
		}

		log.Info("page ", page, " count video ", countVideoOnPage)
		pageCount, ok := c.videoRepository.GetPageCount(countVideoOnPage)
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Failed get page countVideoOnPage")
			return
		}

		videos, err := c.videoRepository.FindVideosByPage(page, countVideoOnPage)
		if err != nil {
			log.Error(err.Error())
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, err.Error())
			return
		}

		responseData := make(map[string]interface{})
		responseData["pageCount"] = pageCount
		responseData["videos"] = videos

		c.BaseController.WriteResponseData(writer, responseData)
	}
}

func (c *VideoController) UploadVideo() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*request).Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		authorization := request.Header.Get("Authorization")
		userDto, ok := util.ValidateToken(authorization, "http://localhost:8001")
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusUnauthorized, false, "Not grant permission")
			return
		}

		title := request.FormValue("title")
		if len(title) == 0 {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Cannot get title")
			return
		}
		description := request.FormValue("description")
		if len(description) == 0 {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Cannot get description")
			return
		}
		fileReader, header, err := request.FormFile("file")
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Cannot get file")
			return
		}

		videoId, err := c.videoService.UploadVideo(userDto.UserId, title, description, fileReader, header)
		if err != nil {
			log.Error(err.Error())
			http.Error(writer, "Failed upload video", http.StatusInternalServerError)
			//c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Failed upload video")
			return
		}

		writer.WriteHeader(http.StatusOK)
		c.BaseController.JsonResponse(writer, videoId)
	}
}

func (c *VideoController) UpdateTitleAndDescription() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*request).Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		authorization := request.Header.Get("Authorization")
		userDto, ok := util.ValidateToken(authorization, "http://localhost:8001")
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusUnauthorized, false, "Not grant permission")
			return
		}

		vars := mux.Vars(request)

		var videoId string
		if videoId, ok = vars["videoId"]; !ok {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "videoId on path not found")
			return
		}

		var videoDto dto.VideoMetadata
		err := json.NewDecoder(request.Body).Decode(&videoDto)
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "cannot decode videoId|title|description struct")
			return
		}

		err = c.videoService.UpdateTitleAndDescription(userDto, videoId, videoDto)
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Failed update title and description")
			return
		}

		c.BaseController.WriteResponse(&writer, http.StatusOK, true, "success update title")
	}
}

func (c *VideoController) DeleteVideo() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*request).Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		authorization := request.Header.Get("Authorization")
		userDto, ok := util.ValidateToken(authorization, "http://localhost:8001")
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusUnauthorized, false, "Not grant permission")
			return
		}

		vars := mux.Vars(request)
		var videoId string

		if videoId, ok = vars["videoId"]; !ok {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "401 id not present")
			return
		}

		err := c.videoService.DeleteVideo(userDto, videoId)
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, err.Error())
			return
		}

		writer.WriteHeader(http.StatusOK)
		c.BaseController.WriteResponse(&writer, http.StatusOK, true, "success delete video with id "+videoId)
	}
}

func (c *VideoController) SearchVideo() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*request).Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		var page int
		page, err := c.GetIntRouteParameter(request, "page")
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "page parameter not present")
			return
		}
		var countVideoOnPage int
		countVideoOnPage, err = c.GetIntRouteParameter(request, "limit")
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "limit parameter not present")
			return
		}
		var search string
		search, ok := c.ParseRouteParameter(request, "search")
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "countVideoOnPage parameter not present")
			return
		}

		log.Info("page ", page, " count video ", countVideoOnPage)
		pageCount, ok := c.videoRepository.GetPageCount(countVideoOnPage)
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Failed get page countVideoOnPage")
			return
		}

		videos, err := c.videoRepository.SearchVideo(search, page, countVideoOnPage)
		if err != nil {
			log.Error(err.Error())
			c.BaseController.WriteResponse(&writer, http.StatusInternalServerError, false, "Video not found")
			return
		}

		c.responseVideoListItems(writer, pageCount, videos)
	}
}

func (c *VideoController) IncrementViews() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*request).Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		vars := mux.Vars(request)
		var id string
		var ok bool

		if id, ok = vars["videoId"]; !ok {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Id route parameter not present")
			return
		}

		responseData := make(map[string]interface{})
		responseData["success"] = c.videoRepository.IncrementViews(id)

		c.JsonResponse(writer, responseData)
	}
}

func (c *VideoController) responseVideoListItems(writer http.ResponseWriter, pageCount int, videos []model.VideoListItem) {
	responseData := make(map[string]interface{})
	responseData["pageCount"] = pageCount
	responseData["videos"] = videos

	c.BaseController.WriteResponseData(writer, responseData)
}

func (c VideoController) AddQuality() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*request).Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		authorization := request.Header.Get("Authorization")
		userDto, ok := util.ValidateToken(authorization, "http://localhost:8001")
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusUnauthorized, false, "Not grant permission")
			return
		}

		vars := mux.Vars(request)

		var videoId string
		if videoId, ok = vars["videoId"]; !ok {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "videoId on path not found")
			return
		}

		var quality dto.Quality
		err := json.NewDecoder(request.Body).Decode(&quality)

		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "cannot decode quality struct")
			return
		}

		err = c.videoService.AddQuality(videoId, userDto, quality)
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Failed add quality")
			return
		}

		c.BaseController.WriteResponse(&writer, http.StatusOK, true, "Success add quality")
	}
}
