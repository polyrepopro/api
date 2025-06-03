package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mateothegreat/go-util/files"
	"github.com/mateothegreat/go-util/validation"
	"gopkg.in/yaml.v3"
)

type DefaultArgs struct {
	Config string
}

type Config struct {
	Path       string       `yaml:"-"`
	URL        string       `yaml:"url" required:"false"`
	Default    string       `yaml:"default" required:"false"`
	Synced     time.Time    `yaml:"synced" required:"false"`
	Workspaces *[]Workspace `yaml:"workspaces" required:"false"`
}

type Defaults struct {
	Workspace string   `yaml:"workspace" required:"false"`
	Tags      []string `yaml:"tags" required:"false"`
}

type Repository struct {
	Name    string    `yaml:"name" required:"true"`
	URL     string    `yaml:"url" required:"true"`
	Origin  string    `yaml:"origin,omitempty" required:"false"`
	Branch  string    `yaml:"branch,omitempty" required:"false"`
	Path    string    `yaml:"path" required:"true"`
	Auth    *Auth     `yaml:"auth,omitempty" required:"false"`
	Hooks   *[]Hook   `yaml:"hooks,omitempty" required:"false"`
	Runners *[]Runner `yaml:"runners,omitempty" required:"false"`
	Tags    []string  `yaml:"tags,omitempty" required:"false"`
}

type HookType string

const (
	CloneHook   HookType = "clone"
	PullHook    HookType = "pull"
	PushHook    HookType = "push"
	PrePushHook HookType = "pre_push"
)

type Command struct {
	Name        string            `yaml:"name" required:"true"`
	Cwd         string            `yaml:"cwd" required:"false"`
	ExitOnError bool              `yaml:"exitOnError" required:"false"`
	Command     []string          `yaml:"command" required:"true"`
	Env         map[string]string `yaml:"env" required:"false"`
}

type Hook struct {
	Type     HookType  `yaml:"type" required:"true"`
	Commands []Command `yaml:"commands" required:"true"`
}

type Runner struct {
	Cwd      string    `yaml:"cwd" required:"false"`
	Watch    bool      `yaml:"watch" required:"false"`
	Matchers []Matcher `yaml:"matchers" required:"false"`
	Commands []Command `yaml:"commands" required:"true"`
}

type Matcher struct {
	Path    string `yaml:"path" required:"true"`
	Include string `yaml:"include" required:"false"`
	Ignore  string `yaml:"ignore" required:"false"`
}

type Auth struct {
	Key string  `yaml:"key,omitempty" required:"false"`
	Env AuthEnv `yaml:"env,omitempty" required:"false"`
}

type AuthEnv struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

func (c *Config) GetWorkspaces(names []string) (*[]Workspace, error) {
	if len(*c.Workspaces) == 0 {
		return nil, fmt.Errorf("no workspaces found in config")
	}
	if len(names) == 0 {
		return c.Workspaces, nil
	} else if len(names) == 1 && names[0] == "" {
		return c.Workspaces, nil
	}
	var workspaces []Workspace
	for _, name := range names {
		workspace, err := c.GetWorkspace(name)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, *workspace)
	}
	return &workspaces, nil
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
			if files.IsSubPath(cwd, files.ExpandPath(workspace.Path)) {
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
		if files.IsSubPath(cwd, repository.GetAbsolutePath()) {
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
	return files.ExpandPath(w.Path)
}

// GetAbsolutePath returns the absolute path of the repository.
//
// Returns:
//   - string: The absolute path of the repository.
func (r *Repository) GetAbsolutePath() string {
	return files.ExpandPath(r.Path)
}

// GetConfig returns a config hydrated by reading from a path.
//
// Returns:
//   - *Config: The hydrated config.
//   - error: An error if the config could not be found.
func GetConfig(path string) (*Config, error) {
	if path != "" {
		return GetAbsoluteConfig(path)
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
		Path: files.ExpandPath(path),
	}

	err := cleanenv.ReadConfig(files.ExpandPath(path), &config)
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
		configPath := files.WalkFile(".polyrepo.yaml", 10)
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
	emptyFields, err := validation.ValidateStructFields(config, "")
	if err != nil {
		return nil, err
	}
	if len(emptyFields) > 0 {
		return nil, fmt.Errorf("empty fields: %v", emptyFields)
	}

	return config, nil
}
