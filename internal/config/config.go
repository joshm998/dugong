package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ContentDir  string `toml:"content_dir"`
	OutputDir   string `toml:"output_dir"`
	TemplateDir string `toml:"template_dir"`
	AssetsDir   string `toml:"assets_dir"`
}

func LoadFromFile(path string) (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file does not exist: %w", err)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	err = toml.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	err = cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("validation Failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.ContentDir == "" {
		return fmt.Errorf("Content Directory is Empty")
	}
	if c.OutputDir == "" {
		return fmt.Errorf("Output Directory is Empty")
	}
	if c.TemplateDir == "" {
		return fmt.Errorf("Template Directory is Empty")
	}
	if c.AssetsDir == "" {
		return fmt.Errorf("Assets Directory is Empty")
	}
	return nil
}

func (c *Config) WriteToFile(path string) error {
	var buf bytes.Buffer

	buf.WriteString("# Dugong Configuration\n\n")

	writeFieldWithComment(&buf, "Content directory containing your content files (.md, .adoc)", "content_dir", c.ContentDir)
	writeFieldWithComment(&buf, "Output directory for generated HTML files", "output_dir", c.OutputDir)
	writeFieldWithComment(&buf, "Directory containing your HTML templates", "template_dir", c.TemplateDir)
	writeFieldWithComment(&buf, "Directory containing static assets (images, CSS, JS)", "assets_dir", c.AssetsDir)

	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func writeFieldWithComment(buf *bytes.Buffer, comment, key, value string) {
	buf.WriteString("# ")
	buf.WriteString(comment)
	buf.WriteString("\n")

	data := map[string]string{key: value}
	encoder := toml.NewEncoder(buf)
	if err := encoder.Encode(data); err != nil {
		fmt.Printf("Unable to write field: %s", err)
	}
	buf.WriteString("\n")
}
