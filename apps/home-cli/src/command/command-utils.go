package command

import (
	"bufio"
	"fmt"
	"io"
	"mgarnier11/home-cli/config"
	"mgarnier11/home-cli/utils"

	globalUtils "mgarnier11.fr/go/libs/utils"

	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/fatih/color"
)

type CliCommand struct {
	action string
	stack  string
	host   string
}

type OsCommand struct {
	cliCommand    *CliCommand
	osCommand     string
	osCommandArgs []string
	dir           string
}

func ExecCommand(stacks []string, hosts []string, action string, args []string) {
	compose := utils.GetCompose()

	commands := []*CliCommand{}

	fmt.Println(stacks)
	fmt.Println(hosts)

	for _, stack := range stacks {
		for _, host := range compose.Stacks[stack] {
			if slices.Contains(hosts, host) {
				commands = append(commands, &CliCommand{action, stack, host})
			}
		}
	}

	results := make(map[*CliCommand]error)

	// if slices.Contains(args, "parallel") {
	// 	var wg sync.WaitGroup

	// 	for _, command := range commands {
	// 		wg.Add(1)

	// 		go (func(command *CliCommand, wg *sync.WaitGroup) {
	// 			defer wg.Done()

	// 			results[command] = execCliCommand(command)
	// 		})(command, &wg)
	// 	}

	// 	wg.Wait()
	// } else {
	for _, command := range commands {
		results[command] = execCommand(command)
	}
	// }

	for command, err := range results {
		if err != nil {
			color.Red(fmt.Sprintf("%s %s Error executing command %s", command.host, command.stack, err))
		} else {
			color.Green(fmt.Sprintf("%s %s Successfully executed command", command.host, command.stack))
		}
	}

	err := execOsCommand(&OsCommand{
		cliCommand:    nil,
		osCommand:     "docker",
		osCommandArgs: []string{"context", "use", "default"},
	})

	if err != nil {
		color.Red(fmt.Sprintf("Error resetting context %s", err))
	} else {
		color.Green("Successfully reset context")
	}
}

func getEnvFilesPaths(composeDir string, stack string) []string {
	dir := ""
	if stack != "" {
		dir = fmt.Sprintf("%s/%s", composeDir, stack)
	} else {
		dir = composeDir
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	envFiles := []string{}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".env") {
			if stack != "" {
				envFiles = append(envFiles, fmt.Sprintf("%s/%s", stack, file.Name()))
			} else {
				envFiles = append(envFiles, file.Name())
			}
		}
	}

	return envFiles
}

func getEnvFiles(command *CliCommand) []string {
	composeDir := config.Env.ComposeDir
	globalEnvFiles := getEnvFilesPaths(composeDir, "")
	stackEnvFiles := getEnvFilesPaths(composeDir, command.stack)

	return append(globalEnvFiles, stackEnvFiles...)
}

func getComposeCommandArgs(command *CliCommand) []string {
	args := []string{
		"compose",
	}

	envFiles := getEnvFiles(command)

	for _, envFile := range envFiles {
		args = append(args, "--env-file", envFile)
	}

	args = append(args,
		"-f",
		fmt.Sprintf("%s/%s.%s.yml", command.stack, command.host, command.stack),
	)

	if command.action == "up" {
		args = append(args, "up", "-d", "--pull", "always")
	} else if command.action == "down" {
		args = append(args, "down", "-v")
	} else if command.action == "restart" {
		args = append(args, "up", "-d", "--pull", "always", "--force-recreate")
	}

	return args
}

func setContext(host string) error {
	dockerContextCreateCommand := &OsCommand{
		cliCommand:    nil,
		osCommand:     "docker",
		osCommandArgs: []string{"context", "create", host, "--docker", "host=" + globalUtils.GetEnv(strings.ToUpper(host)+"_HOST", "")},
	}

	err := execOsCommand(dockerContextCreateCommand)

	if err != nil {
		color.Cyan(fmt.Sprintf("%s", err))
	}

	dockerContextUseCommand := &OsCommand{
		cliCommand:    nil,
		osCommand:     "docker",
		osCommandArgs: []string{"context", "use", host},
	}

	err = execOsCommand(dockerContextUseCommand)

	if err != nil {
		return err
	}

	return nil
}

func execCommand(command *CliCommand) error {

	err := setContext(command.host)

	if err != nil {
		return err
	}

	dockerComposeCommand := &OsCommand{
		cliCommand:    command,
		osCommand:     "docker",
		osCommandArgs: getComposeCommandArgs(command),
		dir:           config.Env.ComposeDir,
	}

	err = execOsCommand(dockerComposeCommand)

	return err
}

func execOsCommand(osCommand *OsCommand) error {
	color.Blue(fmt.Sprintf("Executing command %s %s", osCommand.osCommand, osCommand.osCommandArgs))

	cmd := exec.Command(osCommand.osCommand, osCommand.osCommandArgs...)
	cmd.Dir = osCommand.dir

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	go printExecCommandInfo(osCommand.cliCommand, stdout)
	go printExecCommandInfo(osCommand.cliCommand, stderr)

	err := cmd.Run()

	return err
}

func printExecCommandInfo(command *CliCommand, std io.ReadCloser) {
	scanner := bufio.NewScanner(std)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()

		if command == nil {
			color.Yellow(text)
		} else {
			color.Yellow(fmt.Sprintf("%s %s %s", command.host, command.stack, text))
		}
	}
}
