package workspaces

import (
	"sync"

	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
	"github.com/polyrepopro/api/repositories"
)

type PushArgs struct {
	Workspace *config.Workspace
}

func Push(args PushArgs) []error {
	var errors []error

	var wg sync.WaitGroup
	for _, repo := range *args.Workspace.Repositories {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repositories.Push(repositories.PushArgs{
				PushArgs: git.PushArgs{
					Remote: repo.Origin,
				},
				Workspace:  args.Workspace,
				Repository: &repo,
			})
			if err != nil {
				errors = append(errors, err)
			}
		}()
	}
	wg.Wait()
	return errors
}
