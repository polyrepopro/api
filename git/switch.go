package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/polyrepopro/api/config"
)

type SwitchArgs struct {
	Repository *config.Repository
	Branch     string
}

func Switch(args *SwitchArgs) error {
	repo, err := git.PlainOpen(args.Repository.Path)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(args.Branch),
	})
	if err != nil {
		return err
	}

	return nil
}
