package repositories

import (
	"fmt"
	"os"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
)

func Doctor(repository *config.Repository, workspace *config.Workspace) error {
	path := fmt.Sprintf("%s/%s", workspace.GetAbsolutePath(), repository.Path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		multilog.Error("repositories.doctor", "repository directory does not exist", map[string]interface{}{
			"repository": repository,
		})

		err := git.Clone(git.CloneArgs{
			URL:  repository.URL,
			Path: path,
		})
		if err != nil {
			multilog.Fatal("repositories.doctor", "failed to clone repository", map[string]interface{}{
				"repository": repository,
			})
			return err
		}
	}

	return nil
}
