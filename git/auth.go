package git

import (
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/mateothegreat/go-util/files"
	"github.com/mateothegreat/go-util/urls"
	"github.com/polyrepopro/api/config"
)

func GetAuth(url string, auth *config.Auth) transport.AuthMethod {
	if auth == nil && urls.GetProtocol(url) == "ssh" {
		// Check if the default SSH key exists
		defaultSSHKey := files.ExpandPath("~/.ssh/id_rsa")
		if _, err := os.Stat(defaultSSHKey); err == nil {
			// Default SSH key exists, use it
			sshAuth, err := ssh.NewPublicKeysFromFile("git", defaultSSHKey, "")
			if err != nil {
				multilog.Fatal("git.getauth", "failed to create SSH auth with default key", map[string]interface{}{
					"error": err.Error(),
				})
				return nil
			}

			multilog.Debug("auth", "using default SSH key", map[string]interface{}{
				"publicKey": defaultSSHKey,
				"url":       url,
			})

			return sshAuth
		} else {
			// Default SSH key doesn't exist, proceed without auth
			multilog.Fatal("git.getauth", "no auth provided and default SSH key not found", map[string]interface{}{
				"path": defaultSSHKey,
			})
		}
	} else if auth != nil && auth.Key != "" {
		// Use SSH key
		sshAuth, err := ssh.NewPublicKeysFromFile("git", auth.Key, "")
		if err != nil {
			multilog.Fatal("git.clone", "failed to create SSH auth with provided key", map[string]interface{}{
				"error": err.Error(),
			})
			return nil
		}

		multilog.Debug("auth", "using provided SSH key", map[string]interface{}{
			"publicKey": auth.Key,
			"url":       url,
		})

		return sshAuth
	} else if auth != nil && auth.Env.Username != "" && auth.Env.Password != "" {
		multilog.Debug("auth", "using HTTP auth", map[string]interface{}{
			"username": auth.Env.Username,
			"password": auth.Env.Password,
			"url":      url,
		})

		return &http.BasicAuth{
			Username: os.Getenv(auth.Env.Username),
			Password: os.Getenv(auth.Env.Password),
		}
	}
	return nil
}
