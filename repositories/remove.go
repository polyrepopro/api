package repositories

import (
	"github.com/polyrepopro/api/config"
)

func Remove(args config.Repository) error {
	cfg, err := config.GetRelativeConfig()
	if err != nil {
		return err
	}

	workspace, err := cfg.GetWorkspaceByWorkingDir()
	if err != nil {
		return err
	}

	newRepositories := []config.Repository{}
	for _, repository := range *workspace.Repositories {
		if repository.Path != args.Path {
			newRepositories = append(newRepositories, repository)
		}
	}
	*workspace.Repositories = newRepositories

	return cfg.SaveConfig()
}
