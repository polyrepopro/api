package workspaces

import (
	"fmt"
	"os"

	"github.com/mateothegreat/go-util/files"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
	"github.com/polyrepopro/api/repositories"
)

type SyncArgs struct {
	config.DefaultArgs
	Name string
}

func SyncAll(args *SyncArgs) ([]string, []error) {
	cfg, err := config.GetRelativeConfig()
	if err != nil {
		return nil, []error{err}
	}

	var msgs []string
	for _, workspace := range *cfg.Workspaces {
		syncArgs := SyncArgs{
			Name: workspace.Name,
		}
		if args != nil {
			syncArgs.DefaultArgs = args.DefaultArgs
		}

		m, err := Sync(syncArgs)
		if err != nil {
			return nil, []error{err}
		}

		msgs = append(msgs, m...)
	}

	return msgs, nil
}

func Sync(args SyncArgs) ([]string, error) {
	var ret []string

	cfg, err := config.GetRelativeConfig()
	if err != nil {
		return nil, err
	}

	workspace, err := cfg.GetWorkspace(args.Name)
	if err != nil {
		return nil, err
	}

	workspacePath := files.ExpandPath(workspace.Path)

	if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
		err = os.MkdirAll(workspacePath, 0755)
		if err != nil {
			return nil, err
		}

		ret = append(ret, fmt.Sprintf("created workspace directory %s", workspacePath))
	}

	for _, repo := range *workspace.Repositories {
		repoPath := fmt.Sprintf("%s/%s", workspacePath, repo.Path)

		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			err = git.Clone(git.CloneArgs{
				URL:  repo.URL,
				Path: repoPath,
				Auth: repo.Auth,
			})
			if err != nil {
				return nil, err
			}

			ret = append(ret, fmt.Sprintf("cloned new repository %s", repo.URL))
		}

		err := repositories.Update(workspace, &repo)
		if err != nil {
			return nil, err
		}

		ret = append(ret, fmt.Sprintf("updated repository %s", repo.URL))
	}

	return ret, nil
}
