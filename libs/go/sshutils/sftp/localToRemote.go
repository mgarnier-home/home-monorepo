package sftp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func LocalToRemote(sshClient *ssh.Client, localDirPath, remoteDirPath string) error {
	// Create SFTP client
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %v", err)
	}
	defer sftpClient.Close()

	// Copy the directory
	return copyDirectoryToRemote(localDirPath, remoteDirPath, sftpClient)
}

func copyDirectoryToRemote(localDir, remoteDir string, sftpClient *sftp.Client) error {
	return filepath.Walk(localDir, func(localPath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk path %s: %v", localPath, err)
		}

		// Determine the relative path to preserve directory structure
		relPath, err := filepath.Rel(localDir, localPath)
		if err != nil {
			return err
		}
		remotePath := filepath.Join(remoteDir, relPath)

		if info.IsDir() {
			// Create remote directory
			return sftpClient.MkdirAll(remotePath)
		}

		return copyFileToRemote(localPath, remotePath, sftpClient)
	})
}

func copyFileToRemote(localPath, remotePath string, sftpClient *sftp.Client) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file %s: %v", localPath, err)
	}
	defer localFile.Close()

	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file %s: %v", remotePath, err)
	}

	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	return nil
}
