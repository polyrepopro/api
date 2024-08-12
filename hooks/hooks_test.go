package hooks

import (
	"context"
	"testing"

	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/test"
)

func TestRun(t *testing.T) {
	test.Setup() // Configures the test environment (like logging).
	ctx := context.Background()
	hook := &config.Hook{
		Type: config.CloneHook,
		Commands: []config.Command{
			{
				Name: "ls",
				Cwd:  ".",
				Command: []string{
					"ls",
					"-l",
					".",
				},
			},
			{
				Name: "sleep",
				Command: []string{
					"sleep",
					"5",
				},
			},
			{
				Name: "echo",
				Command: []string{
					"echo",
					"hello",
				},
			},
		},
	}
	err := Run(ctx, hook)
	if err != nil {
		t.Fatal(err)
	}
}
