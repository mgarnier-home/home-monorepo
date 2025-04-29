package backup

import (
	"context"
	"fmt"
	"time"

	"mgarnier11.fr/go/go-autosaver/config"
	"mgarnier11.fr/go/go-autosaver/external"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/ntfy"
	"mgarnier11.fr/go/libs/utils"
)

var running bool = false

func RunSave(appConfig *config.AppConfigFile) bool {
	if running {
		return false
	}
	running = true

	go func() {
		ctx, cancel := context.WithCancel(context.Background())

		defer func() {
			running = false
			cancel()
		}()

		if appConfig.KeepAliveUrl != "" {
			logger.Infof("Starting keep alive url: %s", appConfig.KeepAliveUrl)
			go utils.RunPeriodic(ctx, 30*time.Second, func() {
				logger.Infof("Keep alive url: %s", appConfig.KeepAliveUrl)
			})
		}

		err := save(appConfig)
		if err != nil {
			logger.Errorf("Failed to save: %s", err)

			err = external.SendMail(
				appConfig.Mail,
				appConfig.Mail.ErrorTo,
				fmt.Sprintf("Error for %s of %s", appConfig.FileName, utils.GetDateOfDay()),
				fmt.Sprintf("Error: %s", err),
			)
			if err != nil {
				logger.Errorf("Failed to send mail: %s", err)
			}

			err = ntfy.SendNotification(
				"Autosaver",
				fmt.Sprintf("Backup of %s failed 🔴", appConfig.FileName),
				"bomb",
			)
			if err != nil {
				logger.Errorf("Failed to send ntfy notification: %s", err)
			}

			return
		}
		logger.Infof("Successfully saved")
	}()

	return true
}

func save(appConfig *config.AppConfigFile) error {

	logger.Infof("Starting backup")
	var err error

	err = zipFolder(appConfig.BackupSrc, appConfig.FileName)
	if err != nil {
		return err
	}

	encryptedFileName := appConfig.FileName + ".gpg"

	password, err := encryptFile(appConfig.FileName, encryptedFileName)
	if err != nil {
		return err
	}

	err = external.SendMail(
		appConfig.Mail,
		appConfig.Mail.InfoTo,
		fmt.Sprintf("Infos for %s of %s", appConfig.FileName, utils.GetDateOfDay()),
		fmt.Sprintf("Archive password is : %s", password),
	)
	if err != nil {
		logger.Errorf("Failed to send mail: %s", err)
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

	err = ntfy.SendNotification(
		"Autosaver",
		fmt.Sprintf("Backup of %s success 🟢", appConfig.FileName),
		"partying_face",
	)
	if err != nil {
		logger.Errorf("Failed to send ntfy notification: %s", err)
	}

	return nil

}
