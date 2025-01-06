package sftp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func RemoteToLocal(sshClient *ssh.Client, localDirPath, remoteDirPath string) error {
	// Create SFTP client
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %v", err)
	}
	defer sftpClient.Close()

	// Copy the directory
	return copyDirectoryToLocal(localDirPath, remoteDirPath, sftpClient)
}

func copyDirectoryToLocal(localDir, remoteDir string, sftpClient *sftp.Client) error {
	walker := sftpClient.Walk(remoteDir)

	for walker.Step() {
		if walker.Err() != nil {
			return fmt.Errorf("failed to walk remote directory: %v", walker.Err())
		}

		remotePath := walker.Path()
		localPath := filepath.Join(localDir, strings.TrimPrefix(remotePath, remoteDir))

		if walker.Stat().IsDir() {
			err := os.MkdirAll(localPath, 0755)
			if err != nil {
				return fmt.Errorf("failed to create local directory: %v", err)
			}
		} else {
			err := copyFileToLocal(localPath, remotePath, sftpClient)
			if err != nil {
				return fmt.Errorf("failed to copy file to local: %v", err)
			}
		}
	}

	return nil
}

func copyFileToLocal(localPath, remotePath string, sftpClient *sftp.Client) error {

	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %v", err)
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer localFile.Close()

	_, err = remoteFile.WriteTo(localFile)
	if err != nil {
		return fmt.Errorf("failed to write remote file to local: %v", err)
	}

	return nil
}
