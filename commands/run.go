package commands

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/mateothegreat/go-util/files"
	"github.com/polyrepopro/api/config"
)

func Run(ctx context.Context, label string, command config.Command, cwd string) *os.Process {
	if command.Cwd != "" {
		if err := os.Chdir(files.ExpandPath(command.Cwd)); err != nil {
			multilog.Error(label, "failed to change directory", map[string]interface{}{
				"command": command,
				"error":   err,
			})
			return nil
		}
	} else if cwd != "" {
		if err := os.Chdir(files.ExpandPath(cwd)); err != nil {
			multilog.Error(label, "failed to change directory", map[string]interface{}{
				"command": command,
				"error":   err,
			})
			return nil
		}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	cmd := exec.CommandContext(ctx, command.Command[0], command.Command[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		multilog.Error(label, "failed to get stdout pipe", map[string]interface{}{
			"command": command,
			"error":   err,
		})
		return nil
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		multilog.Error(label, "failed to get stderr pipe", map[string]interface{}{
			"command": command,
			"error":   err,
		})
		return nil
	}

	err = cmd.Start()
	if err != nil {
		multilog.Error(label, "failed to start command", map[string]interface{}{
			"name":  command.Name,
			"error": err,
		})
		return nil
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case <-sig:
				multilog.Info(label, "stdout", map[string]interface{}{
					"name":   command.Name,
					"output": scanner.Text(),
				})
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
			case <-sig:
				multilog.Info(label, "stderr", map[string]interface{}{
					"name":   command.Name,
					"output": scanner.Text(),
				})
				return
			case <-ctx.Done():
				return
			default:
				multilog.Info(label, "stderr", map[string]interface{}{
					"name":   command.Name,
					"output": scanner.Text(),
				})
			}
		}
	}()

	go func() {
		<-ctx.Done()
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}()

	return cmd.Process
}
