package repositories

import (
	"fmt"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/git"
)

func Add(args config.Repository) error {
	config, err := config.GetRelativeConfig()
	if err != nil {
		return err
	}

	workspace, err := config.GetWorkspaceByWorkingDir()
	if err != nil {
		return err
	}

	*workspace.Repositories = append(*workspace.Repositories, args)

	err = config.SaveConfig()
	if err != nil {
		return err
	}

	err = git.Clone(git.CloneArgs{
		URL:  args.URL,
		Path: fmt.Sprintf("%s/%s", workspace.GetAbsolutePath(), args.Path),
	})
	if err != nil {
		multilog.Fatal("repositories.doctor", "failed to clone repository", map[string]interface{}{
			"repository": args,
		})
		return err
	}

	multilog.Info("repositories.add", "✅ added repository", map[string]interface{}{
		"repository": args,
	})

	return nil
}
