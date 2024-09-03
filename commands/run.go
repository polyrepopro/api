package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/mateothegreat/go-util/files"
	"github.com/polyrepopro/api/config"
)

func Run(ctx context.Context, label string, command config.Command, cwd string) error {
	if command.Cwd != "" {
		if err := os.Chdir(files.ExpandPath(command.Cwd)); err != nil {
			multilog.Error(label, "failed to change directory", map[string]interface{}{
				"command": command,
				"error":   err,
			})
			return err
		}
	} else if cwd != "" {
		if err := os.Chdir(files.ExpandPath(cwd)); err != nil {
			multilog.Error(label, "failed to change directory", map[string]interface{}{
				"command": command,
				"error":   err,
			})
			return err
		}
	}

	cmd := exec.CommandContext(ctx, command.Command[0], command.Command[1:]...)
	env := os.Environ()
	for k, v := range command.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = env

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		multilog.Error(label, "failed to get stdout pipe", map[string]interface{}{
			"command": command,
			"error":   err,
		})
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		multilog.Error(label, "failed to get stderr pipe", map[string]interface{}{
			"command": command,
			"error":   err,
		})
		return err
	}

	err = cmd.Start()
	if err != nil {
		multilog.Error(label, "failed to start command", map[string]interface{}{
			"name":  command.Name,
			"error": err,
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
				multilog.Info(label, "stdout", map[string]interface{}{
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
				return
			default:
				multilog.Error(label, "stderr", map[string]interface{}{
					"name":   command.Name,
					"output": scanner.Text(),
				})
			}
		}
	}()

	go func() {
		<-ctx.Done()
		syscall.Kill(cmd.Process.Pid, syscall.SIGKILL)
	}()

	return cmd.Wait()
}
