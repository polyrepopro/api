package repositories

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
	localgit "github.com/polyrepopro/api/git"
)

type progressReporter struct{}

func (pr *progressReporter) Write(p []byte) (n int, err error) {
	multilog.Debug("repositories", "update", map[string]interface{}{
		"message": string(p),
	})
	return len(p), nil
}

// Update updates a repository by fetching all remotes and pulling the latest changes.
// It also prunes all tags and branches that are no longer present.
//
// Arguments:
//   - workspace: The workspace to update.
//   - repo: The repository to update.
//
// Returns:
//   - error: An error if something went wrong.
func Update(workspace *config.Workspace, repo *config.Repository) error {
	auth := localgit.GetAuth(repo.URL, repo.Auth)
	repoPath := fmt.Sprintf("%s/%s", workspace.Path, repo.Path)

	// Open the repository.
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	fetchOpts := &git.FetchOptions{
		RemoteName: "origin",
		Tags:       git.AllTags,
		Prune:      true,
		Progress:   &progressReporter{},
	}

	if auth.Name() != "" {
		fetchOpts.Auth = auth
	}

	// Fetch all remotes
	err = r.Fetch(fetchOpts)
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to fetch: %w", err)
	}

	// Get the working directory.
	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	pullOpts := &git.PullOptions{
		RemoteName: "origin",
		Progress:   &progressReporter{},
	}

	if auth.Name() != "" {
		pullOpts.Auth = auth
	}

	// Check for unstaged changes
	status, err := w.Status()
	if err != nil {
		return fmt.Errorf("failed to get worktree status: %w", err)
	}

	if status.IsClean() {
		// Pull the latest changes.
		err = w.Pull(pullOpts)
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("failed to pull: %w", err)
		}

	} else {
		multilog.Warn("repositories.update", "repository has unstaged changes", map[string]interface{}{
			"path": repoPath,
		})
	}

	return nil
}
