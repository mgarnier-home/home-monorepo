package bo

type ServerBo struct {
	Id      string
	Name    string
	Version string
	Map     string
	Url     string
	Memory  string
	Port    uint16
	NewMap  bool
}
