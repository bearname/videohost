package controller

import (
	"errors"
	"fmt"
	"github.com/bearname/videohost/internal/common/infrarstructure/transport/controller"
	"github.com/bearname/videohost/internal/web-api-gateway/domain"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

type GatewayController struct {
	controller.BaseController
	urlMapping *domain.UrlMapping
}

func NewGatewayController(urlMapping *domain.UrlMapping) *GatewayController {
	v := new(GatewayController)
	v.urlMapping = urlMapping
	return v
}

func (c *GatewayController) Handle(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	fmt.Println(r.RequestURI)

	targetServiceUrl, err := getTargetServiceUrl(r, c.urlMapping)
	if err != nil {
		http.Error(w, "not found api", http.StatusNotFound)
		return
	}

	logRequestPayload(targetServiceUrl)
	c.reverseProxy(w, r, targetServiceUrl)
}

func (c *GatewayController) reverseProxy(w http.ResponseWriter, r *http.Request, uri string) {
	rr, err := http.NewRequest(r.Method, uri, r.Body)
	responseIfError(w, err)

	copyHeader(r.Header, &rr.Header)

	var transport http.Transport
	resp, err := transport.RoundTrip(rr)
	responseIfError(w, err)

	fmt.Printf("Resp-Headers: %v\n", resp.Header)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	responseIfError(w, err)

	dH := w.Header()
	copyHeader(resp.Header, &dH)
	dH.Add("Requested-Host", rr.Host)

	w.Write(body)
}

func getTargetServiceUrl(request *http.Request, mapping *domain.UrlMapping) (string, error) {
	split := strings.Split(request.URL.Path, "/")
	if len(split) < 4 {
		return "", errors.New("invalid request to api")
	}
	host := mapping.Get(split[3])
	targetService := host + request.URL.Path
	if len(request.URL.RawQuery) != 0 {
		targetService += "?" + request.URL.RawQuery
	}
	return targetService, nil
}

func responseIfError(w http.ResponseWriter, err error) {
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func copyHeader(source http.Header, dest *http.Header) {
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}

func logRequestPayload(proxyUrl string) {
	log.Printf("proxy_url: %s\n", proxyUrl)
}
