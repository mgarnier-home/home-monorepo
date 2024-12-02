package commandconn

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// New returns net.Conn, establishing an SSH connection using the provided configuration.
func NewSSH(address string, config *ssh.ClientConfig) (net.Conn, error) {
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH: %w", err)
	}
	log.Printf("Connected to SSH server %s", address)

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}
	log.Printf("Created SSH session")

	stdin, err := session.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	err = session.Start("docker system dial-stdio")
	if err != nil {
		return nil, fmt.Errorf("failed to start shell: %w", err)
	}

	log.Printf("Started shell")

	return &sshCommandConn{
		client:     client,
		session:    session,
		stdin:      stdin,
		stdout:     stdout,
		localAddr:  dummyAddr{network: "ssh", s: "local"},
		remoteAddr: dummyAddr{network: "ssh", s: address},
	}, nil
}

type sshCommandConn struct {
	client     *ssh.Client
	session    *ssh.Session
	stdin      io.WriteCloser
	stdout     io.Reader
	closeOnce  sync.Once
	localAddr  net.Addr
	remoteAddr net.Addr
}

func (c *sshCommandConn) Read(p []byte) (int, error) {
	return c.stdout.Read(p)
}

func (c *sshCommandConn) Write(p []byte) (int, error) {
	return c.stdin.Write(p)
}

func (c *sshCommandConn) Close() error {
	var err error
	c.closeOnce.Do(func() {
		err = c.session.Close()
		if cerr := c.client.Close(); cerr != nil {
			err = fmt.Errorf("failed to close client: %w", cerr)
		}
	})
	return err
}

func (c *sshCommandConn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *sshCommandConn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *sshCommandConn) SetDeadline(t time.Time) error {
	// Not implemented for SSH connections
	return nil
}

func (c *sshCommandConn) SetReadDeadline(t time.Time) error {
	// Not implemented for SSH connections
	return nil
}

func (c *sshCommandConn) SetWriteDeadline(t time.Time) error {
	// Not implemented for SSH connections
	return nil
}
