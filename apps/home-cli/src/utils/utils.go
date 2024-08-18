package utils

import (
	"os"
	"path"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

type EnvVariable struct {
	Variable     string
	DefaultValue string
}

var (
	ComposeDir        = EnvVariable{"COMPOSE_DIR", "/workspaces/home-config/compose"}
	EnvDir            = EnvVariable{"ENV_DIR", "/workspaces/home-config/compose"}
	OliveTinConfigDir = EnvVariable{"OLIVETIN_CONFIG_DIR", "/workspaces/home-config/olivetin"}
)

var (
	StackList  = getStacks()
	HostList   = getHosts()
	ActionList = getActions()
)

func GetEnvVariable(env EnvVariable) string {
	envDir := os.Getenv(env.Variable)

	if envDir == "" {
		envDir = env.DefaultValue
	}

	return envDir
}

func GetFileInDir(dir EnvVariable, file string) string {
	return path.Join(GetEnvVariable(dir), file)
}

func getStacks() []string {
	entries, err := os.ReadDir(GetEnvVariable(ComposeDir))
	if err != nil {
		return []string{}
	}

	stacks := []string{}

	for _, entry := range entries {
		if entry.IsDir() {
			stacks = append(stacks, entry.Name())
		}
	}

	slices.Sort(stacks)

	return stacks
}

func getHosts() []string {
	hosts := []string{}

	for _, stack := range getStacks() {
		entries, err := os.ReadDir(path.Join(GetEnvVariable(ComposeDir), stack))
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), "."+stack+".yml") {
				parts := strings.Split(entry.Name(), ".")

				hosts = append(hosts, parts[0])
			}
		}
	}

	slices.Sort(hosts)

	return slices.Compact(hosts)
}

func stackFileExists(stack string, host string) bool {
	_, err := os.Stat(path.Join(GetEnvVariable(ComposeDir), stack, host+"."+stack+".yml"))
	return err == nil
}

func GetHostsByStack(stack string) []string {
	hosts := []string{}

	for _, host := range HostList {
		if stackFileExists(stack, host) {
			hosts = append(hosts, host)
		}
	}

	return hosts
}

func GetStacksByHost(host string) []string {
	stacks := []string{}
	for _, stack := range StackList {
		if stackFileExists(stack, host) {
			stacks = append(stacks, stack)
		}
	}

	return stacks
}

func getActions() []string {
	return []string{"up", "down", "restart"}
}

func GetSubCommandsPaths(commands []*cobra.Command) []string {
	paths := []string{}

	for _, command := range commands {
		if slices.Contains(ActionList, command.Use) {
			paths = append(paths, command.CommandPath())
		}

		paths = append(paths, GetSubCommandsPaths(command.Commands())...)
	}

	return paths
}

func FindAll[T any](slice []T, predicate func(T) bool) []T {
	var result []T

	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}

	return result
}
