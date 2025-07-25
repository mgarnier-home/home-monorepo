package sshutils

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

func GetSSHKeyAuth(sshKey string) (ssh.AuthMethod, error) {
	// Parse the private key
	signer, err := ssh.ParsePrivateKey([]byte(sshKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return ssh.PublicKeys(signer), nil
}

func GetSSHClient(sshUsername, sshIP, sshPort, sshKey string) (*ssh.Client, error) {
	sshAuth, err := GetSSHKeyAuth(sshKey)

	if err != nil {
		return nil, fmt.Errorf("error getting ssh key auth: %v", err)
	}

	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(sshIP, sshPort), &ssh.ClientConfig{
		User:            sshUsername,
		Auth:            []ssh.AuthMethod{sshAuth},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	if err != nil {
		return nil, fmt.Errorf("error connecting to ssh: %v", err)
	}

	return sshClient, nil
}
