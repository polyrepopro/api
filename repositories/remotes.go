package repositories

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/utils"
)

// GetRemotesArgs represents the arguments for getting the remotes of a repository.
type GetRemotesArgs struct {
	Workspace  *config.Workspace
	Repository *config.Repository
}

// Remote represents a remote repository.
type Remote struct {
	Name string
	URLs []string
}

// GetRemotes gets the remotes for a repository.
//
// Arguments:
// - args: the remotes arguments including workspace and repository
//
// Returns:
// - []Remote: the remotes for the repository
func GetRemotes(args GetRemotesArgs) ([]Remote, error) {
	result := make([]Remote, 0)

	repoPath := fmt.Sprintf("%s/%s", args.Workspace.Path, args.Repository.Path)
	expandedPath, err := utils.ExpandPath(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to expand path %q: %w", repoPath, err)
	}

	repo, err := git.PlainOpen(expandedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return nil, fmt.Errorf("failed to get remotes: %w", err)
	}

	for _, remote := range remotes {
		result = append(result, Remote{
			Name: remote.Config().Name,
			URLs: remote.Config().URLs,
		})
	}

	return result, nil
}

// GetDefaultRemote gets the default remote for a repository as if it were set with `git push -u <remote> <main>`.
// This follows Git's logic for determining the upstream remote:
// 1. If HEAD points to a branch with an upstream, use that remote
// 2. If only one remote exists, use it
// 3. If "origin" exists among multiple remotes, prefer it
// 4. Otherwise, return the first remote
//
// Arguments:
// - args: the remotes arguments including workspace and repository
//
// Returns:
// - Remote: the default remote for the repository
func GetDefaultRemote(args GetRemotesArgs) (Remote, error) {
	repoPath := fmt.Sprintf("%s/%s", args.Workspace.Path, args.Repository.Path)
	expandedPath, err := utils.ExpandPath(repoPath)
	if err != nil {
		return Remote{}, fmt.Errorf("failed to expand path %q: %w", repoPath, err)
	}

	repo, err := git.PlainOpen(expandedPath)
	if err != nil {
		return Remote{}, fmt.Errorf("failed to open repository: %w", err)
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return Remote{}, fmt.Errorf("failed to get remotes: %w", err)
	}

	if len(remotes) == 0 {
		return Remote{}, fmt.Errorf("no remotes found")
	}

	// If only one remote, return it
	if len(remotes) == 1 {
		remote := remotes[0]
		return Remote{
			Name: remote.Config().Name,
			URLs: remote.Config().URLs,
		}, nil
	}

	// Check if HEAD points to a branch with upstream tracking
	head, err := repo.Head()
	if err == nil && head.Name().IsBranch() {
		cfg, err := repo.Config()
		if err == nil {
			branchName := head.Name().Short()
			if branch, exists := cfg.Branches[branchName]; exists && branch.Remote != "" {
				// Find the upstream remote
				for _, remote := range remotes {
					if remote.Config().Name == branch.Remote {
						return Remote{
							Name: remote.Config().Name,
							URLs: remote.Config().URLs,
						}, nil
					}
				}
			}
		}
	}

	// Look for "origin" remote (Git's default preference)
	for _, remote := range remotes {
		if remote.Config().Name == "origin" {
			return Remote{
				Name: remote.Config().Name,
				URLs: remote.Config().URLs,
			}, nil
		}
	}

	// Return first remote as fallback
	remote := remotes[0]
	return Remote{
		Name: remote.Config().Name,
		URLs: remote.Config().URLs,
	}, nil
}
