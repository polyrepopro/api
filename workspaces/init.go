package workspaces

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mateothegreat/go-util/files"
	"github.com/polyrepopro/api/config"
)

type InitArgs struct {
	Path string
	URL  string
}

func Init(args InitArgs) (*config.Config, error) {
	var cfg *config.Config

	path := files.ExpandPath(args.Path)

	// Check if the base directory exists, if not create it
	baseDir := filepath.Dir(path)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	// If a URL is provided, download the file and save it to disk.
	if args.URL != "" {
		// Download the URL file and save to disk
		resp, err := http.Get(args.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to download config file: %w", err)
		}
		defer resp.Body.Close()

		// Create the file
		out, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to write file: %w", err)
		}

		cfg, err = config.GetConfig(&path)
		if err != nil {
			return nil, fmt.Errorf("failed to get config: %w", err)
		}
	} else {
		// Create a new config file with defaults.
		cfg := config.Config{
			Path: path,
		}

		// Save the config file to disk.
		cfg.SaveConfig()
	}

	return cfg, nil
}
