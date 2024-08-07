package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/polyrepopro/api/util"
	"gopkg.in/yaml.v3"
)

type DefaultArgs struct {
	Config string
}

type Config struct {
	Path       string       `yaml:"-"`
	Synced     time.Time    `yaml:"synced" required:"false"`
	Workspaces *[]Workspace `yaml:"workspaces" required:"false"`
}

type Workspace struct {
	Name         string        `yaml:"name"`
	Path         string        `yaml:"path"`
	Repositories *[]Repository `yaml:"repositories" required:"false"`
}

type Repository struct {
	URL    string `yaml:"url"`
	Origin string `yaml:"origin,omitempty"`
	Branch string `yaml:"branch,omitempty"`
	Path   string `yaml:"path"`
	Auth   *Auth  `yaml:"auth,omitempty"`
}

type Auth struct {
	Key string  `yaml:"key,omitempty" required:"false"`
	Env AuthEnv `yaml:"env,omitempty" required:"false"`
}

type AuthEnv struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

// SaveConfig saves the config to the path specified by the config.
//
// Returns:
//   - error: An error if the config could not be saved.
func (c *Config) SaveConfig() error {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	err := encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(c.Path, buf.Bytes(), 0744)
}

// GetWorkspace returns a workspace by name.
//
// Returns:
//   - *Workspace: The workspace.
//   - error: An error if the workspace could not be found.
func (c *Config) GetWorkspace(name string) (*Workspace, error) {
	for _, workspace := range *c.Workspaces {
		if workspace.Name == name {
			return &workspace, nil
		}
	}
	return nil, fmt.Errorf("workspace %s not found", name)
}

// GetWorkspaceByWorkingDir returns a workspace by searching for the current working directory in the workspace.
//
// Returns:
//   - *Workspace: The workspace.
//   - error: An error if the workspace could not be found.
func (c *Config) GetWorkspaceByWorkingDir() (*Workspace, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	for {
		for _, workspace := range *c.Workspaces {
			if util.IsSubPath(cwd, util.ExpandPath(workspace.Path)) {
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

// GetRepositoryByWorkingDir returns a repository by searching for the current working directory in the workspace.
//
// Returns:
//   - *Repository: The repository.
//   - error: An error if the repository could not be found.
func (c *Config) GetRepositoryByWorkingDir() (*Repository, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	workspace, err := c.GetWorkspaceByWorkingDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}

	for _, repository := range *workspace.Repositories {
		if util.IsSubPath(cwd, repository.GetAbsolutePath()) {
			return &repository, nil
		}
	}
	return nil, fmt.Errorf("no repository found for current working directory")
}

// GetAbsolutePath returns the absolute path of the workspace.
//
// Returns:
//   - string: The absolute path of the workspace.
func (w *Workspace) GetAbsolutePath() string {
	return util.ExpandPath(w.Path)
}

// GetAbsolutePath returns the absolute path of the repository.
//
// Returns:
//   - string: The absolute path of the repository.
func (r *Repository) GetAbsolutePath() string {
	return util.ExpandPath(r.Path)
}

// GetConfig returns a config hydrated by reading from a path.
//
// Returns:
//   - *Config: The hydrated config.
//   - error: An error if the config could not be found.
func GetConfig(path *string) (*Config, error) {
	if path != nil {
		return GetAbsoluteConfig(*path)
	}
	return GetRelativeConfig()
}

// GetAbsoluteConfig returns a config hydrated by reading from a path.
//
// Arguments:
//   - path: The path to the config file.
//
// Returns:
//   - *Config: The hydrated config.
//   - error: An error if the config could not be found.
func GetAbsoluteConfig(path string) (*Config, error) {
	config := Config{
		Path: util.ExpandPath(path),
	}

	err := cleanenv.ReadConfig(util.ExpandPath(path), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &config, nil
}

// GetRelativeConfig returns a config hydrated by reading from .polyrepo.yaml.
// It will walk up the directory tree to find the nearest .polyrepo.yaml file.
//
// Returns:
//   - *Config: The hydrated config.
//   - error: An error if the config could not be found.
func GetRelativeConfig() (*Config, error) {
	var config *Config

	if os.Getenv("POLYREPO_CONFIG") != "" {
		// If the POLYREPO_CONFIG environment variable is set, use it.
		configPath := os.Getenv("POLYREPO_CONFIG")
		config = &Config{}
		cleanenv.ReadConfig(configPath, &config)
		log.Printf("Using config from POLYREPO_CONFIG: %s", configPath)
	} else {
		configPath := util.WalkFile(".polyrepo.yaml", 10)
		if configPath != "" {
			config = &Config{}
			cleanenv.ReadConfig(configPath, &config)
		} else {
			cleanenv.ReadEnv(&config)
		}
		if config == nil {
			return nil, fmt.Errorf("base config not found in search paths")
		}

		// Set the path to the config path found for things like saving later.
		config.Path = configPath
	}

	// Validate the config against empty fields.
	emptyFields, err := util.ValidateStructFields(config, "")
	if err != nil {
		return nil, err
	}
	if len(emptyFields) > 0 {
		return nil, fmt.Errorf("empty fields: %v", emptyFields)
	}

	return config, nil
}
