package router

import (
	"github.com/bearname/videohost/internal/common/infrarstructure/profile"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/handler"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/middleware"
	"github.com/bearname/videohost/internal/web-api-gateway/infrastructure/transport/controller"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
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
func Router(gatewayController *controller.GatewayController) http.Handler {

	router := mux.NewRouter()
	router.HandleFunc("/swagger/{id}", httpSwagger.Handler(httpSwagger.URL("http://localhost:8020/swagger/doc.json"))).Methods(http.MethodGet)

	router.HandleFunc("/health", handler.HealthHandler).Methods(http.MethodGet)
	router.HandleFunc("/ready", handler.ReadyHandler).Methods(http.MethodGet)

	r := mux.NewRouter()
	appRouter := r.PathPrefix("/api").Subrouter()
	appRouter.Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodOptions, http.MethodDelete).
		HandlerFunc(gatewayController.Handle)

	return middleware.LogMiddleware(profile.BuildHandlers(router))
}
