package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
)

type PushArgs struct {
	URL    string
	Remote string
	Path   string
	Auth   *config.Auth
}

type pushProgress struct{}

func (h *pushProgress) Write(p []byte) (n int, err error) {
	multilog.Debug("git.push", "pushing progress", map[string]interface{}{
		"message": string(p),
	})
	return len(p), nil
}

func Push(args PushArgs) error {
	repo, err := git.PlainOpen(args.Path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w for repo %q", err, args.Path)
	}

	opts := &git.PushOptions{
		RemoteName: args.Remote,
		Progress:   &pushProgress{},
	}

	auth := GetAuth(args.URL, args.Auth)
	if auth.Name() != "" {
		opts.Auth = auth
	}

	err = repo.Push(opts)
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to push changes: %w for repo %q", err, args.Path)
	}

	return nil

}
