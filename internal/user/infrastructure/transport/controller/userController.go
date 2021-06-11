package controller

import (
	"fmt"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/controller"
	userService "github.com/bearname/videohost/internal/user/app/service"
	"github.com/bearname/videohost/internal/user/domain/repository"
	"github.com/bearname/videohost/internal/videoserver/domain"
	"github.com/bearname/videohost/internal/videoserver/domain/dto"
	"github.com/bearname/videohost/internal/videoserver/infrastructure/transport/requestparser"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type UserController struct {
	controller.BaseController
	userRepository  repository.UserRepo
	userService     userService.UserService
	videoRepository domain.VideoRepository
}

func NewUserController(userService userService.UserService, userRepository repository.UserRepo, videoRepository domain.VideoRepository) *UserController {
	v := new(UserController)
	v.userRepository = userRepository
	v.userService = userService
	v.videoRepository = videoRepository
	return v
}

func (c *UserController) GetUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	log.Println("get usernameOrId called")
	variable, err := c.ParseMuxVariable(request, []string{"usernameOrId"})
	if err != nil {
		http.Error(writer, "Cannot find usernameOrId in request", http.StatusBadRequest)
		return
	}

	usernameOrId := variable[0]

	userDto, err := c.userService.Find(usernameOrId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	c.WriteJsonResponse(writer, userDto)
}

func (c *UserController) UpdatePassword(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	err := c.userService.UpdatePassword(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	c.WriteJsonResponse(writer, struct {
		Result bool `json:"result"`
	}{Result: true})
}

func (c *UserController) Follow(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	userId := context.Get(request, "userId")
	context.Clear(request)
	if userId == nil {
		http.Error(writer, "userId not set", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(request)
	followerId := vars["followingToId"]

	isFollowing := request.URL.Query().Get("following")
	if len(isFollowing) == 0 {
		isFollowing = "true"
	}
	parseBool, err := strconv.ParseBool(isFollowing)
	if err != nil {
		http.Error(writer, "following get query parameter not set", http.StatusBadRequest)
		return
	}

	err = c.userService.Follow(followerId, fmt.Sprintf("%v", userId), parseBool)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	c.WriteJsonResponse(writer, struct {
		Result bool `json:"result"`
	}{Result: true})
}

func (c *UserController) GetUserSubscription(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	log.Println("get user videos called")
	vars := mux.Vars(request)
	userId, ok := vars["userId"]
	if !ok {
		c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "Cannot find userId in request")
		return
	}

	statistic, err := c.userService.GetUserStatistic(userId)
	if err != nil {
		log.Error(err.Error())
		c.BaseController.WriteResponse(writer, http.StatusOK, false, err.Error())
		return
	}

	c.WriteJsonResponse(writer, statistic)
}

func (c *UserController) GetUserVideos(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	log.Println("get user videos called")
	vars := mux.Vars(request)
	userId, ok := vars["userId"]
	if !ok {
		c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "Cannot find by userId in request")
		return
	}

	if _, err := c.userRepository.FindById(userId); err != nil {
		c.BaseController.WriteResponse(writer, http.StatusOK, false, "User not exist")
		return
	}

	parser := requestparser.NewCatalogVideoParser()
	result, err := parser.Parse(request)
	if err != nil {
		c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, err.Error())
		return
	}
	searchDto := result.(dto.SearchDto)

	log.Info("page ", searchDto.Page, " count video ", searchDto.Count)
	countAllVideos, ok := c.userRepository.GetCountVideos(userId)
	if !ok {
		c.BaseController.WriteResponse(writer, http.StatusOK, false, "Failed get page countVideoOnPage")
		return
	}

	videos, err := c.videoRepository.FindUserVideos(userId, searchDto)
	if err != nil {
		log.Error(err.Error())
		c.BaseController.WriteResponse(writer, http.StatusOK, false, err.Error())
		return
	}

	responseData := make(map[string]interface{})
	responseData["countAllVideos"] = countAllVideos
	responseData["videos"] = videos

	c.BaseController.WriteResponseData(writer, responseData)
}
