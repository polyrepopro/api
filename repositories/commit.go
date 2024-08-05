package repositories

import (
	"fmt"

	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
)

type CommitArgs struct {
	git.CommitArgs
	Workspace  *config.Workspace
	Repository *config.Repository
	Message    string
}

type CommitResult struct {
	Hash string
}

func Commit(args CommitArgs) (CommitResult, error) {
	hash, err := git.Commit(git.CommitArgs{
		Path: fmt.Sprintf("%s/%s", args.Workspace.Path, args.Repository.Path),
		Auth: args.Repository.Auth,
	})
	if err != nil {
		return CommitResult{}, fmt.Errorf("failed to commit changes: %w", err)
	}

	return CommitResult{
		Hash: func() string {
			if hash == nil {
				return ""
			}
			return hash.String()
		}(),
	}, nil
}
