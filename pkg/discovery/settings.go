package discovery

type Settings struct {
	DefaultRule        string
	DefaultEnable      bool
	DefaultEntrypoints []string
	DefaultLabels      map[string]string
	AddressType        string
}

var DefaultSettings = Settings{
	DefaultRule:   "Host(`{{ .Name }}`)",
	DefaultLabels: DefaultLabels,
}

var DefaultLabels = map[string]string{
	"traefik.http.services.{{ .Name }}.loadBalancer.server.port": "80",
}
