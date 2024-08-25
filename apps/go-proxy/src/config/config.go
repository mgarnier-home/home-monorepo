package config

type ProxyConfig struct {
	ListenPort int
	TargetPort int
	Protocol   string
	Name       string
}

type HostConfig struct {
	Proxies      []*ProxyConfig
	Name         string
	Ip           string
	MacAddress   string
	SSHUername   string
	SSHPassword  string
	Autostop     bool
	MaxAliveTime int
	DockerPort   int
}
