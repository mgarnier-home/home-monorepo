package host

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"mgarnier11.fr/go/libs/sshutils"

	"mgarnier11.fr/go/go-proxy/config"

	"golang.org/x/crypto/ssh"
)

func sendSSHCommand(ctx context.Context, hostConfig *config.HostConfig, command string) error {
	authMethod, err := sshutils.GetSSHKeyAuth(config.Config.SSHPrivateKey)

	if err != nil {
		return fmt.Errorf("failed to get ssh key auth: %v", err)
	}

	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(hostConfig.Ip, hostConfig.SSHPort), &ssh.ClientConfig{
		User:            hostConfig.SSHUsername,
		Auth:            []ssh.AuthMethod{authMethod},
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

// MagicPacket is a slice of 102 bytes containing the magic packet data.
type MagicPacket [102]byte

// NewMagicPacket allocates a new MagicPacket with the specified MAC.
func newMagicPacket(macAddr string) (packet MagicPacket, err error) {
	mac, err := net.ParseMAC(macAddr)
	if err != nil {
		return packet, err
	}

	if len(mac) != 6 {
		return packet, errors.New("invalid EUI-48 MAC address")
	}

	// write magic bytes to packet
	copy(packet[0:], []byte{255, 255, 255, 255, 255, 255})
	offset := 6

	for i := 0; i < 16; i++ {
		copy(packet[offset:], mac)
		offset += 6
	}

	return packet, nil
}

func sendUDPPacket(mp MagicPacket, addr string) (err error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(mp[:])
	return err
}

// Send writes the MagicPacket to the specified address on port 9.
func (mp MagicPacket) send(addr string) error {
	return sendUDPPacket(mp, addr+":9")
}
