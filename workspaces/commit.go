package workspaces

import (
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
	"github.com/polyrepopro/api/repositories"
)

type CommitArgs struct {
	Workspace *config.Workspace
	Message   string
}

// Commit commits the changes for each repository in the workspace.
//
// Arguments:
//   - args: The arguments for the commit.
//
// Returns:
//   - []git.CommitResult: The results of the commit.
func Commit(args CommitArgs) ([]git.CommitResult, []error) {
	var errors []error
	var results []git.CommitResult

	for _, repo := range *args.Workspace.Repositories {
		res, err := repositories.Commit(repositories.CommitArgs{
			Workspace:  args.Workspace,
			Repository: &repo,
			Message:    args.Message,
		})
		if err != nil {
			errors = append(errors, err)
		}
		if res != nil {
			results = append(results, *res)
		}
	}

	return results, errors
}
