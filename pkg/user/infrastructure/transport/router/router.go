package router

import (
	"github.com/bearname/videohost/pkg/common/database"
	middleware2 "github.com/bearname/videohost/pkg/common/infrarstructure/transport/middleware"
	mysql2 "github.com/bearname/videohost/pkg/user/infrastructure/mysql"
	"github.com/bearname/videohost/pkg/user/infrastructure/transport/controller"
	"github.com/bearname/videohost/pkg/videoserver/infrastructure/mysql"
	"github.com/gorilla/mux"
	"net/http"
)

func Router(connector database.Connector) http.Handler {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	videoRepository := mysql.NewMysqlVideoRepository(connector)
	userRepository := mysql2.NewMysqlUserRepository(connector)
	authController := controller.NewAuthController(userRepository)
	userController := controller.NewUserController(userRepository, videoRepository)

	subRouter.HandleFunc("/auth/create-user", authController.CreateUser).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/auth/login", authController.GetTokenUserPassword).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/auth/token", authController.CheckTokenHandler(authController.GetTokenByToken)).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/auth/token/validate", authController.CheckToken).Methods(http.MethodGet, http.MethodOptions)

	subRouter.HandleFunc("/users/updatePassword", authController.CheckTokenHandler(userController.UpdatePassword)).Methods(http.MethodPut, http.MethodOptions)
	subRouter.HandleFunc("/users/{USERNAME}", authController.CheckTokenHandler(userController.GetUser)).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/users/{USERNAME}/videos", authController.CheckTokenHandler(userController.GetUserVideos)).Methods(http.MethodGet, http.MethodOptions)

	return middleware2.LogMiddleware(router)
}
