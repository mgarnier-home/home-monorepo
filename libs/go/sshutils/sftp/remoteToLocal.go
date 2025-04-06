package sftp

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"mgarnier11.fr/go/libs/utils"
)

func RemoveSFTPDir(sshClient *ssh.Client, dirPath string) error {
	// Create SFTP client
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %v", err)
	}
	defer sftpClient.Close()

	log.Println("Removing directory", dirPath)

	// Remove the directory
	return sftpClient.RemoveAll(dirPath)
}

func GetSFTPDirSize(sftpClient *sftp.Client, dirPath string) (int64, error) {
	walker := sftpClient.Walk(dirPath)

	var size int64
	for walker.Step() {
		if walker.Err() != nil {
			return 0, fmt.Errorf("failed to walk remote directory: %v", walker.Err())
		}

		if walker.Stat().IsDir() {
			continue
		}

		fileInfo := walker.Stat()
		if fileInfo == nil {
			return 0, fmt.Errorf("failed to get file info for %s: %v", walker.Path(), walker.Err())
		}

		size += fileInfo.Size()
	}
	return size, nil
}

func RemoteToLocalProgress(sshClient *ssh.Client, remoteDirPath, localDirPath string, progressFunc func(int64, float64, int64)) error {
	// Create SFTP client
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %v", err)
	}
	defer sftpClient.Close()

	// Copy the directory
	return copyDirectoryToLocal(localDirPath, remoteDirPath, sftpClient, progressFunc)
}

func RemoteToLocal(sshClient *ssh.Client, localDirPath, remoteDirPath string) error {
	return RemoteToLocalProgress(sshClient, localDirPath, remoteDirPath, nil)
}

func copyDirectoryToLocal(localDir, remoteDir string, sftpClient *sftp.Client, progressFunc func(int64, float64, int64)) error {
	totalSize, err := GetSFTPDirSize(sftpClient, remoteDir)
	if err != nil {
		return fmt.Errorf("failed to get remote directory size: %v", err)
	}
	copiedSize := int64(0)

	progressFunc(0, 0.0, totalSize)

	fileProgress := func(n int, copiedFileSize, totalFileSize int64) {
		copiedSize += int64(n)
		percent := float64(copiedSize) / float64(totalSize) * 100.0

		if progressFunc != nil {
			progressFunc(copiedSize, percent, totalSize)
		}
	}

	walker := sftpClient.Walk(remoteDir)

	for walker.Step() {
		if walker.Err() != nil {
			return fmt.Errorf("failed to walk remote directory: %v", walker.Err())
		}

		remotePath := walker.Path()
		localPath := filepath.Join(localDir, strings.TrimPrefix(remotePath, remoteDir))

		println("Copying", remotePath, "to", localPath)

		if walker.Stat().IsDir() {
			err := os.MkdirAll(localPath, 0755)
			if err != nil {
				return fmt.Errorf("failed to create local directory: %v", err)
			}
		} else {
			// err := copyFileToLocal(localPath, remotePath, sftpClient)
			err := utils.ParallelCopyFile(
				remotePath,
				localPath,
				func(path string) (utils.ReadWriterAt, error) { return sftpClient.Open(path) },
				func(name string) (utils.ReadWriterAt, error) { return os.Create(name) },
				fileProgress,
			)
			if err != nil {
				return fmt.Errorf("failed to copy file to local: %v", err)
			}
		}
	}

	return nil
}
