package api

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	Error(w, r, errors.New("404 Not Found"), http.StatusNotFound)
}

func Error(w http.ResponseWriter, r *http.Request, err error, code int) {
	HTTPLogError(r, err)
	RespondWithJSON(w, ErrorResponse{err.Error()}, code)
}

func RespondWithJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		log.WithError(err).Error("could not write result to client")
	}
}
