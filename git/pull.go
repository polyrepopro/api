package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/polyrepopro/api/config"
)

type PullArgs struct {
	URL    string
	Remote string
	Path   string
	Auth   *config.Auth
}

func Pull(args PullArgs) error {
	repo, err := git.PlainOpen(args.Path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	opts := &git.PullOptions{
		RemoteName: args.Remote,
	}

	auth := GetAuth(args.URL, args.Auth)
	if auth.Name() != "" {
		opts.Auth = auth
	}

	err = worktree.Pull(opts)
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to pull changes: %w", err)
	}

	return nil

}
