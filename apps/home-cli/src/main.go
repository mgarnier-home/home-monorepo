package main

import (
	"mgarnier11/home-cli/cmd"
	"mgarnier11/home-cli/config"
	"os/exec"
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
