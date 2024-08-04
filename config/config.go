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

type Config struct {
	Path       string      `yaml:"-"`
	SyncedAt   time.Time   `yaml:"synced"`
	Workspaces []Workspace `yaml:"workspaces"`
}

type Workspace struct {
	Name         string        `yaml:"name"`
	Path         string        `yaml:"path"`
	Repositories *[]Repository `yaml:"repositories"`
}

type Repository struct {
	URL    string `yaml:"url"`
	Branch string `yaml:"branch,omitempty"`
	Path   string `yaml:"path"`
	Auth   Auth   `yaml:"auth,omitempty"`
}

type Auth struct {
	Key string  `yaml:"key,omitempty"`
	Env AuthEnv `yaml:"env,omitempty"`
}

type AuthEnv struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

func (c *Config) SaveConfig() error {
	configPath := util.WalkFile(c.Path, 10)
	if configPath == "" {
		return fmt.Errorf("config not found in search paths")
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	err := encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, buf.Bytes(), 0644)
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

func (w *Workspace) GetAbsolutePath() string {
	return util.ExpandPath(w.Path)
}

func (r *Repository) GetAbsolutePath() string {
	return util.ExpandPath(r.Path)
}

func GetAbsoluteConfig(path string) (*Config, error) {
	config := Config{}

	err := cleanenv.ReadConfig(util.ExpandPath(path), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &config, nil
}

// GetRelativeConfig returns a config hydrated by reading from .poly.yaml.
// It will walk up the directory tree to find the nearest .poly.yaml file.
//
// Returns:
//   - *Config: The hydratedconfig.
//   - error: An error if the config could not be found.
func GetRelativeConfig() (*Config, error) {
	var config *Config

	if os.Getenv("POLYREPO_CONFIG") != "" {
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
