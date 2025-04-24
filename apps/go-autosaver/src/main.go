package main

import (
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/sshutils"

	"mgarnier11.fr/go/go-autosaver/backup"
)

func main() {
	logger.InitAppLogger("")

	// api := server.NewServer(3000)

	// api.Start()

	sshClient, err := sshutils.GetSSHClient("u368422", "u368422.your-storagebox.de", "23", "./id_rsa")

	if err != nil {
		panic(err)
	}
	defer sshClient.Close()

	// lastPercent := 0.0
	// startTime := time.Now()

	// err = sftp.LocalToRemoteProgress(
	// 	sshClient,
	// 	"./test-send",
	// 	"./test-send",
	// 	func(current int64, percent float64, total int64) {
	// 		if percent-lastPercent > 0.1 {
	// 			lastPercent = percent
	// 			elapsed := time.Since(startTime)
	// 			speed := float64(current) / elapsed.Seconds() / 1024 / 1024

	// 			logger.Infof("Progress: %0.1f, %0.1f MB/s", percent, speed)
	// 		}
	// 	},
	// )

	// err = sftp.RemoteToLocalProgress(
	// 	sshClient,
	// 	"test-send",
	// 	"test-receive",
	// 	func(current int64, percent float64, total int64) {
	// 		if percent-lastPercent > 0.1 {
	// 			lastPercent = percent
	// 			elapsed := time.Since(startTime)
	// 			speed := float64(current) / elapsed.Seconds() / 1024 / 1024

	// 			logger.Infof("Progress: %0.1f, %0.1f MB/s", percent, speed)
	// 		}
	// 	},
	// )
	// if err != nil {
	// 	panic(err)
	// }

	// logger.Infof("File uploaded successfully")

	// backup.EncryptFileWithPassword(
	// 	"./10GB.bin",
	// 	"./10GB.bin.gpg",
	// 	"test",
	// 	func(bytesRead int, totalBytesRead int64, totalSize int64) {
	// 		if totalBytesRead > 0 {
	// 			percent := float64(totalBytesRead) / float64(totalSize) * 100.0

	// 			if percent-lastPercent > 0.1 {
	// 				lastPercent = percent

	// 				elapsed := time.Since(startTime)
	// 				speed := float64(totalBytesRead) / elapsed.Seconds() / 1024 / 1024

	// 				logger.Infof("Progress: %0.1f, %0.1f MB/s", percent, speed)
	// 			}
	// 		}
	// 	},
	// )

	backup.ZipFolder("./test-send", "./test-send.zip", func(filePath string, bytesRead int, totalBytesRead int64, totalSize int64) {
		logger.Infof("Zipping : %s, %d bytes read, %d total bytes read, %d total size", filePath, bytesRead, totalBytesRead, totalSize)
	})

	logger.Infof("File encrypted successfully")

}
