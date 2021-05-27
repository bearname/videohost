package controller

import (
	"encoding/json"
	"errors"
	"github.com/bearname/videohost/pkg/common/infrarstructure/transport"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type BaseController struct {
}

func (c *BaseController) ParseMuxVariable(request *http.Request, keys []string) ([]string, error) {
	vars := mux.Vars(request)
	var result []string
	var videoId string
	var ok bool

	for _, key := range keys {
		if videoId, ok = vars[key]; !ok {
			return nil, errors.New(key + " not present")
		}
		result = append(result, videoId)
	}

	return result, nil
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
