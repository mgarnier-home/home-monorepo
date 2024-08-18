package command

import (
	"bufio"
	"fmt"
	"io"
	"mgarnier11/home-cli/compose"
	"mgarnier11/home-cli/utils"
	"os/exec"
	"slices"
	"sync"

	"github.com/fatih/color"
)

type Command struct {
	action string
	stack  string
	host   string
}

type CommandV2 struct {
	Action string
	Stacks []string
	Hosts  []string
}

func ExecCommand(stacks []string, hosts []string, action string, args []string) {
	config := compose.GetConfig()

	commands := []*Command{}

	fmt.Println(stacks)
	fmt.Println(hosts)

	for _, stack := range stacks {
		for _, host := range config.Stacks[stack] {
			if slices.Contains(hosts, host) {
				commands = append(commands, &Command{action, stack, host})
			}
		}
	}

	results := make(map[*Command]error)

	if slices.Contains(args, "parallel") {
		var wg sync.WaitGroup

		for _, command := range commands {
			wg.Add(1)

			go (func(command *Command, wg *sync.WaitGroup) {
				defer wg.Done()

				results[command] = execCliCommand(command)
			})(command, &wg)
		}

		wg.Wait()
	} else {
		for _, command := range commands {
			results[command] = execCliCommand(command)
		}
	}

	for command, err := range results {
		if err != nil {
			color.Red(fmt.Sprintf("%s %s Error executing command %s", command.host, command.stack, err))
		} else {
			color.Green(fmt.Sprintf("%s %s Successfully executed command", command.host, command.stack))
		}
	}

}

func print(command *Command, std io.ReadCloser) {
	scanner := bufio.NewScanner(std)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()

		color.Yellow(fmt.Sprintf("%s %s %s", command.host, command.stack, text))
	}
}

func execCliCommand(command *Command) error {
	var commandArgs = []string{
		"ansible-playbook",
		"playbooks/compose.playbook.yml",
		"--extra-vars",
		"stack=" + command.stack,
		"--extra-vars",
		"command=" + command.action,
		"-l",
		command.host,
	}

	color.Blue(fmt.Sprint(commandArgs))

	cmd := exec.Command(commandArgs[0], commandArgs[1:]...)

	cmd.Dir = utils.GetDir(utils.AnsibleDir)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	go print(command, stdout)
	go print(command, stderr)

	err := cmd.Run()

	return err

}
