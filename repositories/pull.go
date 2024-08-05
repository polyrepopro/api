package repositories

import (
	"fmt"

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

	err := git.Pull(git.PullArgs{
		Path:   fmt.Sprintf("%s/%s", args.Workspace.Path, args.Repository.Path),
		Remote: args.Remote,
		URL:    args.Repository.URL,
		Auth:   args.Repository.Auth,
	})
	if err != nil {
		return fmt.Errorf("failed to pull remote %q: %w", args.Remote, err)
	}
	return nil
}
