package workspaces

import (
	"fmt"
	"os"

	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
)

type SyncArgs struct {
}

func Sync(workspace *config.Workspace, args SyncArgs) error {
	for _, repo := range *workspace.Repositories {
		// repoPath := repo.GetAbsolutePath()
		if _, err := os.Stat(repo.Path); os.IsNotExist(err) {
			// err := os.MkdirAll(filepath.Dir(repoPath), 0755)
			// if err != nil {
			// 	return fmt.Errorf("failed to create directory for repository: %w", err)
			// }

			err = git.Clone(git.CloneArgs{
				URL:  repo.URL,
				Path: repo.Path,
				Auth: &repo.Auth,
			})
			if err != nil {
				return fmt.Errorf("failed to clone repository: %w", err)
			}
		}
	}
}
