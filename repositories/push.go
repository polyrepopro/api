package repositories

import (
	"fmt"

	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
)

type PushArgs struct {
	git.PushArgs
	Workspace  *config.Workspace
	Repository *config.Repository
}

func Push(args PushArgs) error {
	if args.Remote == "" {
		args.Remote = "origin"
	}

	err := git.Push(git.PushArgs{
		Path:   fmt.Sprintf("%s/%s", args.Workspace.Path, args.Repository.Path),
		Remote: args.Remote,
		URL:    args.Repository.URL,
		Auth:   args.Repository.Auth,
	})
	if err != nil {
		return fmt.Errorf("failed to push remote %q: %w", args.Remote, err)
	}
	return nil
}
