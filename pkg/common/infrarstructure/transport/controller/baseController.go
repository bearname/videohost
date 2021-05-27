package controller

import (
	"encoding/json"
	"errors"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport"
	"github.com/bearname/videohost/pkg/common/util"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type BaseController struct {
}

func (c *BaseController) AllowCorsRequest(writer *http.ResponseWriter) {
	(*writer).Header().Set("Access-Control-Allow-Origin", "*")
	(*writer).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*writer).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
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
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *BaseController) WriteResponse(w *http.ResponseWriter, statusCode int, success bool, message string) {
	(*w).WriteHeader(statusCode)
	response := transport.Response{
		Success: success,
		Message: message,
	}
	c.JsonResponse(*w, response)
}

func (c *BaseController) GetIntRouteParameter(request *http.Request, key string) (int, error) {
	pageStr, done := c.ParseRouteParameter(request, key)
	if !done {
		return 0, errors.New("Invalid " + key + " parameter not found")
	}

	page, b := c.validate(pageStr)
	if b {
		return 0, errors.New("Invalid " + key + " parameter not found")
	}

	return page, nil
}

func (c *BaseController) ParseRouteParameter(request *http.Request, key string) (string, bool) {
	query := request.URL.Query()
	keys, ok := query[key]

	if !ok || len(keys) != 1 {
		return "", false
	}

	return keys[0], true
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

func (c *BaseController) validate(pageStr string) (int, bool) {
	page, b := util.StrToInt(pageStr)
	if !b || page < 0 {
		return 0, true
	}
	return page, false
}
