package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"dugong/internal/config"
	"dugong/internal/generator"
	"dugong/internal/watcher"
)

var (
	configFile string
)

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "dugong",
	Short: "A static site generator",
	Long:  "Dugong is a static site generator that can be run as a build step or a service.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new dugong.toml configuration file",
	Long:  `Initialize a new Dugong project by creating a dugong.toml configuration file with interactive prompts.`,
	Run:   runInit,
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the site once",
	Long:  `Build the static site once from source files and exit.`,
	Run:   runBuild,
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Build the site and watch for changes",
	Long:  `Build the static site and watch for file changes, rebuilding automatically when changes are detected.`,
	Run:   runWatch,
}

func init() {
	initCmd.Flags().Bool("force", false, "Overwrite existing dugong.toml if present")
	initCmd.Flags().String("output", "dugong.toml", "Output path for config file")
	buildCmd.Flags().StringVar(&configFile, "config", "dugong.toml", "Path to config file")
	watchCmd.Flags().StringVar(&configFile, "config", "dugong.toml", "Path to config file")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(watchCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	force, _ := cmd.Flags().GetBool("force")
	outputPath, _ := cmd.Flags().GetString("output")

	if _, err := os.Stat(outputPath); err == nil && !force {
		fmt.Fprintf(os.Stderr, "Error: %s already exists. Use --force to overwrite.\n", outputPath)
		os.Exit(1)
	}

	cfg := config.Config{}
	fmt.Println("Dugong Static Site Generator - Configuration Setup")
	fmt.Println("Press Enter to accept the default value shown in [brackets]")
	reader := bufio.NewReader(os.Stdin)
	cfg.ContentDir = promptWithDefault(reader, "Content directory (for .md and .adoc files)", "./content")
	cfg.TemplateDir = promptWithDefault(reader, "Template directory (for HTML templates)", "./templates")
	cfg.AssetsDir = promptWithDefault(reader, "Assets directory (for CSS, JS, images)", "./assets")
	cfg.OutputDir = promptWithDefault(reader, "Output directory (for generated HTML)", "./output")
	fmt.Println()

	if err := cfg.WriteToFile(outputPath); err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}

	fmt.Printf("Created %s\n", outputPath)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Create your content directories:")
	fmt.Printf("     - %s (for .md and .adoc files)\n", cfg.ContentDir)
	fmt.Printf("     - %s (for HTML templates)\n", cfg.TemplateDir)
	fmt.Printf("     - %s (for CSS, JS, images)\n", cfg.AssetsDir)
	fmt.Println("  2. Run 'dugong' to generate your site")
}

func promptWithDefault(reader *bufio.Reader, prompt, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func runBuild(cmd *cobra.Command, args []string) {
	cfg := loadConfig()

	err := os.MkdirAll(cfg.OutputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	siteGenerator := generator.NewGenerator(cfg)

	log.Println("Building site...")
	err = siteGenerator.GenerateAll()
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}
	log.Println("Build complete")
}

func runWatch(cmd *cobra.Command, args []string) {
	cfg := loadConfig()

	err := os.MkdirAll(cfg.OutputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	siteGenerator := generator.NewGenerator(cfg)

	log.Println("Building site...")
	err = siteGenerator.GenerateAll()
	if err != nil {
		log.Fatalf("Initial build failed: %v", err)
	}
	log.Println("Initial build complete")

	fileWatcher, err := watcher.NewWatcher(cfg, siteGenerator)
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}

	err = fileWatcher.Watch()
	if err != nil {
		log.Fatalf("Watch failed: %v", err)
	}
}

func loadConfig() *config.Config {
	cfg, err := config.LoadFromFile(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	err = cfg.Validate()
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	return cfg
}
