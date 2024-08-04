package repositories

import (
	"os"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
)

func Doctor(repository *config.Repository) error {
	path := repository.GetAbsolutePath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		multilog.Error("repositories.doctor", "repository directory does not exist", map[string]interface{}{
			"repository": repository,
		})

		err := git.Clone(git.CloneArgs{
			URL:  repository.URL,
			Path: repository.GetAbsolutePath(),
		})
		if err != nil {
			multilog.Fatal("repositories.doctor", "failed to clone repository", map[string]interface{}{
				"repository": repository,
			})
		}

		multilog.Info("repositories.doctor", "cloned repository", map[string]interface{}{
			"url":  repository.URL,
			"path": path,
		})
	}

	return nil
}
