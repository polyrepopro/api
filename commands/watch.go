package commands

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/mateothegreat/go-util/files"
	"github.com/polyrepopro/api/config"
)

func Watch(watch config.Watch) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		multilog.Fatal("commands.watch", "failed to create watcher", map[string]interface{}{
			"error": err,
		})
	}
	defer watcher.Close()

	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	for _, path := range watch.Paths {
		if watch.Cwd != "" {
			path = filepath.Join(watch.Cwd, path)
		}

		matches, err := filepath.Glob(files.ExpandPath(path))
		if err != nil {
			multilog.Fatal("commands.watch", "failed to glob path", map[string]interface{}{
				"path":  path,
				"error": err,
			})
		}

		for _, match := range matches {
			err = watcher.Add(match)
			if err != nil {
				multilog.Fatal("commands.watch", "failed to add path to watcher", map[string]interface{}{
					"path":  match,
					"error": err,
				})
			}
		}

		multilog.Debug("commands.watch", "added path to watcher", map[string]interface{}{
			"path": path,
		})

		go func(path string) {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Op&fsnotify.Write == fsnotify.Write {
						multilog.Info("commands.watch", "file changed", map[string]interface{}{
							"path": event.Name,
						})
						cancel()
						ctx.Done()
						// Allow the OS time to catch up with the process being killed.
						time.Sleep(100 * time.Millisecond)
						ctx, cancel = RestartCommands(watch.Commands, files.ExpandPath(watch.Cwd))
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					multilog.Error("commands.watch", "watcher error", map[string]interface{}{
						"error": err,
					})
				}
			}
		}(path)
	}

	ctx, cancel = RestartCommands(watch.Commands, watch.Cwd)

	<-make(chan struct{})
}

func RestartCommands(commands []config.Command, cwd string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	for _, command := range commands {
		if command.Cwd != "" {
			command.Cwd = files.ExpandPath(filepath.Join(cwd, command.Cwd))
		}

		go func(command config.Command) {
			select {
			case <-ctx.Done():
				multilog.Info("commands.watch", "command execution cancelled via context", map[string]interface{}{
					"command": command,
				})
				return
			default:
				if command.Cwd != "" {
					if err := os.Chdir(files.ExpandPath(command.Cwd)); err != nil {
						multilog.Error("commands.watch", "failed to change directory", map[string]interface{}{
							"command": command,
							"error":   err,
						})
					}
				} else if cwd != "" {
					if err := os.Chdir(files.ExpandPath(cwd)); err != nil {
						multilog.Error("commands.watch", "failed to change directory", map[string]interface{}{
							"command": command,
							"error":   err,
						})
					}
				}

				cmd := exec.CommandContext(ctx, command.Command[0], command.Command[1:]...)
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					multilog.Error("commands.watch", "failed to get stdout pipe", map[string]interface{}{
						"command": command,
						"error":   err,
					})
					return
				}
				stderr, err := cmd.StderrPipe()
				if err != nil {
					multilog.Error("commands.watch", "failed to get stderr pipe", map[string]interface{}{
						"command": command,
						"error":   err,
					})
					return
				}

				go func() {
					scanner := bufio.NewScanner(stdout)
					for scanner.Scan() {
						select {
						case <-ctx.Done():
							return
						default:
							multilog.Info("commands.watch", "stdout", map[string]interface{}{
								"name":   command.Name,
								"output": scanner.Text(),
							})
						}
					}
					if err := scanner.Err(); err != nil {
						multilog.Error("commands.watch", "error reading stdout", map[string]interface{}{
							"command": command,
							"error":   err,
						})
					}
				}()

				go func() {
					scanner := bufio.NewScanner(stderr)
					for scanner.Scan() {
						select {
						case <-ctx.Done():
							return
						default:
							multilog.Info("commands.watch", "stderr", map[string]interface{}{
								"name":   command.Name,
								"output": scanner.Text(),
							})
						}
					}
					if err := scanner.Err(); err != nil {
						multilog.Error("commands.watch", "error reading stderr", map[string]interface{}{
							"command": command,
							"error":   err,
						})
					}
				}()

				err = cmd.Start()
				if err != nil {
					multilog.Error("commands.watch", "failed to start command", map[string]interface{}{
						"name":  command.Name,
						"error": err,
					})
					if command.ExitOnError {
						os.Exit(1)
					}
					return
				}

				multilog.Info("commands.watch", "started command", map[string]interface{}{
					"command": command,
					"cwd":     command.Cwd,
					"pid":     cmd.Process.Pid,
				})
			}
		}(command)
	}
	return ctx, cancel
}
