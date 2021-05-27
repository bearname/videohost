package router

import (
	mysql2 "github.com/bearname/videohost/pkg/common/infrarstructure/mysql"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/handler"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/middleware"
	userRepo "github.com/bearname/videohost/pkg/user/infrastructure/mysql"
	"github.com/bearname/videohost/pkg/user/infrastructure/transport/controller"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/mysql"
	"github.com/gorilla/mux"
	"net/http"
)

func Router(connector mysql2.MysqlConnector) http.Handler {
	router := mux.NewRouter()
	apiV1Route := router.PathPrefix("/api/v1").Subrouter()

	videoRepository := mysql.NewMysqlVideoRepository(connector)
	userRepository := userRepo.NewMysqlUserRepository(connector)
	authController := controller.NewAuthController(userRepository)
	userController := controller.NewUserController(userRepository, videoRepository)

	router.HandleFunc("/health", handler.HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/ready", handler.ReadyHandler).Methods(http.MethodGet)

	apiV1Route.HandleFunc("/auth/create-user", authController.CreateUser).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/auth/login", authController.GetTokenUserPassword).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/auth/token", authController.CheckTokenHandler(authController.GetTokenByToken)).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/auth/token/validate", authController.CheckToken).Methods(http.MethodGet, http.MethodOptions)

	apiV1Route.HandleFunc("/users/updatePassword", authController.CheckTokenHandler(userController.UpdatePassword)).Methods(http.MethodPut, http.MethodOptions)
	apiV1Route.HandleFunc("/users/{usernameOrId}", authController.CheckTokenHandler(userController.GetUser)).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/users/{username}/videos", authController.CheckTokenHandler(userController.GetUserVideos)).Methods(http.MethodGet, http.MethodOptions)

	return middleware.LogMiddleware(router)
}
