package controllers

import (
	"fmt"
	"mgarnier11/go/sshutils"
	"mgarnier11/go/sshutils/sftp"
	"mgarnier11/mineager/config"
	"mgarnier11/mineager/server/objects/bo"
	"os"
)

func getServerLocalMapPath(serverName string) string {
	return fmt.Sprintf("%s/map", getServerLocalPath(serverName))
}

func getServerLocalPath(serverName string) string {
	return fmt.Sprintf("%s/%s", config.Config.ServersFolderPath, serverName)
}

func createServerDirectory(serverName string) error {
	return os.MkdirAll(getServerLocalPath(serverName), 0755)
}

func deleteServerDirectory(serverName string) error {
	return os.RemoveAll(getServerLocalPath(serverName))
}

func getServerHostMapPath(host *bo.HostBo, serverName string) string {
	return fmt.Sprintf("%s/world", getServerHostPath(host, serverName))
}

func getServerHostPath(host *bo.HostBo, serverName string) string {
	return fmt.Sprintf("%s/%s", host.MineagerPath, serverName)
}

func getMapPath(mapName string) string {
	return fmt.Sprintf("%s/%s", config.Config.MapsFolderPath, mapName)
}

func sendServerMapToHost(serverName string, host *bo.HostBo) error {
	sshClient, err := sshutils.GetSSHClient(host.SSHUsername, host.Ip, host.SSHPort, config.Config.SSHKeyPath)
	if err != nil {
		return fmt.Errorf("error connecting to ssh: %v", err)
	}
	defer sshClient.Close()

	return sftp.LocalToRemote(
		sshClient,
		getServerLocalMapPath(serverName),
		getServerHostMapPath(host, serverName),
	)
}

func getServerMapFromHost(serverName string, host *bo.HostBo) error {
	sshClient, err := sshutils.GetSSHClient(host.SSHUsername, host.Ip, host.SSHPort, config.Config.SSHKeyPath)
	if err != nil {
		return fmt.Errorf("error connecting to ssh: %v", err)
	}
	defer sshClient.Close()

	return sftp.RemoteToLocal(
		sshClient,
		getServerLocalMapPath(serverName),
		getServerHostMapPath(host, serverName),
	)
}

func deleteHostDirectory(serverName string, host *bo.HostBo) error {
	sshClient, err := sshutils.GetSSHClient(host.SSHUsername, host.Ip, host.SSHPort, config.Config.SSHKeyPath)
	if err != nil {
		return fmt.Errorf("error connecting to ssh: %v", err)
	}
	defer sshClient.Close()

	return sftp.RemoveDir(sshClient, getServerHostPath(host, serverName))
}
