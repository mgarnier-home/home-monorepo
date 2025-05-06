package main

import (
	"mgarnier11/home-cli/cmd"
	"os/exec"
)

func main() {
	err := exec.Command("chmod", "600", "/run/secrets/ssh_private_key").Run()

	if err != nil {
		panic(err)
	}

	cmd.Execute()
}
