package repositories

import (
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
)

type StatusResult struct {
	Dirty bool
}

func Status(repository *config.Repository) (StatusResult, error) {
	repoPath := repository.GetAbsolutePath()
	status, err := git.Status(repoPath)
	if err != nil {
		return StatusResult{}, err
	}

	return StatusResult{
		Dirty: !status.IsClean(),
	}, nil
}
