package discovery

import (
	"bytes"
	"fmt"
	osServers "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/config/label"
	"text/template"
)

type Server struct {
	osServers.Server
}

func (s Server) GetAddress(kind string) (address string, err error) {
	var first string

	// Decode addresses
	addresses, err := GetServerAddresses(s.Addresses)
	if err != nil {
		return
	}

	// Find the first address of requested
	for _, poolAddresses := range addresses {
		for _, a := range poolAddresses {
			if a.Type == kind {
				address = a.Address
				return
			}
			if first == "" {
				first = a.Address
			}
		}
	}

	// Return first found of any type as fallback
	return first, nil
}

func (s Server) RegisterConfiguration(configurations map[string]*dynamic.Configuration, settings Settings) (err error) {
	labels := make(map[string]string)

	// Is traffic enabled via meta data flag?
	if enable, ok := s.Metadata[TraefikEnable]; ok && enable == "true" {
		// Use labels from OpenStack
		labels = make(map[string]string)

		// TODO: filter by traefik prefix?
		for key, value := range s.Metadata {
			labels[key] = value
		}
	} else if settings.DefaultEnable {
		// Build labels for defaults while evaluating templates
		for key, value := range settings.DefaultLabels {
			key, err = s.EvalTemplate(key)
			if err != nil {
				err = fmt.Errorf("could not execute template for label key '%s': %w", key, err)
			}

			value, err = s.EvalTemplate(value)
			if err != nil {
				err = fmt.Errorf("could not execute template for label value '%s': %w", value, err)
			}

			labels[key] = value
		}
	} else {
		// Ignore this host
		return
	}

	// Parse server main address from list
	address, err := s.GetAddress(settings.AddressType)
	if err != nil {
		return
	}

	// Decode from labels
	config, err := label.DecodeConfiguration(labels)
	if err != nil {
		err = fmt.Errorf("could not decode labels: %w", err)
		return
	}

	// Create a router when a service is present
	if len(config.HTTP.Services) > 0 && len(config.HTTP.Routers) == 0 {
		if config.HTTP.Routers == nil {
			config.HTTP.Routers = map[string]*dynamic.Router{}
		}

		config.HTTP.Routers[s.Name] = &dynamic.Router{}
	}

	// Update all routers with unset values
	for _, router := range config.HTTP.Routers {
		// set entrypoint when not set
		if len(router.EntryPoints) == 0 {
			router.EntryPoints = settings.DefaultEntrypoints
		}

		// set default rule when not set
		if router.Rule == "" {
			rule, err := s.EvalTemplate(settings.DefaultRule)
			if err != nil {
				return err
			}

			router.Rule = rule
		}

		// Set service when not set
		if router.Service == "" {
			for name, _ := range config.HTTP.Services {
				// Just use the first and break
				router.Service = name
				break
			}
		}
	}

	// Check and update services with defaults
	for _, service := range config.HTTP.Services {
		if service.LoadBalancer != nil {
			for i := range service.LoadBalancer.Servers {
				lbServer := &service.LoadBalancer.Servers[i]

				// build URL when port is set
				if lbServer.URL == "" && lbServer.Port != "" {
					if lbServer.Scheme == "" {
						lbServer.Scheme = "http"
					}

					url := fmt.Sprintf("%s://%s:%s", lbServer.Scheme, address, lbServer.Port)
					lbServer.URL = url
				}
			}
		}
	}

	configurations[s.Name] = config

	/*
		name := s.Name
		host := s.Name
		hostCockpit := s.Name + "-cockpit"

		// TODO: filter Status
		// TODO: filter Tags?

		if settings.Domain != "" {
			host += "." + settings.Domain
			hostCockpit += "." + settings.Domain
		}

		router := "traefik.http.routers." + name
		service := "traefik.http.services." + name

		address, err := s.GetAddress(settings.AddressType)
		if err != nil {
			return
		}

		labels := map[string]string{
			router + ".entrypoints":                      "http", // TODO
			router + ".rule":                             fmt.Sprintf("Host(`%s`)", host),
			router + ".service":                          name,
			service + ".loadBalancer.server.url":         fmt.Sprintf("http://%s/", address),
			router + "-cockpit.entrypoints":              "http", // TODO
			router + "-cockpit.rule":                     fmt.Sprintf("Host(`%s`)", hostCockpit),
			router + "-cockpit.service":                  name + "-cockpit",
			service + "-cockpit.loadBalancer.server.url": fmt.Sprintf("https://%s:9090/", address),
		}

		config, err := label.DecodeConfiguration(labels)
		if err != nil {
			err = fmt.Errorf("could not decode labels: %w", err)
			return
		}

		configurations[name] = config
	*/

	return
}

func (s Server) EvalTemplate(t string) (string, error) {
	tmpl, err := template.New("temp").Parse(t)
	if err != nil {
		return t, err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, s)
	if err != nil {
		return t, err
	}

	return buf.String(), nil
}

func (s Server) BuildStandardLabels(settings Settings) map[string]string {
	return map[string]string{
		//"traefik.http.routers."+s.Name+".entrypoints":                       "{ENTRYPOINT}",
		//"traefik.http.routers."+s.Name+"-cockpit.entrypoints":               "{ENTRYPOINT}",
		//"traefik.http.routers." + s.Name + ".service":                           s.Name,
		//"traefik.http.routers." + s.Name + "-cockpit.service":                   s.Name + "-cockpit",
		//"traefik.http.routers." + s.Name + ".rule":                              settings.DefaultRule,
		//"traefik.http.services." + s.Name + ".loadBalancer.server.port":         "80",
		//"traefik.http.routers." + s.Name + ".rule":                      settings.DefaultRule,
		"traefik.http.services." + s.Name + ".loadBalancer.server.port": "9090",
	}
}
