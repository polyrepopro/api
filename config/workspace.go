package config

import (
	"slices"
)

// Workspace allows you to group repositories together.
type Workspace struct {
	Name         string        `yaml:"name"`
	Path         string        `yaml:"path"`
	Repositories *[]Repository `yaml:"repositories" required:"false"`
	Auth         *Auth         `yaml:"auth,omitempty" required:"false"`
	Tags         []string      `yaml:"tags" required:"false"`
}

// GetRepositories returns the repositories for the workspace.
// If no tags are provided, all repositories are returned.
// If tags are provided, only repositories with the tags are returned.
//
// Arguments:
//   - tags: The tags to filter the repositories by.
//
// Returns:
//   - *[]Repository: The repositories.
func (w *Workspace) GetRepositories(tags []string) *[]Repository {
	var repositories []Repository
	for _, repository := range *w.Repositories {
		if len(tags) == 0 {
			repositories = append(repositories, repository)
			continue
		}
		for _, tag := range tags {
			if slices.Contains(repository.Tags, tag) {
				repositories = append(repositories, repository)
			}
		}
	}
	return &repositories
}
