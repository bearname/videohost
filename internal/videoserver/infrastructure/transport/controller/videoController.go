package controller

import (
	"encoding/json"
	"github.com/bearname/videohost/internal/common/infrarstructure/amqp"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/controller"
	"github.com/bearname/videohost/internal/common/util"
	"github.com/bearname/videohost/internal/videoserver/app/service"
	"github.com/bearname/videohost/internal/videoserver/domain"
	dto2 "github.com/bearname/videohost/internal/videoserver/domain/dto"
	"github.com/bearname/videohost/internal/videoserver/domain/model"
	"github.com/bearname/videohost/internal/videoserver/infrastructure/transport/requestparser"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type VideoController struct {
	controller.BaseController
	videoRepository   domain.VideoRepository
	messageBroker     amqp.MessageBroker
	videoService      service.VideoServiceImpl
	authServerAddress string
}

func NewVideoController(videoRepository domain.VideoRepository,
	videoService *service.VideoServiceImpl,
	messageBroker amqp.MessageBroker,
	authServerAddress string,
) *VideoController {
	v := new(VideoController)

	v.videoRepository = videoRepository
	v.videoService = *videoService
	v.messageBroker = messageBroker
	v.authServerAddress = authServerAddress
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

		//TODO send needed field on get parameter
		result, err := c.BaseController.ParseMuxVariable(request, []string{"videoId"})
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "401 id not present")
			return
		}
		video, err := c.videoService.FindVideo(result[0])
		if err != nil {
			log.Error(err.Error())
			c.BaseController.WriteResponse(&writer, http.StatusNotFound, false, "404 video not found")
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

		parser := requestparser.NewCatalogVideoParser()
		result, err := parser.Parse(request)
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, err.Error())
			return
		}
		searchDto := result.(dto2.SearchDto)

		log.Info("page ", searchDto.Page, " count video ", searchDto.Count)
		onPage, err := c.videoService.FindVideoOnPage(&searchDto)
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, err.Error())
			return
		}

		responseData := make(map[string]interface{})
		responseData["pageCount"] = onPage.PageCount
		responseData["videos"] = onPage.Videos

		c.BaseController.WriteResponseData(writer, responseData)
	}
}

func (c *VideoController) UploadVideo() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*request).Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		authorization := request.Header.Get("Authorization")
		userDto, ok := util.ValidateToken(authorization, c.authServerAddress)
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusUnauthorized, false, "Not grant permission")
			return
		}

		parser := requestparser.NewUploadVideoRequestParser()
		uploadVideoDto, err := parser.Parse(request)
		if err != nil {
			log.Error(err.Error())
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		videoDto := uploadVideoDto.(*dto2.UploadVideoDto)
		videoId, err := c.videoService.UploadVideo(userDto.UserId, videoDto)
		if err != nil {
			log.Error(err.Error())
			http.Error(writer, "Failed upload http", http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
		c.BaseController.WriteJsonResponse(writer, videoId)
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
		userDto, ok := util.ValidateToken(authorization, c.authServerAddress)
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusUnauthorized, false, "Not grant permission")
			return
		}

		result, err := c.BaseController.ParseMuxVariable(request, []string{"videoId"})
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "401 id not present")
			return
		}

		var videoDto dto2.VideoMetadata
		err = json.NewDecoder(request.Body).Decode(&videoDto)
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "cannot decode videoId|title|description struct")
			return
		}

		err = c.videoService.UpdateTitleAndDescription(userDto, result[0], videoDto)
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
		userDto, ok := util.ValidateToken(authorization, c.authServerAddress)
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusUnauthorized, false, "Not grant permission")
			return
		}
		result, err := c.BaseController.ParseMuxVariable(request, []string{"videoId"})
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "401 id not present")
			return
		}

		videoId := result[0]
		err = c.videoService.DeleteVideo(userDto, videoId)
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, err.Error())
			return
		}

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

		parser := requestparser.NewSearchVideoParser()
		result, err := parser.Parse(request)
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, err.Error())
			return
		}
		searchDto := result.(dto2.SearchDto)

		log.Info("page ", searchDto.Page, " count video ", searchDto.Count)
		pageCount, ok := c.videoRepository.GetPageCount(searchDto.Count)
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Failed get page countVideoOnPage")
			return
		}

		videos, err := c.videoRepository.SearchVideo(searchDto.SearchString, searchDto.Page, searchDto.Count)
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

		result, err := c.BaseController.ParseMuxVariable(request, []string{"videoId"})
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "401 id not present")
			return
		}

		responseData := make(map[string]interface{})
		responseData["success"] = c.videoRepository.IncrementViews(result[0])

		c.WriteJsonResponse(writer, responseData)
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
		userDto, ok := util.ValidateToken(authorization, c.authServerAddress)
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusUnauthorized, false, "Not grant permission")
			return
		}

		result, err := c.BaseController.ParseMuxVariable(request, []string{"videoId"})
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "401 id not present")
			return
		}

		var quality model.Quality
		err = json.NewDecoder(request.Body).Decode(&quality)
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "cannot decode quality struct")
			return
		}

		err = c.videoService.AddQuality(result[0], userDto, quality)
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Failed add quality")
			return
		}

		c.BaseController.WriteResponse(&writer, http.StatusOK, true, "Success add quality")
	}
}

func (c VideoController) LikeVideo() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		parser := requestparser.NewLikeVideoRequestParser()
		likeVideo, err := parser.Parse(request)
		if err != nil {
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, err.Error())
			return
		}
		likeVideoReq := likeVideo.(*requestparser.LikeVideoRequest)
		action, err := c.videoService.LikeVideo(model.Like{IdVideo: likeVideoReq.VideoId, OwnerId: likeVideoReq.OwnerId, IsLike: likeVideoReq.IsLike})
		if err != nil {
			c.BaseController.WriteError(writer, err, TranslateError(err))
			return
		}

		c.BaseController.WriteJsonResponse(writer, controller.Response{
			Code:    int(action),
			Message: "Success " + model.ActionToString(action),
		})
	}
}
