package controller

import (
	"encoding/json"
	"github.com/bearname/videohost/internal/common/caching"
	commonDto "github.com/bearname/videohost/internal/common/dto"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/controller"
	"github.com/bearname/videohost/internal/common/util"
	"github.com/bearname/videohost/internal/videoserver/domain"
	"github.com/bearname/videohost/internal/videoserver/domain/dto"
	"github.com/bearname/videohost/internal/videoserver/domain/model"
	"github.com/bearname/videohost/internal/videoserver/infrastructure/transport/requestparser"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type PlayListController struct {
	controller.BaseController
	cache             caching.Cache
	playListService   domain.PlayListService
	authServerAddress string
}

func NewPlayListController(playListService domain.PlayListService, authServerAddress string) *PlayListController {
	v := new(PlayListController)

	v.playListService = playListService
	v.authServerAddress = authServerAddress
	return v
}

func (c *PlayListController) CreatePlaylist() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		userId := context.Get(request, "userId").(string)
		context.Clear(request)
		if len(userId) == 0 {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "cannot get userId form context")
			return
		}
		var createPlayListRequest CreatePlayListRequest
		err := json.NewDecoder(request.Body).Decode(&createPlayListRequest)
		if err != nil {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "cannot decode name|privacy|videosId struct")
			return
		}
		playlistId, err := c.playListService.CreatePlaylist(dto.CreatePlaylistDto{
			Name:    createPlayListRequest.Name,
			OwnerId: userId,
			Privacy: createPlayListRequest.Privacy,
			Videos:  createPlayListRequest.VideosId})
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

func (c *PlayListController) GetUserPlaylists() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ownerId, ok := requestparser.GetStringQueryParameter(request, "ownerId")
		if !ok {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "query parameter ownerId not set")
			return
		}

		var privacyType []model.PrivacyType
		privacyType = append(privacyType, model.Public)
		privacyType = append(privacyType, model.Unlisted)
		authorization := request.Header.Get("Authorization")
		if len(authorization) != 0 {
			var userDto commonDto.UserDto

			userDto, ok = util.ValidateToken(authorization, c.authServerAddress)
			if ok && userDto.UserId == ownerId {
				privacyType = append(privacyType, model.Private)
			}
		}

		userPlaylists, err := c.playListService.FindUserPlaylists(ownerId, privacyType)
		if err != nil {
			//TODO add translate method
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, err.Error())
			return
		}

		c.BaseController.WriteJsonResponse(writer, userPlaylists)
	}
}

func (c *PlayListController) GetPlayList() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		playlistId, err := strconv.Atoi(vars["playlistId"])
		if err != nil {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "cannot get playlistId")
			return
		}

		playlist, err := c.playListService.FindPlaylist(playlistId)
		if err != nil {
			//TODO add translate method
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, err.Error())
			return
		}
		if playlist.Privacy == model.Private {
			authorization := request.Header.Get("Authorization")
			var userDto commonDto.UserDto

			var ok bool
			userDto, ok = util.ValidateToken(authorization, c.authServerAddress)
			if !ok && userDto.UserId != playlist.OwnerId {
				c.BaseController.WriteResponse(writer, http.StatusUnauthorized, false, "Not grant permission")
				return
			}
		}

		c.BaseController.WriteJsonResponse(writer, struct {
			Id      string            `json:"id"`
			Name    string            `json:"name"`
			OwnerId string            `json:"owner_id"`
			Created time.Time         `json:"created"`
			Privacy model.PrivacyType `json:"privacy"`
			Videos  string            `json:"videos"`
		}{playlist.Id, playlist.Name, playlist.OwnerId, playlist.Created, playlist.Privacy, playlist.VideoString})
	}
}

func (c *PlayListController) ModifyVideoToPlaylist() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		userId := context.Get(request, "userId").(string)
		context.Clear(request)
		if len(userId) == 0 {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "cannot get userId form context")
			return
		}
		vars := mux.Vars(request)
		playlistId, err := strconv.Atoi(vars["playlistId"])
		if err != nil {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "playlistId not set")
			return
		}

		var modificationRequest PlayListVideoModificationRequest
		err = json.NewDecoder(request.Body).Decode(&modificationRequest)
		if err != nil {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "cannot decode videos struct")
			return
		}
		action := modificationRequest.Action
		if len(action.String()) == 0 {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "invalid action")
			return
		}

		err = c.playListService.ModifyVideosOnPlaylist(playlistId, userId, modificationRequest.Videos, action)
		if err != nil {
			if err == domain.ErrPlaylistDuplicate {
				c.BaseController.WriteResponse(writer, http.StatusConflict, false, err.Error())
			} else {
				c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, err.Error())
			}

			return
		}

		c.BaseController.WriteResponse(writer, http.StatusOK, true, "success "+action.String()+" videos")
	}
}

func (c *PlayListController) ChangePrivacy() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		userId := context.Get(request, "userId").(string)
		context.Clear(request)
		if len(userId) == 0 {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "cannot get userId form context")
			return
		}
		vars := mux.Vars(request)
		playlistId, err := strconv.Atoi(vars["playlistId"])
		if err != nil {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "playlistId not set")
			return
		}
		privacyType, err := strconv.Atoi(vars["privacyType"])
		if err != nil {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "playlistId not set")
			return
		}

		err = c.playListService.ChangePrivacy(userId, playlistId, model.PrivacyType(privacyType))
		if err != nil {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, err.Error())
			return
		}

		c.BaseController.WriteResponse(writer, http.StatusOK, true, "success change playlist privacy with id "+strconv.Itoa(playlistId))
	}
}

func (c *PlayListController) DeletePlaylist() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		userId := context.Get(request, "userId").(string)
		context.Clear(request)
		if len(userId) == 0 {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "cannot get userId form context")
			return
		}
		vars := mux.Vars(request)
		playlistId, err := strconv.Atoi(vars["playlistId"])
		if err != nil {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, "playlistId not set")
			return
		}

		err = c.playListService.Delete(userId, playlistId)
		if err != nil {
			c.BaseController.WriteResponse(writer, http.StatusBadRequest, false, err.Error())
			return
		}

		c.BaseController.WriteResponse(writer, http.StatusOK, true, "success delete playlist with id "+strconv.Itoa(playlistId))
	}
}
