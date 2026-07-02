package main

import (
	"os"

	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/s3"
	"mgarnier11.fr/go/orchestrator-api/config"
	"mgarnier11.fr/go/orchestrator-api/server"
	common "mgarnier11.fr/go/orchestrator-common"
)

func main() {
	logger.InitAppLogger("dashboard")

	if config.Env.SSHPrivateKey != "" {
		// Create /ssh directory if it doesn't exist
		err := os.MkdirAll("ssh", 0700)
		if err != nil {
			panic(err)
		}
		// Create the id_rsa file with the content of SSH_PRIVATE_KEY
		err = os.WriteFile("ssh/ssh_private_key", []byte(config.Env.SSHPrivateKey), 0600)
		if err != nil {
			panic(err)
		}
	}

	commonLib := common.NewCommonLib(
		config.Env.ComposeDirPath,
		&s3.Config{
			Endpoint:        config.Env.S3Endpoint,
			AccessKeyID:     config.Env.S3AccessKey,
			SecretAccessKey: config.Env.S3SecretKey,
			Bucket:          config.Env.S3Bucket,
		},
	)

	api, err := server.NewServer(config.Env.ServerPort, commonLib)

	if err != nil {
		panic(err)
	}

	err = api.Start()
	if err != nil {
		panic(err)
	}
}
