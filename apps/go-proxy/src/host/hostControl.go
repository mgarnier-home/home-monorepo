package host

import (
	"context"
	"fmt"
	"mgarnier11/go-proxy/config"
	"net"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-ping/ping"
	"golang.org/x/crypto/ssh"
)

func sendSSHCommand(ctx context.Context, config *config.HostConfig, command string) error {
	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(config.Ip, "22"), &ssh.ClientConfig{
		User:            config.SSHUsername,
		Auth:            []ssh.AuthMethod{ssh.Password(config.SSHPassword)},
		Timeout:         2 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to ssh: %v", err)
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()

	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	go session.CombinedOutput(command)

	<-ctx.Done()

	return nil
}

func getHostStatus(ip string) (bool, error) {
	log.Infof("Checking host status: %s", ip)

	pinger, err := ping.NewPinger(ip)

	if err != nil {
		return false, fmt.Errorf("failed to create pinger: %v", err)
	}

	pinger.Count = 1
	pinger.Timeout = 1 * time.Second

	err = pinger.Run()
	if err != nil {
		return false, fmt.Errorf("failed to run pinger: %v", err)
	}

	return pinger.Statistics().PacketsRecv > 0, nil
}
