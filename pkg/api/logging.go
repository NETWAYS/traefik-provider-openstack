package api

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

var Logger = log.StandardLogger()

func HTTPLogging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			HTTPLogRequest(w, r)
		})
	}
}

func HTTPLogRequest(_ http.ResponseWriter, r *http.Request) {
	fields := HTTPLogFields(r)
	Logger.WithFields(fields).Info("HTTP API request")
}

func HTTPLogError(r *http.Request, err error) {
	fields := HTTPLogFields(r)
	fields["err"] = err

	Logger.WithFields(fields).Error("returning API error to client")
}

func HTTPLogFields(r *http.Request) (fields log.Fields) {
	fields = log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"remote": r.RemoteAddr,
	}

	return
}
