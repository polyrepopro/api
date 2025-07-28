package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/utils"
)

// PushArgs represents the arguments for pushing changes to a remote repository.
type PushArgs struct {
	URL    string
	Remote string
	Path   string
	Auth   *config.Auth
}

// pushProgress represents the progress of a push operation.
type pushProgress struct{}

// Write writes the progress of a push operation.
//
// Arguments:
// - p: the progress of the push operation
//
// Returns:
// - n: the number of bytes written
func (h *pushProgress) Write(p []byte) (n int, err error) {
	multilog.Debug("git.push", "pushing progress", map[string]interface{}{
		"message": string(p),
	})
	return len(p), nil
}

// Push pushes the changes to the remote repository.
//
// Arguments:
// - args: the push arguments including path, remote, url, and auth
//
// Returns:
// - error: any error encountered during the push process
func Push(args PushArgs) error {
	expandedPath, err := utils.ExpandPath(args.Path)
	if err != nil {
		return fmt.Errorf("failed to expand path %q: %w", args.Path, err)
	}

	repo, err := git.PlainOpen(expandedPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w for repo %q", err, expandedPath)
	}

	opts := &git.PushOptions{
		RemoteName: args.Remote,
		Progress:   &pushProgress{},
	}

	// Get the actual remote URL from the repository, not from config
	remotes, err := repo.Remotes()
	if err != nil {
		return fmt.Errorf("failed to get remotes: %w", err)
	}

	var actualRemoteURL string
	for _, remote := range remotes {
		if remote.Config().Name == args.Remote {
			if len(remote.Config().URLs) > 0 {
				actualRemoteURL = remote.Config().URLs[0]
			}
			break
		}
	}

	if actualRemoteURL == "" {
		return fmt.Errorf("remote %q not found in repository", args.Remote)
	}

	auth := GetAuth(actualRemoteURL, args.Auth)
	if auth != nil {
		opts.Auth = auth
		multilog.Debug("git.push", "using auth", map[string]interface{}{
			"auth_name": auth.Name(),
			"auth_type": fmt.Sprintf("%T", auth),
			"for_url":   actualRemoteURL,
		})
	}

	err = repo.Push(opts)
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			return nil
		}

		multilog.Debug("git.push", "push failed", map[string]interface{}{
			"error":       err.Error(),
			"remote_name": opts.RemoteName,
			"path":        expandedPath,
		})

		return fmt.Errorf("failed to push changes: %w for repo %q", err, expandedPath)
	}

	return nil
}
