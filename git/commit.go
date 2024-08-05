package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/polyrepopro/api/config"
)

type CommitArgs struct {
	Path    string
	Auth    *config.Auth
	Message string
}

func Commit(args CommitArgs) (*plumbing.Hash, error) {
	repo, err := git.PlainOpen(args.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree status: %w", err)
	}

	if status.IsClean() {
		return nil, nil // No changes to commit
	}

	_, err = worktree.Add(".")
	if err != nil {
		return nil, fmt.Errorf("failed to add changes: %w", err)
	}

	_, err = worktree.Commit(args.Message, &git.CommitOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to commit changes: %w", err)
	}

	return nil, nil
}
