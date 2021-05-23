package router

import (
	"github.com/bearname/videohost/pkg/common/amqp"
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/handler"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/middleware"
	"github.com/bearname/videohost/pkg/videoserver/app/service"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/mysql"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/transport/controller"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Router(connector database.Connector) http.Handler {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	videoRepository := mysql.NewMysqlVideoRepository(connector)
	messageBroker := amqp.NewRabbitMqService("guest", "guest", "localhost", 5672)
	if messageBroker == nil {
		return nil
	}

	videoService := service.NewVideoService(videoRepository, messageBroker)
	videoController := controller.NewVideoController(videoRepository, videoService)
	if videoController == nil {
		log.Fatal("Failed")
		panic("failed build videocontroller")
	}

	streamController := controller.NewStreamController(videoRepository)
	router.HandleFunc("/media/{videoId}/{quality}/stream/", streamController.StreamHandler).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/media/{videoId}/{quality}/stream/{segName:index-[0-9]+.ts}", streamController.StreamHandler).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/health", handler.HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/ready", handler.ReadyHandler).Methods(http.MethodGet)

	subRouter.HandleFunc("/videos/", videoController.GetVideos()).Methods(http.MethodGet)
	subRouter.HandleFunc("/videos/search", videoController.SearchVideo()).Methods(http.MethodGet)
	subRouter.HandleFunc("/videos/{videoId}", videoController.GetVideo()).Methods(http.MethodGet)
	subRouter.HandleFunc("/videos/", videoController.UploadVideo()).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/videos/{videoId}", videoController.UpdateTitleAndDescription()).Methods(http.MethodPut, http.MethodOptions)
	subRouter.HandleFunc("/videos/{videoId}", videoController.DeleteVideo()).Methods(http.MethodDelete, http.MethodOptions)
	subRouter.HandleFunc("/videos/{videoId}/increment", videoController.IncrementViews()).Methods(http.MethodPost, http.MethodOptions)

	var imgServer = http.FileServer(http.Dir("C:\\Users\\mikha\\go\\src\\videohost\\bin\\videoserver\\content"))
	router.PathPrefix("/content/").Handler(http.StripPrefix("/content/", imgServer))

	return middleware.LogMiddleware(router)
}
