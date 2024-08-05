package workspaces

import (
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/repositories"
)

type PullArgs struct {
	Workspace *config.Workspace
}

func Pull(args PullArgs) []error {
	var errors []error

	for _, repo := range *args.Workspace.Repositories {
		err := repositories.Pull(repositories.PullArgs{
			Workspace:  args.Workspace,
			Repository: &repo,
		})
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
