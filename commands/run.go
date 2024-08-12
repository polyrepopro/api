package commands

import (
	"bufio"
	"context"
	"os"
	"os/exec"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/mateothegreat/go-util/files"
	"github.com/polyrepopro/api/config"
)

func Run(ctx context.Context, command config.Command, cwd string) error {
	if command.Cwd != "" {
		if err := os.Chdir(files.ExpandPath(command.Cwd)); err != nil {
			multilog.Error("commands.run", "failed to change directory", map[string]interface{}{
				"command": command,
				"error":   err,
			})
			return err
		}
	} else if cwd != "" {
		if err := os.Chdir(files.ExpandPath(cwd)); err != nil {
			multilog.Error("commands.run", "failed to change directory", map[string]interface{}{
				"command": command,
				"error":   err,
			})
			return err
		}
	}

	cmd := exec.CommandContext(ctx, command.Command[0], command.Command[1:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		multilog.Error("commands.run", "failed to get stdout pipe", map[string]interface{}{
			"command": command,
			"error":   err,
		})
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		multilog.Error("commands.run", "failed to get stderr pipe", map[string]interface{}{
			"command": command,
			"error":   err,
		})
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				multilog.Info("commands.run", "stdout", map[string]interface{}{
					"name":   command.Name,
					"output": scanner.Text(),
				})
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				multilog.Info("commands.run", "context done", map[string]interface{}{
					"name": command.Name,
				})
				return
			default:
				multilog.Info("commands.run", "stderr", map[string]interface{}{
					"name":   command.Name,
					"output": scanner.Text(),
				})
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		multilog.Error("commands.run", "failed to start command", map[string]interface{}{
			"name":  command.Name,
			"error": err,
		})
		return err
	}

	multilog.Debug("commands.run", "started command", map[string]interface{}{
		"command": command,
		"cwd":     command.Cwd,
		"pid":     cmd.Process.Pid,
	})

	err = cmd.Wait()
	if err != nil {
		multilog.Error("commands.run", "command execution failed", map[string]interface{}{
			"name":  command.Name,
			"error": err,
		})
	}

	return nil
}
