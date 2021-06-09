package transport

import (
	"encoding/json"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/controller"
	"github.com/bearname/videohost/internal/common/util"
	"github.com/bearname/videohost/internal/video-comments/app/service"
	"github.com/bearname/videohost/internal/video-comments/domain"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type CommentController struct {
	controller.BaseController
	commentService    service.CommentService
	authServerAddress string
}

func NewCommentController(commentService *service.CommentService, authServerAddress string) *CommentController {
	v := new(CommentController)

	v.commentService = *commentService
	v.authServerAddress = authServerAddress

	return v
}

func (c *CommentController) CreateComment() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		authorization := request.Header.Get("Authorization")
		userDto, ok := util.ValidateToken(authorization, c.authServerAddress)
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusUnauthorized, false, "Not grant permission")
			return
		}

		var commentRequest AddCommentRequest
		err := json.NewDecoder(request.Body).Decode(&commentRequest)
		if err != nil {
			log.Error(err)
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "cannot decode videoId|message struct")
			return
		}

		commentId, err := c.commentService.Create(&domain.CommentDto{
			UserId:   userDto.UserId,
			VideoId:  commentRequest.VideoId,
			ParentId: commentRequest.ParentId,
			Message:  commentRequest.VideoId,
		})
		if err != nil {
			log.Error(err)
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "failed create comments")
			return
		}

		c.BaseController.WriteResponse(&writer, http.StatusOK, true, "success create comment "+strconv.FormatInt(commentId, 10))
	}
}

func (c *CommentController) FindComments() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		commentsFilter, err := DecodeFindCommentsRequest(request)
		if err != nil {
			log.Error(err)
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "parse query parameter")
			return
		}

		if len(commentsFilter.AuthorId) != 0 {
			var comments domain.Comments
			comments, err = c.commentService.FindUserComments(commentsFilter.AuthorId, &commentsFilter.Page)
			if err != nil {
				log.Error(err)
				c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "failed find user comments")
				return
			}

			c.WriteJsonResponse(writer, comments)
			return
		}
		if len(commentsFilter.VideoId) != 0 {
			var comments domain.VideoComments
			comments, err = c.commentService.FindRootLevel(commentsFilter.VideoId, &commentsFilter.Page)
			if err != nil {
				log.Error(err)
				c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "failed find comments")
				return
			}

			c.WriteJsonResponse(writer, comments)
			return
		}
		if commentsFilter.RootId != NotSet {
			var comments []domain.Comment
			comments, err = c.commentService.FindChildren(commentsFilter.RootId, &commentsFilter.Page)
			if err != nil {
				log.Error(err)
				c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "failed find comments")
				return
			}

			c.WriteJsonResponse(writer, struct {
				Comments []domain.Comment `json:"comments"`
			}{Comments: comments})
		}

		err = controller.ErrRouteNotFound
		c.WriteError(writer, err, controller.TransportError{
			Status: http.StatusNotFound,
			Response: controller.Response{
				Code:    100,
				Message: err.Error(),
			},
		})
	}
}

func (c *CommentController) EditComment() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		authorization := request.Header.Get("Authorization")
		userDto, ok := util.ValidateToken(authorization, c.authServerAddress)
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusUnauthorized, false, "Not grant permission")
			return
		}

		var commentRequest EditCommentRequest
		err := json.NewDecoder(request.Body).Decode(&commentRequest)
		if err != nil {
			log.Error(err)
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "cannot decode videoId|message struct")
			return
		}

		vars := mux.Vars(request)
		commentId, err := strconv.Atoi(vars["commentId"])
		if err != nil {
			log.Error(err)
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "failed parse commentId")
			return
		}

		err = c.commentService.Edit(commentId, commentRequest.Message, userDto.UserId)
		if err != nil {
			log.Error(err)
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "failed edit comments")
			return
		}

		c.BaseController.WriteResponse(&writer, http.StatusOK, true, "success edit comment "+strconv.FormatInt(int64(commentId), 10))
	}
}

func (c *CommentController) DeleteComment() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		authorization := request.Header.Get("Authorization")
		userDto, ok := util.ValidateToken(authorization, c.authServerAddress)
		if !ok {
			c.BaseController.WriteResponse(&writer, http.StatusUnauthorized, false, "Not grant permission")
			return
		}

		vars := mux.Vars(request)
		commentId, err := strconv.Atoi(vars["commentId"])
		if err != nil {
			log.Error(err)
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "failed parse commentId")
			return
		}

		err = c.commentService.Delete(commentId, userDto.UserId)
		if err != nil {
			log.Error(err)
			c.BaseController.WriteResponse(&writer, http.StatusBadRequest, false, "failed find comments")
			return
		}

		c.BaseController.WriteResponse(&writer, http.StatusOK, true, "success delete comment "+strconv.FormatInt(int64(commentId), 10))
	}
}
