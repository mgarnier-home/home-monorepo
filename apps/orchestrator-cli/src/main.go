package main

import (
	"os/exec"

	"mgarnier11.fr/go/orchestrator-cli/cmd"
	"mgarnier11.fr/go/orchestrator-cli/config"
)

func main() {

	if config.Env.SshKeyPath != "" {
		err := exec.Command("chmod", "600", config.Env.SshKeyPath).Run()

		if err != nil {
			panic(err)
		}
	}

	cmd.Execute()
}
