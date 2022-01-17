package discovery

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type ServerAddresses map[string][]ServerAddress

type ServerAddress struct {
	Version   int    `mapstructure:"version"`
	Address   string `mapstructure:"addr"`
	Type      string `mapstructure:"OS-EXT-IPS:type"`
	HWAddress string `mapstructure:"OS-EXT-IPS-MAC:mac_addr"`
}

func GetServerAddresses(input map[string]interface{}) (addresses ServerAddresses, err error) {
	addresses = ServerAddresses{}

	for poolName, pool := range input {
		poolAddresses, ok := pool.([]interface{})
		if !ok {
			err = fmt.Errorf("could not decode pool address to []interface{}%s", pool)
			return
		}

		for _, addressData := range poolAddresses {
			var address ServerAddress

			err = mapstructure.Decode(addressData, &address)
			if err != nil {
				return
			}

			addresses[poolName] = append(addresses[poolName], address)
		}
	}

	return
}
