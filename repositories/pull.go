package repositories

import (
	"fmt"
	"os"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
)

type PullArgs struct {
	git.PullArgs
	Workspace  *config.Workspace
	Repository *config.Repository
}

func Pull(args PullArgs) error {
	if args.Remote == "" {
		args.Remote = "origin"
	}

	multilog.Debug("repositories.pull", "pulling repository", map[string]interface{}{
		"repository": args.Repository,
	})

	stat, err := os.Stat(fmt.Sprintf("%s/%s/.git", args.Workspace.GetAbsolutePath(), args.Repository.Path))
	if os.IsNotExist(err) || !stat.IsDir() {
		multilog.Info("repositories.pull", "repository not found, cloning", map[string]interface{}{
			"path": fmt.Sprintf("%s/%s", args.Workspace.GetAbsolutePath(), args.Repository.Path),
			"url":  args.Repository.URL,
		})

		err = git.Clone(git.CloneArgs{
			URL:  args.Repository.URL,
			Path: fmt.Sprintf("%s/%s", args.Workspace.GetAbsolutePath(), args.Repository.Path),
		})
		if err != nil {
			multilog.Fatal("repositories.pull", "failed to clone repository", map[string]interface{}{
				"repository": args,
			})
			return err
		}

		multilog.Info("repositories.pull", "✅ cloned repository", map[string]interface{}{
			"repository": args,
		})
	} else {
		err = git.Pull(git.PullArgs{
			Path:   fmt.Sprintf("%s/%s", args.Workspace.Path, args.Repository.Path),
			Remote: args.Remote,
			URL:    args.Repository.URL,
			Auth:   args.Repository.Auth,
		})
		if err != nil {
			return fmt.Errorf("failed to pull remote %q: %w", args.Remote, err)
		}
	}
	return nil
}
