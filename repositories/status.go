package repositories

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/polyrepopro/api/config"
)

type StatusResult struct {
	Dirty bool
}

func Status(repository *config.Repository) (StatusResult, error) {
	repoPath := repository.GetAbsolutePath()
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return StatusResult{}, fmt.Errorf("failed to open repository: %w", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return StatusResult{}, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := w.Status()
	if err != nil {
		return StatusResult{}, fmt.Errorf("failed to get status: %w", err)
	}

	return StatusResult{
		Dirty: !status.IsClean(),
	}, nil
}
