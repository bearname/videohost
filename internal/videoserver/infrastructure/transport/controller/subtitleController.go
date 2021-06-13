package controller

import (
	"encoding/json"
	"github.com/bearname/videohost/internal/common/caching"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/controller"
	"github.com/bearname/videohost/internal/videoserver/domain"
	"github.com/bearname/videohost/internal/videoserver/domain/dto"
	"net/http"
	"strconv"
)

type SubtitleController struct {
	controller.BaseController
	cache             caching.Cache
	subtitleService   domain.SubtitleService
	authServerAddress string
}

func NewSubtitleController(playListService domain.SubtitleService, authServerAddress string) *SubtitleController {
	v := new(SubtitleController)

	v.subtitleService = playListService
	v.authServerAddress = authServerAddress
	return v
}

func (c *SubtitleController) CreateSubtitle() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var createRequest dto.CreateSubtitleRequestDto
		err := json.NewDecoder(request.Body).Decode(&createRequest)
		if err != nil {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "cannot decode name|privacy|videosId struct")
			return
		}

		playlistId, err := c.subtitleService.Create(createRequest)
		if err != nil {
			//TODO add translate method
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, err.Error())
			return
		}

		c.BaseController.WriteJsonResponse(writer, struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{1, "Success create playlist with id " + strconv.Itoa(int(playlistId))})
	}
}
