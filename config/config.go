package config

type Config struct {
	Working    string      `yaml:"working"`
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
