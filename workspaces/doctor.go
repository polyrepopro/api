package workspaces

import (
	"os"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
)

func Doctor(name string) error {
	config, err := config.GetConfig()
	if err != nil {
		multilog.Fatal("workspaces.doctor", "failed to get config", map[string]interface{}{
			"error": err,
		})
	}
	// Ensure that the workspace path exists.
	workspace, err := config.GetWorkspaceByWorkingDir()
	if err != nil {
		multilog.Fatal("workspaces.doctor", "failed to get workspace", map[string]interface{}{
			"error": err,
		})
	}

	path := workspace.GetAbsolutePath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		multilog.Error("workspaces.doctor", "workspace directory does not exist", map[string]interface{}{
			"name": name,
		})

		err := os.MkdirAll(path, 0755)
		if err != nil {
			multilog.Fatal("workspaces.doctor", "failed to create workspace directory", map[string]interface{}{
				"error": err,
			})
		}
	}

	return nil
}
