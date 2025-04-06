package sftp

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"mgarnier11.fr/go/libs/utils"
)

func LocalToRemoteProgress(sshClient *ssh.Client, localDirPath, remoteDirPath string, progressFunc func(int64, float64, int64)) error {
	// Create SFTP client
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %v", err)
	}
	defer sftpClient.Close()

	// Copy the directory
	return copyDirectoryToRemote(localDirPath, remoteDirPath, sftpClient, progressFunc)
}

func LocalToRemote(sshClient *ssh.Client, localDirPath, remoteDirPath string) error {
	return LocalToRemoteProgress(sshClient, localDirPath, remoteDirPath, nil)
}

func copyDirectoryToRemote(localDir, remoteDir string, sftpClient *sftp.Client, progressFunc func(int64, float64, int64)) error {
	totalSize, err := utils.GetDirSize(localDir)
	if err != nil {
		return fmt.Errorf("failed to get directory size: %v", err)
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

		return utils.ParallelCopyFile(
			localPath,
			remotePath,
			func(s string) (utils.ReadWriterAt, error) { return os.Open(s) },
			func(s string) (utils.ReadWriterAt, error) { return sftpClient.Create(s) },
			fileProgress,
		)
	})
}
