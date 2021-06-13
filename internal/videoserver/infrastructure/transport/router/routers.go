package router

import (
	"github.com/bearname/videohost/internal/common/db"
	"github.com/bearname/videohost/internal/common/infrarstructure/amqp"
	"github.com/bearname/videohost/internal/common/infrarstructure/profile"
	caching "github.com/bearname/videohost/internal/common/infrarstructure/redis"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/handler"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/middleware"
	"github.com/bearname/videohost/internal/videoserver/app/service"
	"github.com/bearname/videohost/internal/videoserver/infrastructure/mysql"
	"github.com/bearname/videohost/internal/videoserver/infrastructure/transport/controller"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	"net/http"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /api/v1/
func Router(connector db.Connector, messageBrokerAddress string, authServerAddress string, redisAddress string, redisPassword string) http.Handler {

	videoRepository := mysql.NewMysqlVideoRepository(connector)
	messageBroker := amqp.NewRabbitMqService(messageBrokerAddress)
	if messageBroker == nil {
		return nil
	}

	cache := caching.NewRedisCache(db.NewDsn(redisAddress, "", redisPassword, ""))
	videoService := service.NewVideoService(videoRepository, messageBroker, cache)
	videoController := controller.NewVideoController(videoRepository, videoService, messageBroker, authServerAddress)
	if videoController == nil {
		log.Fatal("failed build videoController")
		return nil
	}

	router := mux.NewRouter()
	router.HandleFunc("/swagger/{id}", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"), //The url pointing to API definition"
	)).Methods(http.MethodGet)

	router.HandleFunc("/health", handler.HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/ready", handler.ReadyHandler).Methods(http.MethodGet)

	streamController := controller.NewStreamController(service.NewStreamServiceImpl())
	router.HandleFunc("/media/{videoId}/stream/", streamController.StreamHandler).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/media/{videoId}/{quality}/stream/", streamController.StreamHandler).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/media/{videoId}/{quality}/stream/{segName:index-[0-9]+.ts}", streamController.StreamHandler).Methods(http.MethodGet, http.MethodOptions)

	subRouter := router.PathPrefix("/api/v1").Subrouter()
	subtitleRepo := mysql.NewSubtitleRepository(connector)
	subtitleService := service.NewSubtitleService(subtitleRepo, cache)
	subtitleController := controller.NewSubtitleController(subtitleService, authServerAddress)
	subRouter.HandleFunc("/subtitles", middleware.AllowCors(middleware.AuthMiddleware(subtitleController.CreateSubtitle(), authServerAddress))).Methods(http.MethodPost, http.MethodOptions)

	repository := mysql.NewPlaylistRepository(connector)
	listService := service.NewPlayListService(repository, cache)
	playListController := controller.NewPlayListController(listService, authServerAddress)

	subRouter.HandleFunc("/playlists", middleware.AllowCors(middleware.AuthMiddleware(playListController.CreatePlaylist(), authServerAddress))).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/playlists", middleware.AllowCors(playListController.GetUserPlaylists())).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/playlists/{playlistId}", middleware.AllowCors(playListController.GetPlayList())).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/playlists/{playlistId}/modify", middleware.AuthMiddleware(playListController.ModifyVideoToPlaylist(), authServerAddress)).Methods(http.MethodPut, http.MethodOptions)
	subRouter.HandleFunc("/playlists/{playlistId}/change-privacy/{privacyType}", middleware.AllowCors(middleware.AuthMiddleware(playListController.ChangePrivacy(), authServerAddress))).Methods(http.MethodPut, http.MethodOptions)
	subRouter.HandleFunc("/playlists/{playlistId}", middleware.AllowCors(middleware.AuthMiddleware(playListController.DeletePlaylist(), authServerAddress))).Methods(http.MethodDelete, http.MethodOptions)

	subRouter.HandleFunc("/videos/", videoController.GetVideos()).Methods(http.MethodGet)
	subRouter.HandleFunc("/videos/search", videoController.SearchVideo()).Methods(http.MethodGet)
	subRouter.HandleFunc("/videos/{videoId}", videoController.GetVideo()).Methods(http.MethodGet)
	subRouter.HandleFunc("/videos/", videoController.UploadVideo()).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/videos/{videoId}", videoController.UpdateTitleAndDescription()).Methods(http.MethodPut, http.MethodOptions)
	subRouter.HandleFunc("/videos/{videoId}/add-quality", videoController.AddQuality()).Methods(http.MethodPut, http.MethodOptions)
	subRouter.HandleFunc("/videos/{videoId}", videoController.DeleteVideo()).Methods(http.MethodDelete, http.MethodOptions)
	subRouter.HandleFunc("/videos/{videoId}/increment", videoController.IncrementViews()).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/videos-liked", middleware.AllowCors(middleware.AuthMiddleware(videoController.FindUserLikedVideo(), authServerAddress))).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/videos/{videoId}/like/{isLike:[0-1]}", middleware.AllowCors(middleware.AuthMiddleware(videoController.LikeVideo(), authServerAddress))).Methods(http.MethodPost, http.MethodOptions)

	var imgServer = http.FileServer(http.Dir("C:\\Users\\mikha\\go\\src\\videohost\\bin\\videoserver\\content"))
	router.PathPrefix("/content/").Handler(http.StripPrefix("/content/", imgServer))

	return middleware.LogMiddleware(profile.ProfileHandler(router))
}
