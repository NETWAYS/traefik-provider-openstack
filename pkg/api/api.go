package api

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	goLog "log"
	"net/http"
	"time"
)

type Server struct {
	Server *http.Server
}

const DefaultAddress = ":8080"

func NewServer() (s *Server, err error) {
	s = &Server{}

	router := NewRouter()
	router.Use(HTTPLogging())

	s.Server = &http.Server{
		Addr:         DefaultAddress,
		Handler:      router,
		ErrorLog:     goLog.New(Logger.Writer(), "http: ", goLog.LstdFlags),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return
}

func (s *Server) ListenAndServe() error {
	err := s.Server.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		err = fmt.Errorf("could not start HTTP server: %w", err)
		log.WithError(err).Error("could not start HTTP server")
	}

	return err
}

func HandleHome(w http.ResponseWriter, _ *http.Request) {
	RespondWithJSON(w, MessageResponse{"Welcome, see the /traefik endpoint"}, http.StatusOK)
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(NotFound)
	r.HandleFunc("/", HandleHome)

	r.HandleFunc("/traefik", HandleTraefik)

	return r
}
