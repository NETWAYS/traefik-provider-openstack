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
		// VM's with the traefik METATAG will be ignored
		for key, value := range s.Metadata {
			labels[key] = value
		}
	} else if settings.DefaultEnable {
		// Build labels for defaults while evaluating templates
		for key, value := range settings.DefaultLabels {
			key, err = s.EvalTemplate(key)
			if err != nil {
				err = fmt.Errorf("could not execute template for label key '%s': %w", key, err)
				return
			}

			value, err = s.EvalTemplate(value)
			if err != nil {
				err = fmt.Errorf("could not execute template for label value '%s': %w", value, err)
				return
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
			for name := range config.HTTP.Services {
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

	if len(config.TCP.Services) > 0 && len(config.TCP.Routers) == 0 {
		if config.TCP.Routers == nil {
			config.TCP.Routers = map[string]*dynamic.TCPRouter{}
		}

		config.TCP.Routers[s.Name] = &dynamic.TCPRouter{}
	}

	// Update all routers with unset values
	for _, router := range config.TCP.Routers {
		// set entrypoint when not set
		if len(router.EntryPoints) == 0 {
			router.EntryPoints = settings.DefaultEntrypoints
		}

		// set default rule when not set
		if router.Rule == "" {
			rule, erro := s.EvalTemplate(settings.DefaultRule)
			if erro != nil {
				return erro
			}

			router.Rule = rule
		}

		// Set service when not set
		if router.Service == "" {
			for name := range config.TCP.Services {
				// Just use the first and break
				router.Service = name

				break
			}
		}
	}

	// Check and update services with defaults
	for _, service := range config.TCP.Services {
		if service.LoadBalancer != nil {
			for i := range service.LoadBalancer.Servers {
				lbServer := &service.LoadBalancer.Servers[i]

				// build URL when port is set
				if lbServer.Address == "" && lbServer.Port != "" {
					url := fmt.Sprintf("%s:%s", address, lbServer.Port)
					lbServer.Address = url
				}

			}
		}
	}

	configurations[s.Name] = config

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
