package commands

import (
	"context"
	"io/fs"
	"path/filepath"
	"regexp"

	"github.com/fsnotify/fsnotify"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/mateothegreat/go-util/files"
	"github.com/polyrepopro/api/config"
)

func Watch(ctx context.Context, label string, workspacePath string, runner config.Runner) {
	for {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			multilog.Fatal(label, "failed to create watcher", map[string]interface{}{
				"error": err,
			})
			return
		}

		for _, command := range runner.Commands {
			var c context.Context
			var cancel context.CancelFunc

			var base string
			if runner.Cwd != "" {
				base = files.ExpandPath(filepath.Join(workspacePath, runner.Cwd))
			} else {
				base = files.ExpandPath(filepath.Join(workspacePath, command.Cwd))
			}

			var matches []string
			for _, matcher := range runner.Matchers {
				path := matcher.Path
				if path == "" {
					path = command.Cwd
				}

				err := filepath.WalkDir(filepath.Join(base, path), func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if d.IsDir() {
						return nil
					}

					ignore := false
					if matcher.Ignore != "" {
						matched, err := regexp.MatchString(matcher.Ignore, path)
						if err != nil {
							multilog.Fatal(label, "failed to match path", map[string]interface{}{
								"path":    path,
								"error":   err,
								"pattern": matcher.Ignore,
							})
							return err
						}
						if matched {
							ignore = true
						}
					}

					matched, err := regexp.MatchString(matcher.Include, path)
					if err != nil {
						multilog.Fatal(label, "failed to match path", map[string]interface{}{
							"path":    path,
							"error":   err,
							"pattern": matcher.Include,
						})
						return err
					}
					if matched && !ignore {
						matches = append(matches, path)
					}

					return nil
				})
				if err != nil {
					multilog.Fatal(label, "failed to walk path", map[string]interface{}{
						"path":  path,
						"error": err,
					})
				}
			}

			for _, match := range matches {

				err = watcher.Add(match)
				if err != nil {
					multilog.Fatal(label, "failed to add path to watcher", map[string]interface{}{
						"path":  match,
						"error": err,
					})
				}
			}

			c, cancel = context.WithCancel(context.Background())
			go func() {
				for {
					select {
					case <-ctx.Done():
						println("received sigterm, canceling context")
						cancel()
						return
					case <-c.Done():
						return
					case event, ok := <-watcher.Events:
						if !ok {
							multilog.Fatal(label, "watcher events channel closed", nil)
						}
						if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove {
							multilog.Info(label, "change detected", map[string]interface{}{
								"path": event.Name,
								"op":   event.Op,
							})
							cancel()
						}
					case err, ok := <-watcher.Errors:
						if !ok {
							multilog.Fatal(label, "watcher errors channel closed", nil)
						}
						multilog.Fatal(label, "watcher error", map[string]interface{}{
							"error": err,
						})
						cancel()
					}
				}
			}()
			Run(c, label, command, base)
			<-c.Done()
			watcher.Close()
		}
	}
}
