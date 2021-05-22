package router

import (
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/handler"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/middleware"
	userRepo "github.com/bearname/videohost/pkg/user/infrastructure/mysql"
	"github.com/bearname/videohost/pkg/user/infrastructure/transport/controller"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/mysql"
	"github.com/gorilla/mux"
	"net/http"
)

func Router(connector database.Connector) http.Handler {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	videoRepository := mysql.NewMysqlVideoRepository(connector)
	userRepository := userRepo.NewMysqlUserRepository(connector)
	authController := controller.NewAuthController(userRepository)
	userController := controller.NewUserController(userRepository, videoRepository)

	router.HandleFunc("/health", handler.HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/ready", handler.ReadyHandler).Methods(http.MethodGet)

	subRouter.HandleFunc("/auth/create-user", authController.CreateUser).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/auth/login", authController.GetTokenUserPassword).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/auth/token", authController.CheckTokenHandler(authController.GetTokenByToken)).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/auth/token/validate", authController.CheckToken).Methods(http.MethodGet, http.MethodOptions)

	subRouter.HandleFunc("/users/updatePassword", authController.CheckTokenHandler(userController.UpdatePassword)).Methods(http.MethodPut, http.MethodOptions)
	subRouter.HandleFunc("/users/{username}", authController.CheckTokenHandler(userController.GetUser)).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/users/{username}/videos", authController.CheckTokenHandler(userController.GetUserVideos)).Methods(http.MethodGet, http.MethodOptions)

	return middleware.LogMiddleware(router)
}
