package repositories

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
)

func TestDoctor(t *testing.T) {
	config, err := config.GetConfig()
	if err != nil {
		multilog.Fatal("workspaces.doctor", "failed to get config", map[string]interface{}{
			"error": err,
		})
	}

	workspace, err := config.GetWorkspaceByWorkingDir()
	if err != nil {
		multilog.Fatal("workspaces.doctor", "failed to get workspace", map[string]interface{}{
			"error": err,
		})
	}

	for _, repository := range workspace.Repositories {
		err := Doctor(&repository)
		assert.NoError(t, err)
	}
}
