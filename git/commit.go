package git

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
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

func AddGitignoreToWorktree(wt *git.Worktree, path string) ([]gitignore.Pattern, error) {
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

	var patterns []gitignore.Pattern
	for fileScanner.Scan() {
		patterns = append(patterns, gitignore.ParsePattern(fileScanner.Text(), nil))
	}

	return patterns, nil
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

	patterns, err := AddGitignoreToWorktree(worktree, args.Path)
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

	hash, err := worktree.Commit(args.Message, &git.CommitOptions{
		All: false,
	})
	if err != nil {
		return result, fmt.Errorf("failed to commit changes: %w", err)
	}
	result.Hash = hash.String()

	return result, nil
}
