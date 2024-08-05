package workspaces

import (
	"fmt"
	"os"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
	"github.com/polyrepopro/api/repositories"
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
			repoPath := fmt.Sprintf("%s/%s", workspacePath, repo.Path)
			multilog.Info("workspace.sync", "syncing repository", map[string]interface{}{
				"workspace": workspace.Name,
				"path":      repoPath,
			})
			if _, err := os.Stat(repoPath); os.IsNotExist(err) {
				err = git.Clone(git.CloneArgs{
					URL:  repo.URL,
					Path: repoPath,
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
			err = repositories.Update(&workspace, &repo)
			if err != nil {
				multilog.Fatal("workspace.sync", "failed to update repository", map[string]interface{}{
					"repository": repo.URL,
					"error":      err,
				})
			}
			multilog.Info("workspace.sync", "updated repository", map[string]interface{}{
				"repository": repo.URL,
			})
		}
	}

	return nil
}
