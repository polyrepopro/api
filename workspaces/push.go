package workspaces

import (
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/repositories"
)

type PushArgs struct {
	Workspace *config.Workspace
}

func Push(args PushArgs) []error {
	var errors []error

	for _, repo := range *args.Workspace.Repositories {
		err := repositories.Push(repositories.PushArgs{
			Workspace:  args.Workspace,
			Repository: &repo,
		})
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
