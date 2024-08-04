package repositories

import (
	"fmt"

	"github.com/polyrepopro/api/config"
)

func Get(path string) (config.Repository, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return config.Repository{}, err
	}

	workspace, err := cfg.GetWorkspaceByWorkingDir()
	if err != nil {
		return config.Repository{}, err
	}

	for _, repository := range *workspace.Repositories {
		if repository.Path == path {
			return repository, nil
		}
	}

	return config.Repository{}, fmt.Errorf("repository not found with path: %s", path)
}
