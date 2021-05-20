package router

import (
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/middleware"
	mysql2 "github.com/bearname/videohost/pkg/user/infrastructure/mysql"
	authHandler "github.com/bearname/videohost/pkg/user/infrastructure/transport/controller"
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
	videoController := controller.NewVideoController(videoRepository)
	if videoController == nil {
		log.Fatal("Failed")
		panic("failed build videocontroller")
	}

	userRepository := mysql2.NewMysqlUserRepository(connector)
	authController := authHandler.NewAuthController(userRepository)

	subRouter.HandleFunc("/list", videoController.GetVideoList()).Methods(http.MethodGet)
	subRouter.HandleFunc("/video/search", videoController.SearchVideo()).Methods(http.MethodGet)
	subRouter.HandleFunc("/video/{ID}", videoController.GetVideo()).Methods(http.MethodGet)
	subRouter.HandleFunc("/video", authController.CheckTokenHandler(videoController.UploadVideo())).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/video/{id}/increment", videoController.IncrementViews()).Methods(http.MethodPost, http.MethodOptions)

	streamController := controller.NewStreamController(videoRepository)
	subRouter.HandleFunc("/media/{id}/stream/", streamController.StreamHandler).Methods(http.MethodGet)
	subRouter.HandleFunc("/media/{id}/stream/{segName:index[0-9]+.ts}", streamController.StreamHandler).Methods(http.MethodGet)
	var imgServer = http.FileServer(http.Dir("C:\\Users\\mikha\\go\\src\\videohost\\bin\\videoserver\\content"))
	router.PathPrefix("/content/").Handler(http.StripPrefix("/content/", imgServer))

	return middleware.LogMiddleware(router)
}
