package git

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/utils"
)

// BranchStatus represents the relationship between local and remote branches
type BranchStatus struct {
	Behind int  // Number of commits behind remote
	Ahead  int  // Number of commits ahead of remote  
	NeedsPush bool // Local has commits that need to be pushed
	NeedsPull bool // Remote has commits that need to be pulled
}

// EnhancedStatus includes both working directory status and branch comparison
type EnhancedStatus struct {
	WorkingTree git.Status
	Branch      BranchStatus
	HasChanges  bool // Local working directory has uncommitted changes
}

// Status retrieves the git status for the repository at the specified path.
//
// Arguments:
// - path: the file system path to the git repository
//
// Returns:
// - git.Status: the status of the repository
// - error: any error encountered while getting the status
func Status(path string) (git.Status, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return git.Status{}, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return git.Status{}, err
	}

	// Force refresh the worktree status by checking filesystem
	status, err := worktree.Status()
	if err != nil {
		return git.Status{}, err
	}

	// Verification: check with git command if the library reports dirty but git says clean
	if len(status) > 0 {
		cmd := exec.Command("git", "-C", path, "status", "--porcelain")
		output, gitErr := cmd.Output()
		if gitErr == nil && len(strings.TrimSpace(string(output))) == 0 {
			// Git command says clean, return empty status
			return git.Status{}, nil
		}

		multilog.Warn("git.status", "resorted to git command directly", map[string]interface{}{
			"path":    path,
			"status":  status,
			"output":  string(output),
			"gitErr":  gitErr,
			"command": cmd.String(),
		})
	}

	return status, nil
}

// EnhancedStatusWithRemote retrieves comprehensive status including remote branch comparison.
//
// Arguments:
// - path: the file system path to the git repository
// - remoteName: the name of the remote to compare against (e.g., "origin", "upstream")
//
// Returns:
// - EnhancedStatus: comprehensive status including working tree and branch comparison
// - error: any error encountered while getting the status
func EnhancedStatusWithRemote(path, remoteName string) (EnhancedStatus, error) {
	expandedPath, err := utils.ExpandPath(path)
	if err != nil {
		return EnhancedStatus{}, err
	}

	repo, err := git.PlainOpen(expandedPath)
	if err != nil {
		return EnhancedStatus{}, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return EnhancedStatus{}, err
	}

	// Get working tree status
	workingStatus, err := worktree.Status()
	if err != nil {
		return EnhancedStatus{}, err
	}

	// Get current branch
	head, err := repo.Head()
	if err != nil {
		return EnhancedStatus{
			WorkingTree: workingStatus,
			HasChanges:  len(workingStatus) > 0,
		}, nil // Return partial status if we can't get branch info
	}

	// Get branch status
	branchStatus, err := getBranchStatus(repo, head, remoteName)
	if err != nil {
		multilog.Debug("git.status", "failed to get branch status", map[string]interface{}{
			"path":   path,
			"remote": remoteName,
			"error":  err.Error(),
		})
		// Return partial status without branch comparison
		branchStatus = BranchStatus{}
	}

	return EnhancedStatus{
		WorkingTree: workingStatus,
		Branch:      branchStatus,
		HasChanges:  len(workingStatus) > 0,
	}, nil
}

// getBranchStatus compares local branch with remote branch
func getBranchStatus(repo *git.Repository, head *plumbing.Reference, remoteName string) (BranchStatus, error) {
	if !head.Name().IsBranch() {
		return BranchStatus{}, nil // Not on a branch (detached HEAD)
	}

	branchName := head.Name().Short()
	localCommit := head.Hash()

	// Get remote reference
	remoteRefName := plumbing.NewRemoteReferenceName(remoteName, branchName)
	remoteRef, err := repo.Reference(remoteRefName, true)
	if err != nil {
		// Remote branch doesn't exist, assume all local commits need push
		localCommits, err := countCommitsFromRef(repo, head.Hash())
		if err != nil {
			return BranchStatus{}, err
		}
		return BranchStatus{
			Ahead:     localCommits,
			Behind:    0,
			NeedsPush: localCommits > 0,
			NeedsPull: false,
		}, nil
	}

	remoteCommit := remoteRef.Hash()

	// If commits are the same, branches are in sync
	if localCommit == remoteCommit {
		return BranchStatus{
			Ahead:     0,
			Behind:    0,
			NeedsPush: false,
			NeedsPull: false,
		}, nil
	}

	// Count commits ahead and behind
	ahead, err := countCommitsBetween(repo, remoteCommit, localCommit)
	if err != nil {
		return BranchStatus{}, err
	}

	behind, err := countCommitsBetween(repo, localCommit, remoteCommit)
	if err != nil {
		return BranchStatus{}, err
	}

	return BranchStatus{
		Ahead:     ahead,
		Behind:    behind,
		NeedsPush: ahead > 0,
		NeedsPull: behind > 0,
	}, nil
}

// countCommitsBetween counts commits that are in 'to' but not in 'from' (equivalent to git rev-list from..to)
func countCommitsBetween(repo *git.Repository, from, to plumbing.Hash) (int, error) {
	// Use git rev-list command directly for accuracy - equivalent to "git rev-list from..to --count"
	worktree, err := repo.Worktree()
	if err != nil {
		return 0, err
	}
	
	// Use git command to count commits between references
	cmd := exec.Command("git", "-C", worktree.Filesystem.Root(), "rev-list", "--count", fmt.Sprintf("%s..%s", from.String(), to.String()))
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	
	countStr := strings.TrimSpace(string(output))
	if countStr == "" {
		return 0, nil
	}
	
	count := 0
	_, err = fmt.Sscanf(countStr, "%d", &count)
	return count, err
}

// countCommitsFromRef counts total commits from a reference
func countCommitsFromRef(repo *git.Repository, hash plumbing.Hash) (int, error) {
	worktree, err := repo.Worktree()
	if err != nil {
		return 0, err
	}
	
	// Use git command to count total commits from reference
	cmd := exec.Command("git", "-C", worktree.Filesystem.Root(), "rev-list", "--count", hash.String())
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	
	countStr := strings.TrimSpace(string(output))
	if countStr == "" {
		return 0, nil
	}
	
	count := 0
	_, err = fmt.Sscanf(countStr, "%d", &count)
	return count, err
}
