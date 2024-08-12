package commands

import (
	"testing"

	"github.com/polyrepopro/api/commands"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/test"
)

func TestWatch(t *testing.T) {
	test.Setup()
	watch := config.Watch{
		Cwd: "~/workspace/polyrepo/examples/example-test-repo",
		Paths: []string{
			"**/*.go",
		},
		Commands: []commands.Command{
			{
				Name:    "run",
				Cwd:     "test",
				Command: []string{"go", "run", "main.go"},
			},
		},
	}
	Watch(watch)
}
