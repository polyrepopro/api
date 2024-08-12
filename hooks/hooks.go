package hooks

import (
	"context"

	"github.com/polyrepopro/api/commands"
	"github.com/polyrepopro/api/config"
)

func Run(ctx context.Context, hook *config.Hook) error {
	for _, command := range hook.Commands {
		if err := commands.Run(ctx, command, command.Cwd); err != nil {
			return err
		}
	}
	return nil
}
