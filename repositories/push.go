// Package repositories provides functions for interacting with repositories.
package repositories

import (
	"fmt"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
)

// PushArgs represents the arguments for pushing changes to a remote repository.
type PushArgs struct {
	git.PushArgs
	Workspace  *config.Workspace
	Repository *config.Repository
}

// Push pushes the changes to the remote repository.
//
// Arguments:
// - args: the push arguments including workspace, repository, and remote
//
// Returns:
// - error: any error encountered during the push process
func Push(args PushArgs) error {
	r := args.Remote
	if args.Remote == "" {
		remote, err := GetDefaultRemote(GetRemotesArgs{
			Workspace:  args.Workspace,
			Repository: args.Repository,
		})
		if err != nil {
			return fmt.Errorf("failed to get default remote: %w", err)
		}
		r = remote.Name
	}
	multilog.Debug(args.Workspace.Name, args.Repository.Name, map[string]interface{}{
		"remote": r,
		"url":    args.Repository.URL,
	})
	err := git.Push(git.PushArgs{
		Path:   fmt.Sprintf("%s/%s", args.Workspace.Path, args.Repository.Path),
		Remote: r,
		URL:    args.Repository.URL,
		Auth:   args.Repository.Auth,
	})
	if err != nil {
		return fmt.Errorf("failed to push remote %q: %w", r, err)
	}
	return nil
}
