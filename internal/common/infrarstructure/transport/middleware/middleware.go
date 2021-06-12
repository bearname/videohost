package middleware

import (
	"encoding/json"
	"github.com/bearname/videohost/internal/common/dto"
	"github.com/bearname/videohost/internal/common/util"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func LogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})
}

func AllowCors(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*request).Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		next(writer, request)
	}
}

func AuthMiddleware(next http.HandlerFunc, authServerUrl string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if (*request).Method == "OPTIONS" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		authorization := request.Header.Get("Authorization")
		if len(authorization) == 0 {
			http.Error(writer, "authorization header not set", http.StatusUnauthorized)
			return
		}
		body, err := util.GetRequest(&http.Client{}, authServerUrl+"/api/v1/auth/token/validate", authorization)
		if err != nil {
			http.Error(writer, "valid check writes", http.StatusInternalServerError)
			return
		}

		var userDto dto.UserDto
		err = json.Unmarshal(body, &userDto)
		if err != nil {
			log.Error(err.Error())
			http.Error(writer, "valid check build user dto", http.StatusInternalServerError)
			return
		}

		context.Set(request, "userId", userDto.UserId)
		context.Set(request, "username", userDto.Username)

		next(writer, request)
	}
}
