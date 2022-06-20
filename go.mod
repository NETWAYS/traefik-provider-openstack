module github.com/NETWAYS/traefik-provider-openstack

go 1.16

require (
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/gophercloud/gophercloud v0.25.0
	github.com/gorilla/mux v1.8.0
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cast v1.5.0 // indirect
	github.com/stretchr/testify v1.7.1
	github.com/traefik/traefik/v2 v2.7.1
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c // indirect
	k8s.io/apimachinery v0.24.2 // indirect
	sigs.k8s.io/json v0.0.0-20220525155127-227cbc7cc124 // indirect
)

// Containous forks for traefik - we can't compile without them
// See https://github.com/traefik/traefik/blob/master/go.mod
replace (
	github.com/abbot/go-http-auth => github.com/containous/go-http-auth v0.4.1-0.20200324110947-a37a7636d23e
	github.com/go-check/check => github.com/containous/check v0.0.0-20170915194414-ca0bf163426a
//github.com/gorilla/mux => github.com/containous/mux v0.0.0-20181024131434-c33f32e26898
//github.com/mailgun/minheap => github.com/containous/minheap v0.0.0-20190809180810-6e71eb837595
//github.com/mailgun/multibuf => github.com/containous/multibuf v0.0.0-20190809014333-8b6c9a7e6bba
)
