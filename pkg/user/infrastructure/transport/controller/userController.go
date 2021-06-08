package controller

import (
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/controller"
	userService "github.com/bearname/videohost/pkg/user/app/service"
	"github.com/bearname/videohost/pkg/user/domain/repository"
	"github.com/bearname/videohost/pkg/videoserver/domain"
	dto2 "github.com/bearname/videohost/pkg/videoserver/domain/dto"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/transport/requestparser"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type UserController struct {
	controller.BaseController
	userRepository  repository.UserRepository
	userService     userService.UserService
	videoRepository domain.VideoRepository
}

func NewUserController(userRepository repository.UserRepository, videoRepository domain.VideoRepository) *UserController {
	v := new(UserController)
	v.userRepository = userRepository
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
	username, ok := vars["username"]
	if !ok {
		c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Cannot find username in request")
		return
	}

	if _, err := c.userRepository.FindByUserName(username); err != nil {
		c.BaseController.WriteResponse(&writer, http.StatusOK, false, "User not exist")
		return
	}

	userId, ok := context.Get(request, "userId").(string)
	if !ok {
		c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "Cannot check userId")
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
	countAllVideos, ok := c.userRepository.GetCountVideos(userId)
	if !ok {
		c.BaseController.WriteResponse(&writer, http.StatusOK, false, "Failed get page countVideoOnPage")
		return
	}

	videos, err := c.videoRepository.FindUserVideos(userId, searchDto)
	if err != nil {
		log.Error(err.Error())
		c.BaseController.WriteResponse(&writer, http.StatusOK, false, err.Error())
		return
	}

	responseData := make(map[string]interface{})
	responseData["countAllVideos"] = countAllVideos
	responseData["videos"] = videos

	c.BaseController.WriteResponseData(writer, responseData)
}
