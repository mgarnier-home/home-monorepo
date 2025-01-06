package sshutils

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func GetSSHKeyAuth(sshKeyPath string) (ssh.AuthMethod, error) {
	// Load the private key file
	key, err := os.ReadFile(sshKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	// Parse the private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return ssh.PublicKeys(signer), nil
}
