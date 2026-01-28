package generator

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"dugong/internal/config"
	"dugong/internal/parser"
	"path/filepath"
)

type Generator struct {
	config *config.Config
	graph  *DependencyGraph
}

func NewGenerator(config *config.Config) *Generator {
	return &Generator{
		config: config,
		graph:  NewDependencyGraph(),
	}
}

func (g *Generator) GenerateAll() error {
	var contentFiles []string

	absPath, err := filepath.Abs(g.config.OutputDir)
	if err != nil || absPath == "/" || strings.HasPrefix(absPath, "/etc") || strings.HasPrefix(absPath, "/usr") || strings.HasPrefix(absPath, "/bin") || strings.HasPrefix(absPath, "/home") && !strings.HasPrefix(absPath, os.Getenv("HOME")) {
		return fmt.Errorf("dangerous output directory will not remove")
	}

	err = os.RemoveAll(absPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(absPath, 0755)
	if err != nil {
		return err
	}

	err = filepath.Walk(g.config.ContentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		contentFiles = append(contentFiles, path)
		return nil
	})
	if err != nil {
		return err
	}

	for _, path := range contentFiles {
		err := g.renderContentFile(path)
		if err != nil {
			return err
		}
	}

	var templateFiles []string
	err = filepath.Walk(g.config.TemplateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		templateFiles = append(templateFiles, path)
		return nil
	})
	if err != nil {
		return err
	}

	usedTemplates := make(map[string]bool)
	for _, content := range g.graph.Content {
		usedTemplates[content.Template] = true
	}

	for _, path := range templateFiles {
		if !usedTemplates[path] {
			relPath, err := filepath.Rel(g.config.TemplateDir, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path: %w", err)
			}

			data := TemplateData{
				RelPath:    relPath,
				Pages:      []string{},
				PagesByDir: map[string][]string{},
			}

			outputPath := filepath.Join(g.config.OutputDir, relPath)
			err = renderTemplate(path, data, outputPath)
		}
	}

	err = g.copyAssets()
	if err != nil {
		log.Printf("Warning: failed to copy assets: %v", err)
	}

	return nil
}

func (g *Generator) renderContentFile(path string) error {
	contentParser, err := parser.NewParser(path)
	if err != nil {
		return err
	}

	doc, err := contentParser.Parse(path)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	relPath, err := filepath.Rel(g.config.ContentDir, path)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	templatePath, err := findTemplate(g.config.TemplateDir, g.config.ContentDir, path)
	if err != nil {
		return fmt.Errorf("failed to find template: %w", err)
	}

	data := TemplateData{
		RenderedBody: template.HTML(doc.Content),
		Title:        doc.Title,
		RelPath:      relPath,
		Attributes:   doc.Attributes,
		Pages:        []string{},
		PagesByDir:   map[string][]string{},
	}

	outputPath := filepath.Join(g.config.OutputDir, strings.TrimSuffix(relPath, filepath.Ext(relPath))+".html")

	err = renderTemplate(templatePath, data, outputPath)

	contentMeta := &Content{
		Path:      path,
		Template:  templatePath,
		IncludeIn: []string{},
	}
	g.graph.AddContent(contentMeta)
	return nil
}
