package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/sftp"
	"mgarnier11.fr/go/go-autosaver/config"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/sshutils"
	sftpUtils "mgarnier11.fr/go/libs/sshutils/sftp"
	"mgarnier11.fr/go/libs/utils"
)

func osReadDir(path string) ([]os.FileInfo, error) {
	//function to cast os.ReadDir to os.FileInfo
	// os.ReadDir returns a slice of os.DirEntry, we need to convert it to []os.FileInfo
	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var fileInfos []os.FileInfo
	for _, dir := range dirs {
		info, err := dir.Info()
		if err != nil {
			return nil, err
		}
		fileInfos = append(fileInfos, info)
	}
	return fileInfos, nil
}

func getDateOfDay() string {
	return fmt.Sprintf("%d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day())
}

func createBackupFolder(
	stat func(string) (os.FileInfo, error),
	dirPath string,
) error {

	_, err := stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
			}
		} else {
			return fmt.Errorf("failed to stat directory %s: %w", dirPath, err)
		}
	}
	return nil
}

func deleteOldFolders(
	readDir func(string) ([]os.FileInfo, error),
	removeAll func(string) error,
	dirPath string,
) error {
	directories, err := readDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	logger.Infof("Checking %d directories in %s", len(directories), dirPath)

	for _, dir := range directories {
		if dir.IsDir() {
			// check if the dir name corresponds to the date format	using a regex
			dirName := dir.Name()
			if dateRegex.MatchString(dirName) {
				//check if the dir is older than 14 days
				dirDate, err := time.Parse("2006-01-02", dirName)
				if err != nil {
					return fmt.Errorf("failed to parse directory name %s as date: %w", dirName, err)
				}

				_14days := 14 * 24 * time.Hour
				_1year := 365 * 24 * time.Hour

				// delete the directory if it is older than 14 days and not the first day of the month, delete if it is older than 1 year
				if time.Since(dirDate) > _14days && dirDate.Day() != 1 || time.Since(dirDate) > _1year {
					// delete the directory
					err = removeAll(filepath.Join(dirPath, dirName))
					if err != nil {
						return fmt.Errorf("failed to remove directory %s: %w", dirName, err)
					}
				}
			}
		}
	}

	return nil
}

func copyToRemote(remoteDest *config.RemoteDestConfig, srcFile string) error {
	logger.Infof("Copying backup to remote dest")

	sshClient, err := sshutils.GetSSHClient(
		remoteDest.SSHUser,
		remoteDest.SSHHost,
		strconv.Itoa(remoteDest.SSHPort),
		remoteDest.SSHKeyPath,
	)
	if err != nil {
		return fmt.Errorf("failed to get SSH client: %w", err)
	}
	defer sshClient.Close()

	logger.Infof("Connected to remote dest")

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer sftpClient.Close()

	err = createBackupFolder(
		sftpClient.Stat,
		filepath.Join(remoteDest.SSHPath, getDateOfDay()),
	)
	if err != nil {
		return fmt.Errorf("failed to create backup folder: %w", err)
	}

	err = deleteOldFolders(
		sftpClient.ReadDir,
		sftpClient.RemoveAll,
		remoteDest.SSHPath,
	)
	if err != nil {
		return fmt.Errorf("failed to delete old folders: %w", err)
	}

	lastCopyPercent := 0.0

	err = sftpUtils.LocalToRemoteProgress(
		sshClient,
		srcFile,
		filepath.Join(remoteDest.SSHPath, getDateOfDay(), filepath.Base(srcFile)),
		func(current int64, percent float64, total int64) {
			if percent-lastCopyPercent > 1 {
				lastCopyPercent = percent
				logger.Infof("Copying backup to remote dest: %d", int(percent))

			}
		},
	)

	if err != nil {
		return fmt.Errorf("failed to copy backup: %w", err)
	}
	logger.Infof("Successfully copied backup to remote dest")

	return nil
}

func copyToLocal(localDest, srcFile string) error {
	logger.Infof("Copying backup to local dest")

	err := createBackupFolder(
		os.Stat,
		filepath.Join(localDest, getDateOfDay()),
	)
	if err != nil {
		return fmt.Errorf("failed to create backup folder: %w", err)
	}

	err = deleteOldFolders(
		osReadDir,
		os.RemoveAll,
		filepath.Join(localDest, getDateOfDay()),
	)
	if err != nil {
		return fmt.Errorf("failed to delete old folders: %w", err)
	}

	lastCopyPercent := 0.0

	err = utils.ParallelCopyFile(
		srcFile,
		filepath.Join(localDest, getDateOfDay(), filepath.Base(srcFile)),
		func(s string) (utils.ReadWriterAt, error) { return os.Open(s) },
		func(s string) (utils.ReadWriterAt, error) { return os.Create(s) },
		func(written int, totalWritten, totalSize int64) {
			if totalWritten > 0 {
				percent := float64(totalWritten) / float64(totalSize) * 100.0

				if percent-lastCopyPercent > 1 {
					lastCopyPercent = percent
					logger.Infof("Copying backup to local dest: %d", int(percent))
				}
			}
		},
	)

	if err != nil {
		return fmt.Errorf("failed to copy backup: %w", err)
	}

	logger.Infof("Successfully copied backup to local dest")

	return nil
}
