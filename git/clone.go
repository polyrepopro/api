package git

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/config"
	"github.com/polyrepopro/api/util"
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
	var auth transport.AuthMethod
	var err error

	if args.Auth == nil {
		// Check if the default SSH key exists
		if _, err := os.Stat(util.ExpandPath("~/.ssh/id_rsa")); err == nil {
			// Default SSH key exists, use it
			auth, err = ssh.NewPublicKeysFromFile("git", "/Users/matthewdavis/.ssh/id_rsa", "")
			if err != nil {
				multilog.Error("git.clone", "failed to create SSH auth with default key", map[string]interface{}{
					"error": err,
				})
				return err
			}
		} else {
			// Default SSH key doesn't exist, proceed without auth
			multilog.Warn("git.clone", "no auth provided and default SSH key not found", map[string]interface{}{
				"path": "/Users/matthewdavis/.ssh/id_rsa",
			})
		}
	} else if args.Auth.Key != "" {
		auth, err = ssh.NewPublicKeysFromFile("git", args.Auth.Key, "")
		if err != nil {
			multilog.Fatal("git.clone", "failed to create SSH auth with provided key", map[string]interface{}{
				"error": err,
			})
			return err
		}
	} else if args.Auth.Env.Username != "" && args.Auth.Env.Password != "" {
		auth = &http.BasicAuth{
			Username: os.Getenv(args.Auth.Env.Username),
			Password: os.Getenv(args.Auth.Env.Password),
		}
	}

	multilog.Info("git.clone", "cloning repository", map[string]interface{}{
		"url":  args.URL,
		"path": args.Path,
	})

	_, err = git.PlainClone(args.Path, false, &git.CloneOptions{
		URL:               args.URL,
		Progress:          &progress{},
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              auth,
	})
	if err != nil {
		multilog.Fatal("git.clone", "failed to clone repository", map[string]interface{}{
			"error": err,
		})
	}

	return nil
}
