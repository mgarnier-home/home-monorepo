package backup

import (
	"mgarnier11.fr/go/go-autosaver/config"
	"mgarnier11.fr/go/libs/logger"
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

func save(appConfig *config.AppConfigFile) error {
	defer func() { running = false }()

	logger.Infof("Starting backup")
	var err error

	err = zipFolder(appConfig.BackupSrc)
	if err != nil {
		return err
	}

	_, err = encryptFile("./backup.zip", "./backup.zip.gpg")
	if err != nil {
		return err
	}

	if appConfig.LocalDest != "" {
		err = copyToLocal(
			appConfig.LocalDest,
			"./backup.zip",
		)
		if err != nil {
			return err
		}
	}

	if appConfig.RemoteDest != nil {
		err = copyToRemote(
			appConfig.RemoteDest,
			"./backup.zip.gpg",
		)
		if err != nil {
			return err
		}
	}

	return nil

}
