package git

import (
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/mateothegreat/go-multilog/multilog"
)

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
