package generator

import (
	"os"

	"path/filepath"
)

func (g *Generator) copyAssets() error {
	if _, err := os.Stat(g.config.AssetsDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(g.config.AssetsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(g.config.AssetsDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(g.config.OutputDir, "assets", relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, data, info.Mode())
	})
}
