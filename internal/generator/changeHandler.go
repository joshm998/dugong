package generator

import (
	"dugong/internal/config"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (g *Generator) HandleContentChange(change ContentChange) error {
	switch change.Op {
	case OpCreate, OpUpdate:
		return g.createOrUpdateContent(change.Path)
	case OpDelete:
		return g.deleteContent(change.Path)
	case OpRename:
		return g.renameContent(change.OldPath, change.Path)
	}
	return nil
}

func (g *Generator) createOrUpdateContent(path string) error {
	return g.renderContentFile(path)

	// Rebuild aggregates if needed
	// return nil
}

func (g *Generator) deleteContent(source string) error {
	return removeOldOutput(source, g.graph, g.config)
	// Rebuild aggregates
	//return nil
}

func (g *Generator) renameContent(oldPath, newPath string) error {
	err := removeOldOutput(oldPath, g.graph, g.config)
	if err != nil {
		return err
	}

	// Create new content (will rebuild and add to graph)
	return g.renderContentFile(newPath)

	//TODO: aggregates
}

func removeOldOutput(source string, graph *DependencyGraph, config *config.Config) error {
	contentMeta := graph.Content[source]
	if contentMeta == nil {
		return nil // already deleted
	}

	relPath, err := filepath.Rel(config.ContentDir, source)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	outputPath := filepath.Join(config.OutputDir, strings.TrimSuffix(relPath, filepath.Ext(relPath))+".html")

	if err := os.Remove(outputPath); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: failed to remove output file: %v", err)
	}

	graph.RemoveContent(source)
	return nil
}
