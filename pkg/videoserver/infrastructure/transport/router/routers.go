package router

import (
	"github.com/bearname/videohost/pkg/common/amqp"
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/common/infrarstructure/caching"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/handler"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/middleware"
	"github.com/bearname/videohost/pkg/common/model"
	"github.com/bearname/videohost/pkg/videoserver/app/service"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/mysql"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/transport/controller"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Router(connector database.Connector, messageBrokerAddress string, authServerAddress string, redisAddress string, redisPassword string) http.Handler {

	videoRepository := mysql.NewMysqlVideoRepository(connector)
	messageBroker := amqp.NewRabbitMqService(messageBrokerAddress)
	if messageBroker == nil {
		return nil
	}

	cache := caching.NewRedisCache(model.NewDsn(redisAddress, "", redisPassword, ""))
	videoService := service.NewVideoService(videoRepository, messageBroker, cache)
	videoController := controller.NewVideoController(videoRepository, videoService, messageBroker, authServerAddress)
	if videoController == nil {
		log.Fatal("failed build videoController")
		return nil
	}

	streamController := controller.NewStreamController(service.NewStreamServiceImpl())

	router := mux.NewRouter()
	router.HandleFunc("/health", handler.HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/ready", handler.ReadyHandler).Methods(http.MethodGet)

	router.HandleFunc("/media/{videoId}/stream/", streamController.StreamHandler).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/media/{videoId}/{quality}/stream/", streamController.StreamHandler).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/media/{videoId}/{quality}/stream/{segName:index-[0-9]+.ts}", streamController.StreamHandler).Methods(http.MethodGet, http.MethodOptions)

	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subRouter.HandleFunc("/videos/", videoController.GetVideos()).Methods(http.MethodGet)
	subRouter.HandleFunc("/videos/search", videoController.SearchVideo()).Methods(http.MethodGet)
	subRouter.HandleFunc("/videos/{videoId}", videoController.GetVideo()).Methods(http.MethodGet)
	subRouter.HandleFunc("/videos/", videoController.UploadVideo()).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/videos/{videoId}", videoController.UpdateTitleAndDescription()).Methods(http.MethodPut, http.MethodOptions)
	subRouter.HandleFunc("/videos/{videoId}/add-quality", videoController.AddQuality()).Methods(http.MethodPut, http.MethodOptions)
	subRouter.HandleFunc("/videos/{videoId}", videoController.DeleteVideo()).Methods(http.MethodDelete, http.MethodOptions)
	subRouter.HandleFunc("/videos/{videoId}/increment", videoController.IncrementViews()).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/videos/{videoId}/like/{isLike:[0-1]}", middleware.AllowCors(middleware.AuthMiddleware(videoController.LikeVideo(), authServerAddress))).Methods(http.MethodPost, http.MethodOptions)

	var imgServer = http.FileServer(http.Dir("C:\\Users\\mikha\\go\\src\\videohost\\bin\\videoserver\\content"))
	router.PathPrefix("/content/").Handler(http.StripPrefix("/content/", imgServer))

	return middleware.LogMiddleware(router)
}
