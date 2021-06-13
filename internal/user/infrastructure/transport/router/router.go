package router

import (
	"github.com/bearname/videohost/internal/common/db"
	"github.com/bearname/videohost/internal/common/infrarstructure/profile"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/handler"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/middleware"
	"github.com/bearname/videohost/internal/user/app/service"
	userRepo "github.com/bearname/videohost/internal/user/infrastructure/mysql"
	"github.com/bearname/videohost/internal/user/infrastructure/transport/controller"
	"github.com/bearname/videohost/internal/videoserver/infrastructure/mysql"
	"github.com/gorilla/mux"
	"net/http"
)

func Router(connector db.Connector) http.Handler {

	videoRepo := mysql.NewMysqlVideoRepository(connector)
	userRepository := userRepo.NewMysqlUserRepository(connector)
	followingRepository := userRepo.NewFollowerRepo(connector)
	authService := service.NewAuthService(userRepository)
	authController := controller.NewAuthController(authService)
	userService := service.NewUserService(userRepository, followingRepository)
	userController := controller.NewUserController(*userService, userRepository, videoRepo)

	router := mux.NewRouter()

	router.HandleFunc("/health", handler.HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/ready", handler.ReadyHandler).Methods(http.MethodGet)

	apiV1Route := router.PathPrefix("/api/v1").Subrouter()
	apiV1Route.HandleFunc("/auth/create-user", authController.CreateUser).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/auth/login", authController.Login).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/auth/token", authController.CheckTokenHandler(authController.RefreshToken)).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/auth/token/validate", authController.ValidateToken).Methods(http.MethodGet, http.MethodOptions)

	apiV1Route.HandleFunc("/users/updatePassword", authController.CheckTokenHandler(userController.UpdatePassword)).Methods(http.MethodPut, http.MethodOptions)
	apiV1Route.HandleFunc("/users/{usernameOrId}", authController.CheckTokenHandler(userController.GetUser)).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/users/{userId}/subscriber", userController.GetUserSubscription).Methods(http.MethodGet, http.MethodOptions)
	apiV1Route.HandleFunc("/users/{userId}/follow", authController.CheckTokenHandler(userController.Follow)).Methods(http.MethodPost, http.MethodOptions)
	apiV1Route.HandleFunc("/users/{userId}/videos", userController.GetUserVideos).Methods(http.MethodGet, http.MethodOptions)

	return middleware.LogMiddleware(profile.BuildHandlers(router))
}
