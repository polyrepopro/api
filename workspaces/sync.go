package workspaces

import "github.com/polyrepopro/api/config"

type SyncArgs struct {
}

func Sync(workspace *config.Workspace, args SyncArgs) error {
	for _, repo := range *workspace.Repositories {
	}
}
