package config

import (
	"slices"
)

type Workspace struct {
	Name         string        `yaml:"name"`
	Path         string        `yaml:"path"`
	Repositories *[]Repository `yaml:"repositories" required:"false"`
	Tags         []string      `yaml:"tags" required:"false"`
}

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
