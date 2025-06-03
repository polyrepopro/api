package repositories

import (
	"fmt"

	"github.com/mateothegreat/go-util/files"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
)

type CommitArgs struct {
	git.CommitArgs
	Workspace  *config.Workspace
	Repository *config.Repository
	Message    string
}

func Commit(args CommitArgs) (*git.CommitResult, error) {
	result, err := git.Commit(git.CommitArgs{
		Path:    files.ExpandPath(fmt.Sprintf("%s/%s", args.Workspace.Path, args.Repository.Path)),
		Auth:    args.Repository.Auth,
		Message: args.Message,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to commit changes: %w", err)
	}

	return result, nil
}
