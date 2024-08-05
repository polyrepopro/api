package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
)

type CloneArgs struct {
	URL  string
	Path string
	Auth *config.Auth
}

type progress struct{}

func (h *progress) Write(p []byte) (n int, err error) {
	multilog.Info("git.clone", "cloning progress", map[string]interface{}{
		"message": string(p),
	})
	return len(p), nil
}

func Clone(args CloneArgs) error {
	var err error

	multilog.Info("git.clone", "cloning repository", map[string]interface{}{
		"url":  args.URL,
		"path": args.Path,
	})

	opts := &git.CloneOptions{
		URL:               args.URL,
		Progress:          &progress{},
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	auth := GetAuth(args.URL, args.Auth)
	if auth.Name() != "" {
		opts.Auth = auth
	}

	_, err = git.PlainClone(args.Path, false, opts)
	if err != nil {
		multilog.Fatal("git.clone", "failed to clone repository", map[string]interface{}{
			"error": err,
		})
	}

	return nil
}
