package transport

import (
	"github.com/bearname/videohost/internal/common/db"
	caching "github.com/bearname/videohost/internal/common/infrarstructure/redis"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/handler"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/middleware"
	"github.com/bearname/videohost/internal/video-comments/app/service"
	"github.com/bearname/videohost/internal/video-comments/infrastructure/mysql"

	"github.com/gorilla/mux"
	"net/http"
)

func Router(connector db.Connector, userServerAddress string, videoServerAddress string, redisAddress string, redisPassword string) http.Handler {
	commentRepo := mysql.NewCommentRepo(connector)
	cache := caching.NewRedisCache(db.NewDsn(redisAddress, "", redisPassword, ""))
	commentService := service.NewCommentService(commentRepo, cache, videoServerAddress)
	videoController := NewCommentController(commentService, userServerAddress)

	router := mux.NewRouter()
	router.HandleFunc("/health", handler.HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/ready", handler.ReadyHandler).Methods(http.MethodGet)

	subRouter := router.PathPrefix("/api/v1/comments").Subrouter()
	subRouter.HandleFunc("", middleware.AllowCors(videoController.FindComments())).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("", middleware.AllowCors(videoController.CreateComment())).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/{commentId}", middleware.AllowCors(videoController.EditComment())).Methods(http.MethodPut, http.MethodOptions)
	subRouter.HandleFunc("/{commentId}", middleware.AllowCors(videoController.DeleteComment())).Methods(http.MethodDelete, http.MethodOptions)

	return middleware.LogMiddleware(router)
}
