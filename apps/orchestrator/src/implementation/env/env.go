package env

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/s3"
)

type EnvService struct {
	logger *logger.Logger

	envFilesDir string

	s3Config *s3.Config
}

var (
	instance *EnvService
)

func InitEnvService(s3Config *s3.Config) *EnvService {
	tempEnvFilesDir, err := os.MkdirTemp("", "env_files")
	if err != nil {
		fmt.Printf("Error creating temporary directory for env files: %v\n", err)
		os.Exit(1)
	}

	instance = &EnvService{
		logger:      logger.NewLogger("[SERVICE:ENV]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil),
		s3Config:    s3Config,
		envFilesDir: tempEnvFilesDir,
	}

	err = instance.RefreshEnvFiles(context.Background())
	if err != nil {
		instance.logger.Errorf("Error initializing env service: %v", err)
		os.Exit(1)
	}

	return instance
}

func GetEnvService() *EnvService {
	if instance == nil {
		fmt.Println("EnvService is not initialized. Please call InitEnvService first.")
		os.Exit(1)
	}

	return instance
}

func (service *EnvService) Destroy() {
	if service.envFilesDir != "" {
		err := os.RemoveAll(service.envFilesDir)
		if err != nil {
			service.logger.Errorf("Error removing temporary env files directory: %v", err)
		} else {
			service.logger.Debugf("Temporary env files directory removed: %s", service.envFilesDir)
		}
	}
}

func (service *EnvService) GetEnvFilesDir() string {
	return service.envFilesDir
}

func (service *EnvService) RefreshEnvFiles(ctx context.Context) error {

	client, err := s3.NewClient(ctx, *service.s3Config)
	if err != nil {
		return fmt.Errorf("error creating S3 client: %w", err)
	}

	// 1. List all objects in the bucket
	objects, err := client.ListObjects(ctx, "")
	if err != nil {
		return fmt.Errorf("error listing objects in bucket: %w", err)
	}

	for _, object := range objects {
		service.logger.Debugf("Checking object: %s", object.Key)

		// 2. Ensure it's a .env file
		if !strings.HasSuffix(object.Key, ".env") {
			continue
		}

		localPath := path.Join(service.envFilesDir, object.Key)

		// 3. Ensure local subdirectories exist (e.g., targetDir/stack_1/)
		if err := os.MkdirAll(path.Dir(localPath), 0755); err != nil {
			service.logger.Errorf("Error creating local directory for %s: %v", localPath, err)
			continue
		}

		// 4. Download the file
		err := client.DownloadToFile(ctx, object.Key, localPath)
		if err != nil {
			service.logger.Errorf("Error downloading file %s: %v", object.Key, err)
			continue
		}
		service.logger.Infof("Downloaded file %s to %s", object.Key, localPath)
	}

	return nil
}
