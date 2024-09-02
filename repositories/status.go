package repositories

import (
	"github.com/polyrepopro/api/git"
)

type StatusResult struct {
	Dirty bool
}

func Status(path string) (StatusResult, error) {
	status, err := git.Status(path)
	if err != nil {
		return StatusResult{}, err
	}

	return StatusResult{
		Dirty: !status.IsClean(),
	}, nil
}
