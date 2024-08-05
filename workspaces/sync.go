package workspaces

import (
	"fmt"
	"os"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
	"github.com/polyrepopro/api/util"
)

type SyncArgs struct {
	Name string
}

func Sync(args SyncArgs) error {
	cfg, err := config.GetRelativeConfig()
	if err != nil {
		return fmt.Errorf("failed to get relative config: %w", err)
	}

	for _, workspace := range *cfg.Workspaces {
		multilog.Info("workspace.sync", "syncing workspace", map[string]interface{}{
			"name": workspace.Name,
		})
		workspacePath := util.ExpandPath(workspace.Path)
		if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
			err = os.MkdirAll(workspacePath, 0755)
			if err != nil {
				multilog.Fatal("workspace.sync", "failed to create workspace directory", map[string]interface{}{
					"path":  workspacePath,
					"error": err,
				})
			}
			multilog.Info("workspace.sync", "created workspace directory", map[string]interface{}{
				"path": workspacePath,
			})
		}
		for _, repo := range *workspace.Repositories {
			multilog.Info("workspace.sync", "syncing repository", map[string]interface{}{
				"workspace": workspace.Name,
				"path":      util.ExpandPath(repo.Path),
			})
			if _, err := os.Stat(util.ExpandPath(repo.Path)); os.IsNotExist(err) {
				err = git.Clone(git.CloneArgs{
					URL:  repo.URL,
					Path: util.ExpandPath(repo.Path),
					Auth: repo.Auth,
				})
				if err != nil {
					multilog.Fatal("workspace.sync", "failed to clone repository", map[string]interface{}{
						"repository": repo.URL,
						"error":      err,
					})
				}
				multilog.Info("workspace.sync", "cloned repository", map[string]interface{}{
					"repository": repo.URL,
				})
			}
		}
	}

	return nil
}
