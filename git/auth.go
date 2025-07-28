package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/mateothegreat/go-util/files"
	"github.com/mateothegreat/go-util/urls"
	"github.com/polyrepopro/api/config"
	gossh "golang.org/x/crypto/ssh"
)

var defaultKeys = []string{
	"~/.ssh/id_rsa",
	"~/.ssh/id_ed25519",
}

// getCredentialsFromHelper calls git credential fill to get credentials from configured helpers
func getCredentialsFromHelper(url string) (username, password string, err error) {
	// Create a context with timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, "git", "credential", "fill")
	
	// Prepare input for git credential fill
	input := fmt.Sprintf("url=%s\n\n", url)
	cmd.Stdin = strings.NewReader(input)
	
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get credentials from git credential helper: %w", err)
	}
	
	// Parse the output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "username=") {
			username = strings.TrimPrefix(line, "username=")
		} else if strings.HasPrefix(line, "password=") {
			password = strings.TrimPrefix(line, "password=")
		}
	}
	
	if username == "" || password == "" {
		return "", "", fmt.Errorf("git credential helper did not provide username/password")
	}
	
	return username, password, nil
}

func GetAuth(url string, auth *config.Auth) transport.AuthMethod {
	if auth == nil {
		protocol := urls.GetProtocol(url)
		multilog.Debug("GetAuth", "protocol detection", map[string]interface{}{
			"url":      url,
			"protocol": protocol,
		})
		
		if protocol == "ssh" {
			// Try SSH keys directly
			for _, key := range defaultKeys {
				defaultSSHKey := files.ExpandPath(key)
				if _, err := os.Stat(defaultSSHKey); err == nil {
					// Default SSH key exists, use it
					sshAuth, err := ssh.NewPublicKeysFromFile("git", defaultSSHKey, "")
					if err != nil {
						multilog.Debug("git.getauth", "failed to create SSH auth with default key", map[string]interface{}{
							"error": err.Error(),
							"key":   defaultSSHKey,
						})
						continue // Try next key instead of returning nil
					}

					// Set up host key callback to accept any host key (less secure but more compatible)
					sshAuth.HostKeyCallback = gossh.InsecureIgnoreHostKey()

					multilog.Debug("GetAuth", "using SSH key file", map[string]interface{}{
						"publicKey": defaultSSHKey,
						"url":       url,
					})

					return sshAuth
				}
			}
		} else if protocol == "https" || protocol == "http" {
			// Try Git credential helpers for HTTPS URLs
			username, password, err := getCredentialsFromHelper(url)
			if err == nil {
				multilog.Debug("GetAuth", "using credential helper", map[string]interface{}{
					"url":      url,
					"username": username,
				})
				return &http.BasicAuth{
					Username: username,
					Password: password,
				}
			}
			
			multilog.Debug("GetAuth", "credential helper failed", map[string]interface{}{
				"url":   url,
				"error": err.Error(),
			})
		}

		// Try SSH agent as fallback
		sshAuth, err := ssh.NewSSHAgentAuth("git")
		if err == nil {
			multilog.Debug("GetAuth", "using SSH agent as fallback", map[string]interface{}{
				"url": url,
			})
			return sshAuth
		}

		// No valid SSH keys found
		multilog.Debug("git.getauth", "no auth provided and no valid SSH keys or agent found", map[string]interface{}{
			"keys":        defaultKeys,
			"agent_error": err.Error(),
		})
		return nil
	} else if auth != nil && auth.Key != "" {
		// Use SSH key
		sshAuth, err := ssh.NewPublicKeysFromFile("git", files.ExpandPath(auth.Key), "")
		if err != nil {
			multilog.Fatal("git.clone", "failed to create SSH auth with provided key", map[string]interface{}{
				"error": err.Error(),
			})
			return nil
		}

		// Set up host key callback to accept any host key (less secure but more compatible)
		sshAuth.HostKeyCallback = gossh.InsecureIgnoreHostKey()

		multilog.Debug("GetAuth", "using provided SSH key", map[string]interface{}{
			"publicKey": auth.Key,
			"url":       url,
		})

		return sshAuth
	} else if auth != nil && auth.Env.Username != "" && auth.Env.Password != "" {
		multilog.Debug("GetAuth", "using HTTP auth", map[string]interface{}{
			"username": auth.Env.Username,
			"password": auth.Env.Password,
			"url":      url,
		})

		return &http.BasicAuth{
			Username: os.Getenv(auth.Env.Username),
			Password: os.Getenv(auth.Env.Password),
		}
	}

	multilog.Debug("git.getauth", "no auth could be found", map[string]interface{}{
		"url": url,
	})
	return nil
}
