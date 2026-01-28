package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	textTemplate "text/template"
)

type TemplateData struct {
	RenderedBody template.HTML
	Title        string
	RelPath      string
	Attributes   map[string]interface{}
	Pages        []string
	PagesByDir   map[string][]string
}

func TemplateFunctions() template.FuncMap {
	return template.FuncMap{}
}

func TextTemplateFunctions() textTemplate.FuncMap {
	return textTemplate.FuncMap{}
}

func isXMLTemplate(templatePath string) bool {
	return strings.HasSuffix(templatePath, ".xml")
}

// TODO: Rewrite this logic
func findTemplate(templateDir, contentDir, contentPath string) (string, error) {
	relPath, err := filepath.Rel(contentDir, contentPath)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(relPath)
	baseName := filepath.Base(contentDir)

	isIndex := strings.HasPrefix(baseName, "index.") || strings.HasPrefix(baseName, "_index.")

	if isIndex && dir != "." {
		indexTemplate := filepath.Join(templateDir, dir, "index_template.html")
		_, err := os.Stat(indexTemplate)
		if err == nil {
			return indexTemplate, nil
		}
	}

	if dir != "." {
		dirTemplate := filepath.Join(templateDir, dir, "page_template.html")
		_, err := os.Stat(dirTemplate)
		if err == nil {
			return dirTemplate, nil
		}
	}

	if isIndex {
		defaultIndexTemplate := filepath.Join(templateDir, "index_template.html")
		_, err := os.Stat(defaultIndexTemplate)
		if err == nil {
			return defaultIndexTemplate, nil
		}
	}

	defaultTemplate := filepath.Join(templateDir, "page_template.html")
	_, err = os.Stat(defaultTemplate)
	if err == nil {
		return defaultTemplate, nil
	}

	return "", fmt.Errorf("no template found")
}

func renderTemplate(templatePath string, data TemplateData, outputPath string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if isXMLTemplate(templatePath) {
		tmpl, err := textTemplate.New(filepath.Base(templatePath)).Funcs(TextTemplateFunctions()).ParseFiles(templatePath)
		if err != nil {
			return fmt.Errorf("failed to parse XML template: %w", err)
		}

		err = tmpl.Execute(outFile, data)
		if err != nil {
			return fmt.Errorf("failed to execute XML template: %w", err)
		}
	} else {
		tmpl, err := template.New(filepath.Base(templatePath)).Funcs(TemplateFunctions()).ParseFiles(templatePath)
		if err != nil {
			return fmt.Errorf("failed to parse HTML template: %w", err)
		}

		err = tmpl.Execute(outFile, data)
		if err != nil {
			return fmt.Errorf("failed to execute HTML template: %w", err)
		}
	}

	return nil
}
