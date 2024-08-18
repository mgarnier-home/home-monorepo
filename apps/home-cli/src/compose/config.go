package compose

import (
	"mgarnier11/home-cli/utils"
)

type Config struct {
	Stacks  map[string][]string
	Hosts   map[string][]string
	Actions []string
}

func GetConfig() *Config {
	config := &Config{}

	config.Stacks = make(map[string][]string)
	config.Hosts = make(map[string][]string)
	config.Actions = utils.ActionList

	for _, stack := range utils.StackList {
		config.Stacks[stack] = utils.GetHostsByStack(stack)
	}

	for _, host := range utils.HostList {
		config.Hosts[host] = utils.GetStacksByHost(host)
	}

	return config
}
