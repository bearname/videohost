package controller

import (
	"encoding/json"
	"github.com/bearname/videohost/videoserver/dto"
	"github.com/bearname/videohost/videoserver/repository"
	"github.com/bearname/videohost/videoserver/util"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserController struct {
	BaseController
	userRepository  repository.UserRepository
	videoRepository repository.VideoRepository
}

func NewUserController(userRepository repository.UserRepository, videoRepository repository.VideoRepository) *UserController {
	v := new(UserController)
	v.userRepository = userRepository
	v.videoRepository = videoRepository
	return v
}

func (c *UserController) GetUser(writer http.ResponseWriter, request *http.Request) {
	writer = *c.BaseController.AllowCorsRequest(&writer)
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	log.Println("get username called")
	vars := mux.Vars(request)
	username, ok := vars["USERNAME"]
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
	writer = *c.BaseController.AllowCorsRequest(&writer)
	var userDto dto.ChangePasswordUserDto
	err := json.NewDecoder(request.Body).Decode(&userDto)
	if err != nil {
		http.Error(writer, "cannot decode username/password struct", http.StatusBadRequest)
		return
	}
	if !util.IsUsernameContextOk(userDto.Username, request) {
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
	username, ok := vars["USERNAME"]
	if !ok {
		http.Error(writer, "Cannot find username in request", http.StatusBadRequest)
		return
	}
	if _, err := c.userRepository.FindByUserName(username); err != nil {
		http.Error(writer, "Cannot find username", http.StatusNotFound)
		return
	}

	userId, ok := context.Get(request, "userId").(string)
	if !ok {
		http.Error(writer, "Cannot check userId", http.StatusInternalServerError)
		return
	}

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
	countAllVideos, ok := c.userRepository.GetCountVideos(userId)
	if !ok {
		http.Error(writer, "Failed get page countVideoOnPage", http.StatusInternalServerError)
		return
	}

	videos, err := c.videoRepository.FindUserVideos(userId, page, countVideoOnPage)
	if err != nil {
		log.Error(err.Error())
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := make(map[string]interface{})
	responseData["countAllVideos"] = countAllVideos
	responseData["videos"] = videos

	c.writeResponseData(writer, responseData)
}
