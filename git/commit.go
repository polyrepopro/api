package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/mateothegreat/go-util/files"
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

// GetGitUser retrieves the git user configuration from the system git config.
//
// Arguments:
// - path: the repository path to check for local git config first
//
// Returns:
// - *object.Signature: the signature containing user name, email, and timestamp
func GetGitUser(path string) (*object.Signature, error) {
	var name, email string
	var err error

	// Try to get user.name from git config
	cmd := exec.Command("git", "-C", path, "config", "user.name")
	if output, cmdErr := cmd.Output(); cmdErr == nil {
		name = strings.TrimSpace(string(output))
	}

	// If no local config, try global
	if name == "" {
		cmd = exec.Command("git", "config", "--global", "user.name")
		if output, cmdErr := cmd.Output(); cmdErr == nil {
			name = strings.TrimSpace(string(output))
		}
	}

	// Try to get user.email from git config
	cmd = exec.Command("git", "-C", path, "config", "user.email")
	if output, cmdErr := cmd.Output(); cmdErr == nil {
		email = strings.TrimSpace(string(output))
	}

	// If no local config, try global
	if email == "" {
		cmd = exec.Command("git", "config", "--global", "user.email")
		if output, cmdErr := cmd.Output(); cmdErr == nil {
			email = strings.TrimSpace(string(output))
		}
	}

	// Fallback to system defaults if git config is not available
	if name == "" {
		if envUser := os.Getenv("USER"); envUser != "" {
			name = envUser
		} else {
			name = "polyrepo-user"
		}
	}

	if email == "" {
		if hostname, err := os.Hostname(); err == nil {
			email = fmt.Sprintf("%s@%s", name, hostname)
		} else {
			email = fmt.Sprintf("%s@localhost", name)
		}
	}

	return &object.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}, err
}

// AddGitignoreToWorktree parses and adds gitignore patterns to the worktree excludes.
//
// Arguments:
// - wt: the git worktree instance
// - path: the repository path containing the .gitignore file
//
// Returns:
// - []gitignore.Pattern: the parsed gitignore patterns
// - error: any error encountered while reading or parsing the .gitignore file
func AddGitignoreToWorktree(paths ...string) ([]gitignore.Pattern, error) {
	patterns := make([]gitignore.Pattern, 0)

	for _, path := range paths {
		if !files.FileExists(filepath.Join(path, ".gitignore")) {
			return nil, nil
		}

		f, err := os.Open(filepath.Join(path, ".gitignore"))
		if err != nil {
			return nil, fmt.Errorf("failed to read .gitignore: %w", err)
		}
		defer f.Close()

		fileScanner := bufio.NewScanner(f)
		fileScanner.Split(bufio.ScanLines)

		for fileScanner.Scan() {
			patterns = append(patterns, gitignore.ParsePattern(fileScanner.Text(), nil))
		}

	}

	return patterns, nil
}

// Commit creates a git commit with the provided message and current staged changes.
//
// Arguments:
// - args: the commit arguments including path, auth, and message
//
// Returns:
// - *CommitResult: the result containing commit hash and changed files
// - error: any error encountered during the commit process
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

	patterns, err := AddGitignoreToWorktree(args.Path)
	if err != nil {
		return result, fmt.Errorf("failed to add gitignore to worktree: %w", err)
	}

	worktree.Excludes = append(worktree.Excludes, patterns...)

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

	// Get the git user information for the commit signature
	signature, err := GetGitUser(args.Path)
	if err != nil {
		return result, fmt.Errorf("failed to get git user information: %w", err)
	}

	hash, err := worktree.Commit(args.Message, &git.CommitOptions{
		All:       false,
		Author:    signature,
		Committer: signature,
	})
	if err != nil {
		return result, fmt.Errorf("failed to commit changes: %w", err)
	}
	result.Hash = hash.String()

	return result, nil
}
