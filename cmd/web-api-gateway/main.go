package main

import (
	"fmt"
	"github.com/bearname/videohost/internal/web-api-gateway/domain"
	"github.com/bearname/videohost/internal/web-api-gateway/infrastructure/transport/controller"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	urlMapping := make(map[string]string)
	urlMapping["videos"] = "http://localhost:8000"
	urlMapping["playlists"] = "http://localhost:8000"
	urlMapping["subtitles"] = "http://localhost:8000"
	urlMapping["comments"] = "http://localhost:8010"
	urlMapping["auth"] = "http://localhost:8001"
	urlMapping["users"] = "http://localhost:8001"
	urlMapping["media"] = "http://localhost:8020"
	urlMapping["content"] = "http://localhost:8020"

	mapping := domain.NewUrlMapping(urlMapping)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8050"
	}

	gatewayController := controller.NewGatewayController(mapping)
	r := mux.NewRouter()
	appRouter := r.PathPrefix("/api").Subrouter()
	appRouter.Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodOptions, http.MethodDelete).
		HandlerFunc(gatewayController.Handle)

	fmt.Println("start at :" + port)
	err := http.ListenAndServe(":"+port, r)
	log.Fatal(err)
}
