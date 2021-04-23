package discovery

import (
	"github.com/gophercloud/gophercloud"
	osServers "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
)

type ServerList map[string]*Server

type Servers struct {
	Servers ServerList
}

func LoadServers(client *gophercloud.ServiceClient, filter FilterOptions) (servers *Servers, err error) {
	servers = &Servers{}
	servers.Servers = ServerList{}

	err = osServers.List(client, nil).EachPage(func(page pagination.Page) (bool, error) {
		s, err := osServers.ExtractServers(page)
		if err != nil {
			return false, err
		}

		// Inject all servers into our list, that match our filter
		for _, server := range s {
			if filter.MatchesServer(&server) {
				servers.Servers[server.ID] = &Server{server}
			}
		}

		return true, nil
	})

	return
}

func (s Servers) RegisterConfiguration(configurations map[string]*dynamic.Configuration, settings Settings) (err error) {
	for _, server := range s.Servers {
		err = server.RegisterConfiguration(configurations, settings)
		if err != nil {
			return
		}
	}

	return
}
