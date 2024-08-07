package workspaces

import (
	"fmt"
	"time"

	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/repositories"
)

type CommitArgs struct {
	Workspace *config.Workspace
}

func Commit(args CommitArgs) ([]repositories.CommitResult, []error) {
	var errors []error
	var results []repositories.CommitResult

	for _, repo := range *args.Workspace.Repositories {
		res, err := repositories.Commit(repositories.CommitArgs{
			Workspace:  args.Workspace,
			Repository: &repo,
			Message:    fmt.Sprintf("Test commit @ %s", time.Now()),
		})
		if err != nil {
			errors = append(errors, err)
		}
		results = append(results, res)
	}

	return results, errors
}
