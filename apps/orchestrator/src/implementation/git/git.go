package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func gitAuth(token string) *gitHttp.BasicAuth {
	return &gitHttp.BasicAuth{
		Username: "git", // "git" works as a placeholder username for most Git providers
		Password: token,
	}
}

func gitClone(dir string, repoUrl string, token string) (*git.Repository, error) {

	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:      repoUrl,
		Auth:     gitAuth(token),
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	return repo, nil
}

func GitPull(dir string, repoUrl string, token string) error {
	repo, err := gitClone(dir, repoUrl, token)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = worktree.Pull(&git.PullOptions{
		Auth:     gitAuth(token),
		Progress: os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to pull repository: %w", err)
	}

	return nil
}
