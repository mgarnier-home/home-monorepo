package common

var ActionList = []string{"up", "down", "restart"}

type Command struct {
	Command     string       `yaml:"command"`
	ComposeFile *ComposeFile `yaml:"compose_file"`
	Action      string       `yaml:"action"`
}

type ComposeFile struct {
	Name  string `yaml:"name"`
	Path  string `yaml:"path"`
	Host  string `yaml:"host"`
	Stack string `yaml:"stack"`
}

type ComposeConfig struct {
	Host       string                 `yaml:"host"`
	Stack      string                 `yaml:"stack"`
	Action     string                 `yaml:"action"`
	Config     string                 `yaml:"config"`
	HostConfig string                 `yaml:"host_config"`
	Services   map[string]interface{} `yaml:"services"`
}

type ComposeFileSource struct {
	Services map[string]interface{} `yaml:"services"`
}
