package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type SwitchArgs struct {
	Path   string
	Branch string
}

func Switch(args *SwitchArgs) error {
	repo, err := git.PlainOpen(args.Path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(args.Branch),
	})
	if err != nil {
		return fmt.Errorf("failed to checkout branch %q: %w", args.Branch, err)
	}

	return nil
}
