package workspaces

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/polyrepopro/api/util"
)

type InitArgs struct {
	Path string
	URL  string
}

func Init(args InitArgs) error {
	path := util.ExpandPath(args.Path)

	// Check if the base directory exists, if not create it
	baseDir := filepath.Dir(path)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		multilog.Error("workspaces.init", "failed to create base directory", map[string]interface{}{
			"error": err,
			"path":  baseDir,
		})
		return err
	}

	if args.URL != "" {
		// Download the URL file and save to disk
		resp, err := http.Get(args.URL)
		if err != nil {
			multilog.Fatal("workspaces.init", "failed to download config file", map[string]interface{}{
				"error": err,
				"url":   args.URL,
			})
		}
		defer resp.Body.Close()

		// Create the file
		out, err := os.Create(path)
		if err != nil {
			multilog.Fatal("workspaces.init", "failed to create file", map[string]interface{}{
				"error": err,
				"path":  path,
			})
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			multilog.Fatal("workspaces.init", "failed to write file", map[string]interface{}{
				"error": err,
				"path":  path,
			})
		}
	}

	return nil
}
