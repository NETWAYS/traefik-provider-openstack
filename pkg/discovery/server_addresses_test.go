package discovery

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var jsonAddressData = `{
	"private": [
		{
			"OS-EXT-IPS-MAC:mac_addr": "00:0c:29:0d:11:74",
			"OS-EXT-IPS:type": "fixed",
			"addr": "192.168.1.30",
			"version": 4
		}
	]
}`

func TestGetServerAddresses(t *testing.T) {
	var input map[string]interface{}

	err := json.NewDecoder(strings.NewReader(jsonAddressData)).Decode(&input)
	assert.NoError(t, err)

	/*
	i := map[string]interface{}{
		"private": map[string]interface{}{
			"version": 4,
			"address": "10.123.0.1",
			"OS-EXT-IPS:type": "fixed",
			"OS-EXT-IPS-MAC:mac_addr": "fa:ff:ff:ff:ff:ff",
		},
	}
	*/

	a, err := GetServerAddresses(input)
	assert.NoError(t, err)
	assert.Len(t, a, 1)
	assert.Contains(t, a, "private")
	assert.Len(t, a["private"], 1)
	assert.Equal(t, "192.168.1.30", a["private"][0].Address)
	assert.Equal(t, "fixed", a["private"][0].Type)
	assert.Equal(t, "00:0c:29:0d:11:74", a["private"][0].HWAddress)
}
