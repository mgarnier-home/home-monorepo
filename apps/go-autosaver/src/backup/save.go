package backup

import (
	"mgarnier11.fr/go/go-autosaver/config"
	"mgarnier11.fr/go/go-autosaver/external"
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

	err = zipFolder(appConfig.BackupSrc, appConfig.FileName)
	if err != nil {
		return err
	}

	encryptedFileName := appConfig.FileName + ".gpg"

	_, err = encryptFile(appConfig.FileName, encryptedFileName)
	if err != nil {
		return err
	}

	if appConfig.LocalDest != "" {
		err = external.CopyToLocal(
			appConfig.LocalDest,
			appConfig.FileName,
		)
		if err != nil {
			return err
		}
	}

	if appConfig.RemoteDest != nil {
		err = external.CopyToRemote(
			appConfig.RemoteDest,
			appConfig.FileName,
		)
		if err != nil {
			return err
		}
	}

	return nil

}
