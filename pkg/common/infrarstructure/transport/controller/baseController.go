package controller

import (
	"encoding/json"
	util2 "github.com/bearname/videohost/pkg/common/util"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type BaseController struct {
}

func (c *BaseController) AllowCorsRequest(writer *http.ResponseWriter) *http.ResponseWriter {
	(*writer).Header().Set("Access-Control-Allow-Origin", "*")
	(*writer).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*writer).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	return writer
}

func (c *BaseController) JsonResponse(writer http.ResponseWriter, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = writer.Write(jsonData)
	if err != nil {
		return
	}
}

func (c *BaseController) GetIntRouteParameter(writer http.ResponseWriter, request *http.Request, key string) (int, bool) {
	pageStr, done := c.ParseRouteParameter(request, key)
	if !done {
		http.Error(writer, "400 "+key+" parameter not found", http.StatusBadRequest)
		return 0, false
	}

	page, b := c.validate(writer, pageStr)
	if b {
		http.Error(writer, "400 invalid "+key+" parameter", http.StatusBadRequest)
		return 0, false
	}
	return page, true
}

func (c *BaseController) ParseRouteParameter(request *http.Request, key string) (string, bool) {
	query := request.URL.Query()
	keys, ok := query[key]

	if !ok || len(keys) != 1 {
		return "", false
	}

	return keys[0], true
}

func (c *BaseController) validate(writer http.ResponseWriter, pageStr string) (int, bool) {
	page, b := util2.StrToInt(pageStr)
	if !b || page < 0 {
		http.Error(writer, "400 Invalid page parameter", http.StatusBadRequest)
		return 0, true
	}
	return page, false
}

func (c *BaseController) WriteResponseData(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if _, err = io.WriteString(w, string(b)); err != nil {
		log.WithField("err", err).Error("write response error")
	}
}
