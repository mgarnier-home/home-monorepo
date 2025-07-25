package osutils

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"mgarnier11.fr/go/libs/logger"
)

var Logger = logger.NewLogger("[OS]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)

type OsCommand struct {
	OsCommand     string
	OsCommandArgs []string
	Dir           string
}

func ExecOsCommandOutput(osCommand *OsCommand, commandLog string) (string, error) {
	Logger.Debugf("Executing OS command: %s %s in directory: %s", osCommand.OsCommand, strings.Join(osCommand.OsCommandArgs, " "), osCommand.Dir)

	cmd := exec.Command(osCommand.OsCommand, osCommand.OsCommandArgs...)
	cmd.Dir = osCommand.Dir

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("error executing command %s %s: %w", osCommand.OsCommand, strings.Join(osCommand.OsCommandArgs, " "), err)
	}

	return string(output), nil
}

func ExecOsCommand(osCommand *OsCommand, commandLog string) error {
	Logger.Debugf("Executing OS command: %s %s in directory: %s", osCommand.OsCommand, strings.Join(osCommand.OsCommandArgs, " "), osCommand.Dir)

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

var mu sync.Mutex

func streamWithPrefix(reader io.Reader, writer io.Writer, prefix string, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		mu.Lock()
		str := color.YellowString("%s - %s", prefix, line)
		Logger.Infof("%s", str)
		fmt.Fprintf(writer, "%s\n", str)
		if f, ok := writer.(interface{ Flush() }); ok {
			f.Flush()
		}
		mu.Unlock()
	}
}

func ExecOsCommandStream(osCommand *OsCommand, writer io.Writer, prefix string) error {
	Logger.Debugf("Executing command in directory: %s", osCommand.Dir)
	Logger.Debugf("%s %s", osCommand.OsCommand, strings.Join(osCommand.OsCommandArgs, " "))

	cmd := exec.Command(osCommand.OsCommand, osCommand.OsCommandArgs...)
	cmd.Dir = osCommand.Dir

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	var wg sync.WaitGroup
	wg.Add(2)
	go streamWithPrefix(stdout, writer, prefix, &wg)
	go streamWithPrefix(stderr, writer, prefix, &wg)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command %s %s: %w", osCommand.OsCommand, strings.Join(osCommand.OsCommandArgs, " "), err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for command %s %s to finish: %w", osCommand.OsCommand, strings.Join(osCommand.OsCommandArgs, " "), err)
	}

	wg.Wait()

	return nil
}
