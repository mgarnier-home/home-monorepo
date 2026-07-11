package compose

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/orchestrator/implementation/git"
	"mgarnier11.fr/go/orchestrator/models"
)

type ComposeService struct {
	config *models.OrchestratorConfig
	logger *logger.Logger

	composeFileDir string
	composeFiles   []*models.ComposeFile
}

var (
	instance *ComposeService
)

func InitComposeService(config *models.OrchestratorConfig) *ComposeService {
	var (
		composeFileDir string
		err            error
	)

	if config.GitRepo != "" {
		composeFileDir, err = os.MkdirTemp("", "compose-files-*")
		if err != nil {
			fmt.Printf("Error creating temporary directory for compose files: %v\n", err)
			os.Exit(1)
		}
	} else {
		composeFileDir = config.ComposeDirPath
	}

	instance = &ComposeService{
		config:         config,
		logger:         logger.NewLogger("[SERVICE:COMPOSE]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil),
		composeFiles:   []*models.ComposeFile{},
		composeFileDir: composeFileDir,
	}

	err = instance.RefreshComposeFiles()
	if err != nil {
		instance.logger.Errorf("Error initializing compose service: %v", err)
		os.Exit(1)
	}

	return instance
}

func GetComposeService() *ComposeService {
	if instance == nil {
		fmt.Println("ComposeService is not initialized. Please call InitComposeService first.")
		os.Exit(1)
	}

	return instance
}

func (service *ComposeService) Destroy() {
	if service.composeFileDir != "" && service.config.ComposeDirPath == "" {
		err := os.RemoveAll(service.composeFileDir)
		if err != nil {
			service.logger.Errorf("Error removing temporary compose files directory: %v", err)
		}
	}
}

func (service *ComposeService) GetComposeFiles() []*models.ComposeFile {
	return service.composeFiles
}

func (service *ComposeService) GetComposeFilesDir() string {
	return service.composeFileDir
}

func (service *ComposeService) RefreshComposeFiles() error {
	var (
		composeFiles []*models.ComposeFile
		err          error
	)
	if service.config.GitRepo != "" {
		composeFiles, err = service.getComposeFilesFromGitRepo(service.composeFileDir, service.config.GitRepo)
	} else {
		composeFiles, err = service.getComposeFilesFromDir(service.composeFileDir)
	}

	if err != nil {
		return fmt.Errorf("error refreshing compose files: %w", err)
	}

	service.composeFiles = composeFiles
	return nil
}

func (service *ComposeService) getComposeFilesFromGitRepo(dir string, gitRepo string) ([]*models.ComposeFile, error) {
	service.logger.Infof("Getting compose files from git repository: %s %s", gitRepo, dir)
	err := git.GitPull(dir, gitRepo, service.config.GitToken)
	if err != nil {
		return nil, err
	}

	composeFiles, err := service.getComposeFilesFromDir(dir)
	if err != nil {
		return nil, err
	}

	return composeFiles, nil
}

func (service *ComposeService) getComposeFilesFromDir(dir string) ([]*models.ComposeFile, error) {
	service.logger.Debugf("Getting compose files from directory: %s", dir)

	stacks, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	composeFiles := []*models.ComposeFile{}

	for _, stack := range stacks {

		if !stack.IsDir() {
			continue
		}

		stackName := stack.Name()

		service.logger.Verbosef("Found stack: %s", stackName)

		stackPath := path.Join(dir, stackName)
		hostsFiles, err := os.ReadDir(stackPath)

		if err != nil {
			service.logger.Errorf("Error reading directory %s: %v", stackPath, err)
			continue
		}

		for _, hostFile := range hostsFiles {
			hostFileName := hostFile.Name()

			parts := strings.Split(hostFileName, ".")

			if len(parts) != 3 || parts[1] != stackName || parts[2] != "yml" {
				continue
			}

			hostName := parts[0]

			service.logger.Verbosef("Found compose file: %s for host: %s for stack: %s", hostFileName, hostName, stackName)

			composeFiles = append(composeFiles, composeFile(dir, stackName, hostName))
		}
	}

	return composeFiles, nil
}

func composeFile(dir, stackName, hostName string) *models.ComposeFile {
	stackPath := path.Join(dir, stackName)

	return &models.ComposeFile{
		Name:  hostName + "-" + stackName,
		Stack: stackName,
		Host:  hostName,
		Path:  path.Join(stackPath, hostName+"."+stackName+".yml"),
	}
}
