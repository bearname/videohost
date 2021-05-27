package router

import (
	"github.com/bearname/videohost/pkg/common/database"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/handler"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport/middleware"
	"github.com/bearname/videohost/pkg/user/app/service"
	userRepo "github.com/bearname/videohost/pkg/user/infrastructure/mysql"
	"github.com/bearname/videohost/pkg/user/infrastructure/transport/controller"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/mysql"
	"github.com/gorilla/mux"
	"net/http"
)

func Router(connector database.Connector) http.Handler {
	router := mux.NewRouter()
	apiV1Route := router.PathPrefix("/api/v1").Subrouter()

	videoRepo := mysql.NewMysqlVideoRepository(connector)
	userRepository := userRepo.NewMysqlUserRepository(connector)
	authService := service.NewAuthService(userRepository)
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userRepository, videoRepo)

	router.HandleFunc("/health", handler.HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/ready", handler.ReadyHandler).Methods(http.MethodGet)

	apiV1Route.HandleFunc("/auth/create-user", authController.CreateUser).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/auth/login", authController.Login).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/auth/token", authController.CheckTokenHandler(authController.RefreshToken)).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/auth/token/validate", authController.ValidateToken).Methods(http.MethodGet, http.MethodOptions)

	apiV1Route.HandleFunc("/users/updatePassword", authController.CheckTokenHandler(userController.UpdatePassword)).Methods(http.MethodPut, http.MethodOptions)
	apiV1Route.HandleFunc("/users/{usernameOrId}", authController.CheckTokenHandler(userController.GetUser)).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/users/{username}/videos", authController.CheckTokenHandler(userController.GetUserVideos)).Methods(http.MethodGet, http.MethodOptions)

	return middleware.LogMiddleware(router)
}
