package git

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
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
		URL: args.URL,
		// Progress:          &progress{},
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	auth := GetAuth(args.URL, args.Auth)
	if auth != nil && auth.Name() != "" {
		opts.Auth = auth
	}

	if auth != nil {
		multilog.Debug("git.clone", "testing authentication", map[string]interface{}{
			"url": args.URL,
		})

		storer := memory.NewStorage()
		fs := memfs.New()

		_, err := git.Clone(storer, fs, &git.CloneOptions{
			URL:   args.URL,
			Auth:  auth,
			Depth: 1,
		})

		if err != nil && err != git.ErrRepositoryAlreadyExists {
			multilog.Error("git.clone", "authentication test failed", map[string]interface{}{
				"url":   args.URL,
				"error": err.Error(),
			})
			return err
		}

		multilog.Debug("git.clone", "authentication test successful", map[string]interface{}{
			"url": args.URL,
		})
	}

	_, err = git.PlainClone(args.Path, false, opts)
	if err != nil {
		multilog.Error("git.clone", "failed to clone repository", map[string]interface{}{
			"url":   args.URL,
			"path":  args.Path,
			"error": err.Error(),
		})
		return err
	}

	return nil
}
