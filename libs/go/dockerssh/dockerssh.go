package dockerssh

import (
	"context"
	"fmt"
	"io"
	"mgarnier11/go/sshutils"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"golang.org/x/crypto/ssh"
)

// New returns net.Conn, establishing an SSH connection using the provided configuration.
func NewSSHDialer(address string, config *ssh.ClientConfig) (net.Conn, error) {
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH: %w", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %w", err)
	}

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
		closeErr := c.session.Close()

		if closeErr != nil {
			err = fmt.Errorf("failed to close session: %w", closeErr)
		}

		closeErr = c.client.Close()

		if closeErr != nil {
			err = fmt.Errorf("failed to close client: %w", closeErr)
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

type dummyAddr struct {
	network string
	s       string
}

func (d dummyAddr) Network() string {
	return d.network
}

func (d dummyAddr) String() string {
	return d.s
}

func GetDockerClient(sshUsername string, hostIp string, sshPort string, sshKeyPath string) (*client.Client, error) {
	authMethod, err := sshutils.GetSSHKeyAuth(sshKeyPath)

	if err != nil {
		return nil, fmt.Errorf("failed to get ssh key auth: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User:            sshUsername,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Replace with a proper callback in production
	}

	sshDialer := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return NewSSHDialer(
			net.JoinHostPort(hostIp, sshPort),
			sshConfig,
		)
	}

	httpClient := &http.Client{
		// No tls
		// No proxy
		Transport: &http.Transport{
			DialContext: sshDialer,
		},
		Timeout: 2 * time.Second,
	}

	var clientOpts []client.Opt

	clientOpts = append(clientOpts,
		client.WithHTTPClient(httpClient),
		client.WithDialContext(sshDialer),
	)

	client, err := client.NewClientWithOpts(clientOpts...)

	if err != nil {
		return nil, err
	}

	return client, nil
}
