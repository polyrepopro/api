package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/util"
)

type Config struct {
	Workspaces []Workspace `yaml:"workspaces"`
}

type Workspace struct {
	Name         string       `yaml:"name"`
	Path         string       `yaml:"path"`
	Repositories []Repository `yaml:"repositories"`
}

type Repository struct {
	URL    string `yaml:"url"`
	Branch string `yaml:"branch"`
	Path   string `yaml:"path"`
}

func (c *Config) GetWorkspace(name string) (*Workspace, error) {
	for _, workspace := range c.Workspaces {
		if workspace.Name == name {
			return &workspace, nil
		}
	}
	return nil, fmt.Errorf("workspace %s not found", name)
}

func (c *Config) GetWorkspaceByWorkingDir() (*Workspace, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	for {
		for _, workspace := range c.Workspaces {
			expandedPath, err := util.ExpandPath(workspace.Path)
			if err != nil {
				return nil, fmt.Errorf("failed to expand path for workspace %s: %w", workspace.Name, err)
			}

			if util.IsSubPath(cwd, expandedPath) {
				return &workspace, nil
			}
		}

		// Move up one directory
		parent := filepath.Dir(cwd)
		if parent == cwd {
			// Reached the root directory
			break
		}
		cwd = parent
	}
	return nil, fmt.Errorf("no workspace found for current working directory or its parents")
}

func (c *Config) GetRepositoryByWorkingDir() (*Repository, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	workspace, err := c.GetWorkspaceByWorkingDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}

	for _, repository := range workspace.Repositories {
		if util.IsSubPath(cwd, repository.GetAbsolutePath()) {
			return &repository, nil
		}
	}
	return nil, fmt.Errorf("no repository found for current working directory")
}

func (w *Workspace) GetAbsolutePath() string {
	expandedPath, err := util.ExpandPath(w.Path)
	if err != nil {
		multilog.Fatal("config.GetAbsolutePath", "failed to expand path for workspace", map[string]interface{}{
			"error": err,
		})
	}
	return expandedPath
}

func (r *Repository) GetAbsolutePath() string {
	expandedPath, err := util.ExpandPath(r.Path)
	if err != nil {
		multilog.Fatal("config.GetAbsolutePath", "failed to expand path for repository", map[string]interface{}{
			"error": err,
		})
	}
	return expandedPath
}

// GetConfig returns a config hydrated by reading from .poly.yaml.
// It will walk up the directory tree to find the nearest .poly.yaml file.
//
// Returns:
//   - *Config: The hydratedconfig.
//   - error: An error if the config could not be found.
func GetConfig() (*Config, error) {
	var config *Config

	configPath := util.WalkFile(".poly.yaml", 10)
	if configPath != "" {
		config = &Config{}
		cleanenv.ReadConfig(configPath, &config)
	} else {
		cleanenv.ReadEnv(&config)
	}
	if config == nil {
		return nil, fmt.Errorf("base config not found in search paths")
	}

	emptyFields, err := util.ValidateStructFields(config, "")
	if err != nil {
		return nil, err
	}
	if len(emptyFields) > 0 {
		return nil, fmt.Errorf("empty fields: %v", emptyFields)
	}

	emptyFields, err = util.ValidateStructFields(config, "")
	if err != nil {
		return nil, err
	}
	if len(emptyFields) > 0 {
		return nil, fmt.Errorf("empty fields: %v", emptyFields)
	}

	return config, nil
}
