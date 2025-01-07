package bo

type HostBo struct {
	Name         string
	Ip           string
	ProxyIp      string
	SSHUsername  string
	SSHPort      string
	StartPort    int
	MineagerPath string
	Ping         bool
}
