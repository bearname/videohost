package controller

import (
	"encoding/json"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/controller"
	dto2 "github.com/bearname/videohost/pkg/user/app/dto"
	service2 "github.com/bearname/videohost/pkg/user/app/service"
	"github.com/bearname/videohost/pkg/user/domain/repository"
	repository2 "github.com/bearname/videohost/pkg/videoserver/domain/repository"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserController struct {
	controller.BaseController
	userRepository  repository.UserRepository
	videoRepository repository2.VideoRepository
}

func NewUserController(userRepository repository.UserRepository, videoRepository repository2.VideoRepository) *UserController {
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

	log.Println("get username called")
	vars := mux.Vars(request)
	username, ok := vars["username"]
	if !ok {
		http.Error(writer, "Cannot find username in request", http.StatusBadRequest)
		return
	}
	if _, err := c.userRepository.FindByUserName(username); err != nil {
		http.Error(writer, "Cannot find username", http.StatusNotFound)
		return
	}

	c.JsonResponse(writer,
		struct {
			Username    string `json:"username"`
			Description string `json:"description"`
		}{username, ""})
}

func (c *UserController) UpdatePassword(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	var userDto dto2.ChangePasswordUserDto
	err := json.NewDecoder(request.Body).Decode(&userDto)
	if err != nil {
		http.Error(writer, "cannot decode username/password struct", http.StatusBadRequest)
		return
	}
	if !service2.IsUsernameContextOk(userDto.Username, request) {
		http.Error(writer, "Is username context invalid", http.StatusBadRequest)
		return
	}

	userFromDb, err := c.userRepository.FindByUserName(userDto.Username)
	if (err == nil && userFromDb.Username != userDto.Username) || err != nil {
		http.Error(writer, "User not exist", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword(userFromDb.Password, []byte(userDto.OldPassword))
	if err != nil {
		http.Error(writer, "Wrong password", http.StatusUnauthorized)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userDto.NewPassword), bcrypt.DefaultCost)

	if ok := c.userRepository.UpdatePassword(userDto.Username, passwordHash); !ok {
		http.Error(writer, "Failed update password", http.StatusUnauthorized)
		return
	}

	c.JsonResponse(writer, struct {
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

	var page int
	page, err := c.GetIntRouteParameter(request, "page")
	if err != nil {
		c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, err.Error())
		return
	}
	var countVideoOnPage int
	countVideoOnPage, err = c.GetIntRouteParameter(request, "countVideoOnPage")
	if err != nil {
		c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, err.Error())
		return
	}

	log.Info("page ", page, " count video ", countVideoOnPage)
	countAllVideos, ok := c.userRepository.GetCountVideos(userId)
	if !ok {
		c.BaseController.WriteResponse(&writer, http.StatusOK, false, "Failed get page countVideoOnPage")
		return
	}

	videos, err := c.videoRepository.FindUserVideos(userId, page, countVideoOnPage)
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
