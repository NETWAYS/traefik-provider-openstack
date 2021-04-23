package api

import (
	"fmt"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"net/http"
)

var CurrentTraefikConfiguration *dynamic.Configuration

func HandleTraefik(w http.ResponseWriter, r *http.Request) {
	if CurrentTraefikConfiguration == nil {
		Error(w, r, fmt.Errorf("configuration is empty!"), 500)
	}

	RespondWithJSON(w, CurrentTraefikConfiguration, http.StatusOK)
}

func init() {
	Router.HandleFunc("/traefik", HandleTraefik)
}
