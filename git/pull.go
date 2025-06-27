package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
)

type PullArgs struct {
	URL    string
	Remote string
	Path   string
	Auth   *config.Auth
}

type pullProgress struct{}

func (h *pullProgress) Write(p []byte) (n int, err error) {
	multilog.Debug("git.pull", "pulling progress", map[string]interface{}{
		"message": string(p),
	})
	return len(p), nil
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
		RemoteName:        args.Remote,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		// Progress:          &pullProgress{},
		Force: true, // Allow non-fast-forward updates
	}

	multilog.Debug("git.pull", "pulling", map[string]interface{}{
		"url":    args.URL,
		"remote": args.Remote,
		"path":   args.Path,
	})

	auth := GetAuth(args.URL, args.Auth)
	if auth.Name() != "" {
		opts.Auth = auth
	}

	err = worktree.Pull(opts)
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to pull changes: %w for %q", err, args.Path)
	}

	return nil
}
