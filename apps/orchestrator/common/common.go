package common

import (
	"github.com/charmbracelet/lipgloss"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/s3"
	"mgarnier11.fr/go/orchestrator-common/config"
	"mgarnier11.fr/go/orchestrator-common/exec"
	"mgarnier11.fr/go/orchestrator-common/files"
)

type CommonLib struct {
	composeDir string
	s3Config   *s3.Config
	Exec       *exec.Exec
	Files      *files.Files
	Config     *config.Config
	logger     *logger.Logger
}

func NewCommonLib(composeDir string, s3Config *s3.Config) *CommonLib {
	logger := logger.NewLogger("[COMMON-LIB]", "%-10s ", lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), nil)
	files := files.NewFiles(composeDir)
	config := config.NewConfig(composeDir, s3Config)
	exec := exec.NewExec()

	logger.Infof("CommonLib initialized with composeDir: %s", composeDir)

	return &CommonLib{
		composeDir: composeDir,
		s3Config:   s3Config,
		Exec:       exec,
		Files:      files,
		Config:     config,
		logger:     logger,
	}
}
