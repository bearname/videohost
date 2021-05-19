package controller

import (
	"github.com/bearname/videohost/videoserver/repository"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type UserController struct {
	BaseController
	userRepository *repository.UserRepository
}

func NewUserController(userRepository *repository.UserRepository) *UserController {
	v := new(UserController)

	v.userRepository = userRepository
	return v
}

func (c *UserController) GetUser(writer http.ResponseWriter, request *http.Request) {
	writer = *c.BaseController.AllowCorsRequest(&writer)
	if (*request).Method == "OPTIONS" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	log.Println("get user called")
	vars := mux.Vars(request)
	user, ok := vars["USERNAME"]
	if !ok {
		http.Error(writer, "Cannot find username in request", http.StatusBadRequest)
		return
	}
	if _, ok := c.userRepository.Users[user]; !ok {
		http.Error(writer, "Cannot find user", http.StatusNotFound)

		return
	}

	c.JsonResponse(writer,
		struct {
			Username    string `json:"username"`
			Description string `json:"description"`
		}{user, ""})
}
