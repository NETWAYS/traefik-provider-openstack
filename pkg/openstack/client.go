package openstack

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

func NewAuthProviderClient() (provider *gophercloud.ProviderClient, err error) {
	ao, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		err = fmt.Errorf("could not build auth options: %w", err)
		return
	}

	provider, err = openstack.AuthenticatedClient(ao)
	if err != nil {
		err = fmt.Errorf("could not build provider client: %w", err)
		return
	}

	return
}

func NewComputeClient() (client *gophercloud.ServiceClient, err error) {
	provider, err := NewAuthProviderClient()
	if err != nil {
		return
	}

	client, err = openstack.NewComputeV2(provider, gophercloud.EndpointOpts{})
	if err != nil {
		err = fmt.Errorf("could not build compute client: %w", err)
		return
	}

	return
}
