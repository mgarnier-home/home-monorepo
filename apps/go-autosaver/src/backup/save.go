package backup

import (
	"fmt"
	"math"
	"strconv"

	"golang.org/x/crypto/ssh"
	"mgarnier11.fr/go/go-autosaver/config"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/sshutils"
	"mgarnier11.fr/go/libs/sshutils/sftp"
	"mgarnier11.fr/go/libs/utils"
)

var running bool = false

func RunSave(appConfig *config.AppConfigFile) bool {
	if running {
		return false
	}
	running = true

	go func() {

		err := save(appConfig)
		if err != nil {
			logger.Errorf("Failed to save: %s", err)
			return
		}
		logger.Infof("Successfully saved")
	}()

	return true
}

func zipFolder(backupSrc string) error {
	filePercent, lastFilePercent := 0.0, 0.0
	totalPercent, lastTotalPercent := 0.0, 0.0

	logger.Infof("Zipping folder %s", backupSrc)

	err := ZipFolder(
		backupSrc,
		"backup.zip",
		func(
			fileName string,
			written int,
			fileWritten,
			fileSize,
			totalWritten,
			totalSize int64,
		) {
			filePercent = float64(fileWritten) / float64(fileSize) * 100
			totalPercent = float64(totalWritten) / float64(totalSize) * 100

			if math.Abs(filePercent-lastFilePercent) > 1 {
				lastFilePercent = filePercent
				logger.Debugf("Zipping file %s: %d", fileName, int(filePercent))
			}

			if totalPercent-lastTotalPercent > 1 {
				lastTotalPercent = totalPercent
				logger.Infof("Zipping folder: %d", int(totalPercent))
			}

		})

	if err != nil {
		return fmt.Errorf("failed to zip folder: %w", err)
	}

	logger.Infof("Successfully zipped folder")

	return nil
}

func encryptBackup(zipFile, outputFile string) (string, error) {
	password, err := utils.GenerateRandomString(20)

	if err != nil {
		return "", fmt.Errorf("failed to generate password: %w", err)
	}

	logger.Infof("Encrypting backup with password: %s", password)

	encryptPercent, lastEncryptPercent := 0.0, 0.0

	err = EncryptFileWithPassword(
		zipFile,
		outputFile,
		password,
		func(written int, totalWritten int64, totalSize int64) {
			encryptPercent = float64(totalWritten) / float64(totalSize) * 100
			if encryptPercent-lastEncryptPercent > 1 {
				lastEncryptPercent = encryptPercent
				logger.Infof("Encrypting backup: %d", int(encryptPercent))
			}
		})

	if err != nil {
		return "", fmt.Errorf("failed to encrypt backup: %w", err)
	}

	logger.Infof("Successfully encrypted backup")

	return password, nil
}

func copyBackup(sshClient *ssh.Client, srcFile, destFile string) error {
	logger.Infof("Copying backup to remote server")

	lastCopyPercent := 0.0

	err := sftp.LocalToRemoteProgress(
		sshClient,
		srcFile,
		destFile,
		func(current int64, percent float64, total int64) {
			if percent-lastCopyPercent > 1 {
				lastCopyPercent = percent
				logger.Infof("Copying backup: %d", int(percent))

			}
		},
	)

	if err != nil {
		return fmt.Errorf("failed to copy backup: %w", err)
	}
	logger.Infof("Successfully copied backup")

	return nil
}

func save(appConfig *config.AppConfigFile) error {
	defer func() { running = false }()

	logger.Infof("Starting backup")

	err := zipFolder(appConfig.BackupSrc)
	if err != nil {
		return err
	}

	_, err = encryptBackup("./backup.zip", "./backup.zip.gpg")
	if err != nil {
		return err
	}

	sshClient, err := sshutils.GetSSHClient(
		appConfig.BackupDest.SSHUser,
		appConfig.BackupDest.SSHHost,
		strconv.Itoa(appConfig.BackupDest.SSHPort),
		appConfig.BackupDest.SSHKeyPath,
	)
	if err != nil {
		return fmt.Errorf("failed to get SSH client: %w", err)
	}
	defer sshClient.Close()

	err = copyBackup(
		sshClient,
		"./backup.zip.gpg",
		appConfig.BackupDest.SSHPath+"/backup.zip.gpg",
	)

	if err != nil {
		return err
	}

	return nil

}
