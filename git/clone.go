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

	if args.Auth == nil && util.GetProtocol(args.URL) == "ssh" {
		// Check if the default SSH key exists
		defaultSSHKey := util.ExpandPath("~/.ssh/id_rsa")
		if _, err := os.Stat(defaultSSHKey); err == nil {
			// Default SSH key exists, use it
			auth, err = ssh.NewPublicKeysFromFile("git", defaultSSHKey, "")
			if err != nil {
				multilog.Error("git.clone", "failed to create SSH auth with default key", map[string]interface{}{
					"error": err,
				})
				return err
			}
		} else {
			// Default SSH key doesn't exist, proceed without auth
			multilog.Warn("git.clone", "no auth provided and default SSH key not found", map[string]interface{}{
				"path": defaultSSHKey,
			})
		}
	} else if args.Auth != nil && args.Auth.Key != "" {
		auth, err = ssh.NewPublicKeysFromFile("git", args.Auth.Key, "")
		if err != nil {
			multilog.Fatal("git.clone", "failed to create SSH auth with provided key", map[string]interface{}{
				"error": err,
			})
			return err
		}
	} else if args.Auth != nil && args.Auth.Env.Username != "" && args.Auth.Env.Password != "" {
		auth = &http.BasicAuth{
			Username: os.Getenv(args.Auth.Env.Username),
			Password: os.Getenv(args.Auth.Env.Password),
		}
	}

	multilog.Info("git.clone", "cloning repository", map[string]interface{}{
		"url":  args.URL,
		"path": args.Path,
	})

	opts := &git.CloneOptions{
		URL:               args.URL,
		Progress:          &progress{},
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	if auth != nil {
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
