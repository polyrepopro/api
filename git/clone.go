package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/mateothegreat/go-multilog/multilog"
)

type CloneArgs struct {
	URL  string
	Path string
}

type progress struct{}

func (h *progress) Write(p []byte) (n int, err error) {
	multilog.Info("git.clone", "cloning progress", map[string]interface{}{
		"message": string(p),
	})
	return len(p), nil
}

func Clone(args CloneArgs) error {
	_, err := git.PlainClone(args.Path, true, &git.CloneOptions{
		URL:      args.URL,
		Progress: &progress{},
	})
	if err != nil {
		multilog.Fatal("git.clone", "failed to clone repository", map[string]interface{}{
			"error": err,
		})
	}

	return nil
}
