package handler

import (
	"fmt"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "{\"status\": \"OK\"}")
}

func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "{\"host\": \"%v\"}", r.Host)
}
