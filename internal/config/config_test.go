package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFromFile_NonExistent(t *testing.T) {
	_, err := LoadFromFile("nonexistent-config.toml")
	if err == nil {
		t.Fatalf("LoadFromFile() with non-existent file did not return err")
	}
}


func TestLoadFromFile_Valid(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "test.toml")
	configContent := `
	content_dir = "/my/content"
	output_dir = "/my/output"
	template_dir = "/my/templates"
	assets_dir = "/my/assets"
	`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("LoadFromFile() returned error: %v", err)
	}

	fixtures := []struct {
		name     string
		got      string
		expected string
	}{
		{"SourceDir", cfg.ContentDir, "/my/content"},
		{"OutputDir", cfg.OutputDir, "/my/output"},
		{"TemplateDir", cfg.TemplateDir, "/my/templates"},
		{"AssetsDir", cfg.AssetsDir, "/my/assets"},
	}

	for _, tt := range fixtures {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestLoadFromFile_EmptyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "empty.toml")
	if err := os.WriteFile(configPath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	_, err = LoadFromFile(configPath)
	if err == nil {
		t.Fatalf("LoadFromFile() with empty file did not return error")
	}
}

func TestValidate_MissingFields(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorContains string
	}{
		{
			name: "missing content_dir",
			config: &Config{
				ContentDir:  "",
				OutputDir:   "/my/output",
				TemplateDir: "/my/templates",
				AssetsDir:   "/my/assets",
			},
			expectError:   true,
			errorContains: "Content Directory",
		},
		{
			name: "missing output_dir",
			config: &Config{
				ContentDir:  "/my/content",
				OutputDir:   "",
				TemplateDir: "/my/templates",
				AssetsDir:   "/my/assets",
			},
			expectError:   true,
			errorContains: "Output Directory",
		},
		{
			name: "missing template_dir",
			config: &Config{
				ContentDir:  "/my/content",
				OutputDir:   "/my/output",
				TemplateDir: "",
				AssetsDir:   "/my/assets",
			},
			expectError:   true,
			errorContains: "Template Directory",
		},
		{
			name: "missing assets_dir",
			config: &Config{
				ContentDir:  "/my/content",
				OutputDir:   "/my/output",
				TemplateDir: "/my/templates",
				AssetsDir:   "",
			},
			expectError:   true,
			errorContains: "Assets Directory",
		},
		{
			name: "all fields missing",
			config: &Config{
				ContentDir:  "",
				OutputDir:   "",
				TemplateDir: "",
				AssetsDir:   "",
			},
			expectError:   true,
			errorContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError && err == nil {
				t.Errorf("Validate() expected error but got nil")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}

			if tt.expectError && err != nil && tt.errorContains != "" {
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Validate() error = %q, should contain %q", err.Error(), tt.errorContains)
				}
			}
		})
	}
}

func TestWriteToFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "test-write.toml")
	cfg := &Config{
		ContentDir:   "./content",
		OutputDir:   "./public",
		TemplateDir: "./templates",
		AssetsDir:   "./assets",
	}

	err = cfg.WriteToFile(configPath)
	if err != nil {
		t.Fatalf("WriteToFile() returned error: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("WriteToFile() did not create file")
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)

	expectedStrings := []string{
		"content_dir = \"./content\"",
		"output_dir = \"./public\"",
		"template_dir = \"./templates\"",
		"assets_dir = \"./assets\"",
		"# Dugong Configuration",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("WriteToFile() output missing expected string: %q", expected)
		}
	}
}
