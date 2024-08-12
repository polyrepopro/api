package workspaces

import (
	"fmt"

	"github.com/mateothegreat/go-util/files"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
)

type SwitchArgs struct {
	Workspace *config.Workspace
	Branch    string
}

func Switch(args SwitchArgs) []error {
	var errors []error

	for _, repo := range *args.Workspace.Repositories {
		err := git.Switch(&git.SwitchArgs{
			Path:   fmt.Sprintf("%s/%s", files.ExpandPath(args.Workspace.Path), repo.Path),
			Branch: args.Branch,
		})
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
