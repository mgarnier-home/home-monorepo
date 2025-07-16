package osUtils

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"mgarnier11.fr/go/libs/logger"
)

type OsCommand struct {
	OsCommand     string
	OsCommandArgs []string
	Dir           string
}

func ExecOsCommandOutput(osCommand *OsCommand, commandLog string) (string, error) {
	logger.Debugf("Executing OS command: %s %s in directory: %s", osCommand.OsCommand, strings.Join(osCommand.OsCommandArgs, " "), osCommand.Dir)

	cmd := exec.Command(osCommand.OsCommand, osCommand.OsCommandArgs...)
	cmd.Dir = osCommand.Dir

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("error executing command %s %s: %w", osCommand.OsCommand, strings.Join(osCommand.OsCommandArgs, " "), err)
	}

	return string(output), nil
}

func ExecOsCommand(osCommand *OsCommand, commandLog string) error {
	logger.Debugf("Executing OS command: %s %s in directory: %s", osCommand.OsCommand, strings.Join(osCommand.OsCommandArgs, " "), osCommand.Dir)

	cmd := exec.Command(osCommand.OsCommand, osCommand.OsCommandArgs...)
	cmd.Dir = osCommand.Dir

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	go printExecCommandInfo(commandLog, stdout)
	go printExecCommandInfo(commandLog, stderr)

	err := cmd.Run()

	return err
}

func printExecCommandInfo(command string, std io.ReadCloser) {
	scanner := bufio.NewScanner(std)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()

		if command == "" {
			color.Yellow(text)
		} else {
			color.Yellow(fmt.Sprintf("%s %s", command, text))
		}
	}
}
