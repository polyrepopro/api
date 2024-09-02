package git

import "github.com/go-git/go-git/v5"

func Status(path string) (git.Status, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return git.Status{}, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return git.Status{}, err
	}

	status, err := worktree.Status()
	if err != nil {
		return git.Status{}, err
	}

	return status, nil
}
