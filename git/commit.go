package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/polyrepopro/api/config"
)

type CommitArgs struct {
	Path    string
	Auth    *config.Auth
	Message string
}

type CommitResult struct {
	Path     string
	Hash     string
	Messages *[]string
}

type CommitResultMessage struct {
	Name   string
	Status git.StatusCode
}

func Commit(args CommitArgs) (*CommitResult, error) {
	result := &CommitResult{
		Path:     args.Path,
		Messages: &[]string{},
	}

	repo, err := git.PlainOpen(args.Path)
	if err != nil {
		return result, fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return result, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return result, fmt.Errorf("failed to get worktree status: %w", err)
	}

	for name, change := range status {
		if change.Worktree == git.Unmodified {
			continue
		}
		*result.Messages = append(*result.Messages, fmt.Sprintf("%s (%s)", name, string(change.Worktree)))
	}

	_, err = worktree.Add(".")
	if err != nil {
		return result, fmt.Errorf("failed to add changes: %w", err)
	}

	hash, err := worktree.Commit(args.Message, &git.CommitOptions{
		All: true,
	})
	if err != nil {
		return result, fmt.Errorf("failed to commit changes: %w", err)
	}
	result.Hash = hash.String()

	return result, nil
}
