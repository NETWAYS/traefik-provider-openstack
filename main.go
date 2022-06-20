package main

import (
	"context"
	"github.com/NETWAYS/traefik-provider-openstack/pkg/api"
	"github.com/NETWAYS/traefik-provider-openstack/pkg/discovery"
	"github.com/NETWAYS/traefik-provider-openstack/pkg/openstack"
	log "github.com/sirupsen/logrus"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/provider"
	"os"
	"os/signal"
	"time"
)

func main() {
	err := UpdateConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP server
	err = SetupAndStartServer()
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateConfiguration() (err error) {
	client, err := openstack.NewComputeClient()
	if err != nil {
		return
	}

	settings := discovery.DefaultSettings
	settings.DefaultRule = "Host(`{{ .Name }}.training.netways.de`)"
	settings.DefaultEnable = true
	settings.AddressType = "floating"
	settings.DefaultLabels = map[string]string{
		// HTTP
		"traefik.http.routers.{{ .Name }}-http.service":                   "{{ .Name }}-http",
		"traefik.http.routers.{{ .Name }}-http.rule":                      "Host(`{{ .Name }}.training.netways.de`)",
		"traefik.http.routers.{{ .Name }}-http.entrypoints":               "http",
		"traefik.http.services.{{ .Name }}-http.loadBalancer.server.port": "80",
		"traefik.http.routers.{{ .Name }}-http.middlewares":               "http-to-https",

		// HTTPS
		"traefik.http.routers.{{ .Name }}-https.service":                   "{{ .Name }}-https",
		"traefik.http.routers.{{ .Name }}-https.rule":                      "Host(`{{ .Name }}.training.netways.de`)",
		"traefik.http.routers.{{ .Name }}-https.entrypoints":               "https",
		"traefik.http.services.{{ .Name }}-https.loadBalancer.server.port": "80",
		"traefik.http.routers.{{ .Name }}-https.tls":                       "true",

		// TODO The PathPrefix-/Path-Rule does not work, because it searches for the actual Filepath inside '/var/www/...' of the webiste, which does not exist.
		// COCKPIT
		//"traefik.http.routers.{{ .Name }}-cockpit.service":                     "{{ .Name }}-cockpit",
		//"traefik.http.routers.{{ .Name }}-cockpit.rule":                        "Host(`{{ .Name }}.training.netways.de`) && PathPrefix(`/cockpit`)",
		//"traefik.http.routers.{{ .Name }}-cockpit.entrypoints":                 "https",
		//"traefik.http.services.{{ .Name }}-cockpit.loadBalancer.server.port":   "9090",
		//"traefik.http.services.{{ .Name }}-cockpit.loadBalancer.server.scheme": "https",
		//"traefik.http.routers.{{ .Name }}-cockpit.tls":                         "true",

		// SSH
		"traefik.tcp.routers.{{ .Name }}-ssh.rule":                      "HostSNI(`*`)",
		"traefik.tcp.routers.{{ .Name }}-ssh.entrypoints":               "ssh",
		"traefik.tcp.routers.{{ .Name }}-ssh.service":                   "{{ .Name }}-ssh",
		"traefik.tcp.services.{{ .Name }}-ssh.loadBalancer.server.port": "22",

		// Redirect HTTP to HTTPS
		"traefik.http.middlewares.http-to-https.redirectscheme.scheme":    "https",
		"traefik.http.middlewares.http-to-https.redirectscheme.permanent": "true",
	}

	filter := discovery.FilterOptions{}

	server, err := discovery.LoadServers(client, filter)
	if err != nil {
		return
	}

	configurations := make(map[string]*dynamic.Configuration)

	err = server.RegisterConfiguration(configurations, settings)
	if err != nil {
		return
	}

	config := provider.Merge(context.Background(), configurations)

	if len(config.TCP.Routers) == 0 && len(config.TCP.Services) == 0 {
		config.TCP = nil
	}
	if len(config.UDP.Routers) == 0 && len(config.UDP.Services) == 0 {
		config.UDP = nil
	}

	api.CurrentTraefikConfiguration = config

	return
}

func SetupAndStartServer() (err error) {
	s, err := api.NewServer()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Info("shutting down HTTP API")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		s.Server.SetKeepAlivesEnabled(false)
		if err := s.Server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v", err)
		}
	}()

	log.Info("starting HTTP API server")

	err = s.ListenAndServe()
	if err != nil {
		return
	}

	log.Info("stopped HTTP API server")

	return
}
