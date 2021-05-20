package router

import (
	"github.com/bearname/videohost/videoserver/controller"
	"github.com/bearname/videohost/videoserver/repository/mysql"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Router(connector mysql.Connector) http.Handler {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	videoRepository := mysql.NewMysqlVideoRepository(connector)
	videoController := controller.NewVideoController(videoRepository)
	if videoController == nil {
		log.Fatal("Failed")
		panic("failed build videocontroller")
	}

	userRepository := mysql.NewMysqlUserRepository(connector)
	authController := controller.NewAuthController(userRepository)
	userController := controller.NewUserController(userRepository, videoRepository)

	subRouter.HandleFunc("/list", videoController.GetVideoList()).Methods(http.MethodGet)
	subRouter.HandleFunc("/video/search", videoController.SearchVideo()).Methods(http.MethodGet)
	subRouter.HandleFunc("/video/{ID}", videoController.GetVideo()).Methods(http.MethodGet)
	subRouter.HandleFunc("/video", authController.CheckTokenHandler(videoController.UploadVideo())).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/video/{id}/increment", videoController.IncrementViews()).Methods(http.MethodPost, http.MethodOptions)

	subRouter.HandleFunc("/auth/create-user", authController.CreateUser).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/auth/login", authController.GetTokenUserPassword).Methods(http.MethodPost, http.MethodOptions)
	//subRouter.HandleFunc("/auth/logout",authController.CheckTokenAuthHandler(logout)).Methods("POST")

	subRouter.HandleFunc("/auth/token", authController.CheckTokenHandler(authController.GetTokenByToken)).Methods(http.MethodGet, http.MethodOptions)

	subRouter.HandleFunc("/users/updatePassword", authController.CheckTokenHandler(userController.UpdatePassword)).Methods(http.MethodPut, http.MethodOptions)
	subRouter.HandleFunc("/users/{USERNAME}", authController.CheckTokenHandler(userController.GetUser)).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/users/{USERNAME}/videos", authController.CheckTokenHandler(userController.GetUserVideos)).Methods(http.MethodGet, http.MethodOptions)

	streamController := controller.NewStreamController(videoRepository)
	subRouter.HandleFunc("/media/{id}/stream/", streamController.StreamHandler).Methods(http.MethodGet)
	subRouter.HandleFunc("/media/{id}/stream/{segName:index[0-9]+.ts}", streamController.StreamHandler).Methods(http.MethodGet)
	var imgServer = http.FileServer(http.Dir("C:\\Users\\mikha\\go\\src\\videoserver\\videoserver\\content"))
	router.PathPrefix("/content/").Handler(http.StripPrefix("/content/", imgServer))

	return logMiddleware(router)
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})
}
