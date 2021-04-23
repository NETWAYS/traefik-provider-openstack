package discovery

import (
	osServers "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"path/filepath"
)

type FilterOptions struct {
	Includes []string
	Excludes []string
}

func (f FilterOptions) Matches(name string) bool {
	for _, pattern := range f.Excludes {
		if ok, _ := filepath.Match(pattern, name); ok {
			return false
		}
	}

	if len(f.Includes) == 0 {
		return true
	}

	for _, pattern := range f.Includes {
		if ok, _ := filepath.Match(pattern, name); ok {
			return true
		}
	}

	return false
}

func (f FilterOptions) MatchesServer(server *osServers.Server) bool {
	return f.Matches(server.Name)
}
