package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	config := &CloneConfig{
		Name:        "",
		Destination: "",
		Repo:        "github.com/openshift-online",
	}

	flag.StringVar(&config.Name, "name", "", "Name of the new project")
	flag.StringVar(&config.Destination, "destination", "", "Destination directory")
	flag.StringVar(&config.Repo, "repo", "github.com/openshift-online", "Git repository")
	flag.Parse()

	if config.Name == "" || config.Destination == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s --name <project-name> --destination <path>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s --name user-service --destination ~/projects/user-service\n", os.Args[0])
		os.Exit(1)
	}

	if err := cloneProject(config); err != nil {
		fmt.Fprintf(os.Stderr, "Clone failed: %v\n", err)
		os.Exit(1)
	}
}

type CloneConfig struct {
	Name        string
	Repo        string
	Destination string
}

type FileCategory int

const (
	ModuleFile FileCategory = iota
	GoSourceFile
	OpenAPIFile
	InfrastructureFile
	PathStructureFile
	DocumentationFile
	ConfigurationFile
	SkipFile
)

func cloneProject(config *CloneConfig) error {
	fmt.Printf("🚀 Cloning TRex to %s...\n", config.Name)

	// Get current directory (source) - need to go up two levels from scripts/clone/
	sourceDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Go up to the TRex root directory
	sourceDir = filepath.Join(sourceDir, "..", "..")
	sourceDir, err = filepath.Abs(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to resolve source directory: %v", err)
	}

	// Create destination directory
	if err := os.MkdirAll(config.Destination, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Copy files with replacements
	if err := copyWithReplacements(sourceDir, config.Destination, config); err != nil {
		return fmt.Errorf("failed to copy files: %v", err)
	}

	fmt.Printf("✅ Clone completed successfully!\n")
	fmt.Printf("📁 Location: %s\n", config.Destination)
	fmt.Printf("📋 Next steps (copy and paste these commands):\n\n")
	fmt.Printf("cd %s &&\n", config.Destination)
	fmt.Printf("go mod tidy &&\n")
	fmt.Printf("make db/setup &&\n")
	fmt.Printf("make binary &&\n")
	fmt.Printf("go run scripts/generate/main.go --kind YourEntity &&\n")
	fmt.Printf("make generate &&\n")
	fmt.Printf("make test && make test-integration\n")

	return nil
}

func copyWithReplacements(srcDir, dstDir string, config *CloneConfig) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip scripts directory to avoid recursive copying
		if strings.Contains(path, "scripts/") {
			return nil
		}

		// Skip git and other build artifacts
		if shouldSkipPath(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Create relative path
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFileWithReplacements(path, dstPath, config)
	})
}

func copyFileWithReplacements(srcPath, dstPath string, config *CloneConfig) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	content, err := io.ReadAll(srcFile)
	if err != nil {
		return err
	}

	// Apply replacements based on file type

	fmt.Printf("Processing file %s to %s:", srcPath, dstPath)
	processedContent := processFileContent(string(content), srcPath, config)

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	written, err := dstFile.WriteString(processedContent)
	if err != nil {
		return err
	}

	fmt.Printf("   wrote: %b\n", written)
	fmt.Printf("")

	return err
}

func shouldSkipPath(path string) bool {
	skipPatterns := []string{
		".git",
		"vendor/",
		"build/",
		"*.log",
		".trex.md",
		"clones/",
		"demos/",
	}

	for _, pattern := range skipPatterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

func processFileContent(content, path string, config *CloneConfig) string {
	category := categorizeFile(path)

	switch category {
	case ModuleFile:
		return processModuleFile(content, config)
	case GoSourceFile:
		return processGoSourceFile(content, config)
	case OpenAPIFile:
		return processOpenAPIFile(content, config)
	case InfrastructureFile:
		return processInfrastructureFile(content, config)
	case DocumentationFile:
		return processDocumentationFile(content, path, config)
	default:
		return content
	}
}

// File categorization and processing functions
func categorizeFile(path string) FileCategory {
	filename := filepath.Base(path)

	if filename == "go.mod" {
		return ModuleFile
	}

	if strings.HasSuffix(path, ".go") {
		return GoSourceFile
	}

	if strings.Contains(path, "openapi") && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
		return OpenAPIFile
	}

	if filename == "Makefile" || filename == "Dockerfile" || strings.Contains(path, "templates/") || strings.Contains(path, ".tekton/") {
		return InfrastructureFile
	}

	if strings.HasSuffix(path, ".md") {
		return DocumentationFile
	}

	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		return ConfigurationFile
	}

	return SkipFile
}

func processModuleFile(content string, config *CloneConfig) string {
	// Replace module declaration
	content = strings.ReplaceAll(content, "module github.com/openshift-online/rh-trex", fmt.Sprintf("module %s/%s", config.Repo, config.Name))
	return content
}

func processGoSourceFile(content string, config *CloneConfig) string {
	// Replace import paths
	content = strings.ReplaceAll(content, "github.com/openshift-online/rh-trex/", fmt.Sprintf("%s/%s/", config.Repo, config.Name))

	// Replace package trex declaration
	if strings.Contains(content, "package trex") {
		binaryName := toBinaryName(config.Name)
		content = strings.ReplaceAll(content, "package trex", fmt.Sprintf("package %s", binaryName))
	}

	return content
}

func processOpenAPIFile(content string, config *CloneConfig) string {
	// Replace API paths
	content = strings.ReplaceAll(content, "/api/rh-trex/", fmt.Sprintf("/api/%s/", config.Name))

	// Replace operation IDs
	camelName := toCamelCase(config.Name)
	content = strings.ReplaceAll(content, "ApiRhTrexV1", fmt.Sprintf("Api%sV1", camelName))

	// Replace service descriptions
	content = strings.ReplaceAll(content, "rh-trex Service API", fmt.Sprintf("%s Service API", config.Name))

	return content
}

func processInfrastructureFile(content string, config *CloneConfig) string {
	binaryName := toBinaryName(config.Name)
	sqlName := toSqlSafeName(config.Name)

	// Database and service replacements
	content = strings.ReplaceAll(content, "psql-rhtrex", fmt.Sprintf("psql-%s", sqlName))
	content = strings.ReplaceAll(content, "rhtrex", sqlName)

	// Binary and service name replacements
	content = strings.ReplaceAll(content, "/usr/local/bin/trex", fmt.Sprintf("/usr/local/bin/%s", binaryName))
	content = strings.ReplaceAll(content, "trex-service", fmt.Sprintf("%s-service", binaryName))
	content = strings.ReplaceAll(content, "trex-api", fmt.Sprintf("%s-api", binaryName))
	content = strings.ReplaceAll(content, "trex-metrics", fmt.Sprintf("%s-metrics", binaryName))

	// Template and label replacements
	replacements := map[string]string{
		"name: trex":     fmt.Sprintf("name: %s", binaryName),
		"app: trex":      fmt.Sprintf("app: %s", binaryName),
		"template: trex": fmt.Sprintf("template: %s", binaryName),
	}

	for old, new := range replacements {
		content = strings.ReplaceAll(content, old, new)
	}

	// API path replacements
	content = strings.ReplaceAll(content, "/api/rh-trex", fmt.Sprintf("/api/%s", config.Name))
	content = strings.ReplaceAll(content, "cmd/trex/", fmt.Sprintf("cmd/%s/", binaryName))

	return content
}

func processDocumentationFile(content, path string, config *CloneConfig) string {
	// Replace project references
	content = strings.ReplaceAll(content, "rh-trex Service", fmt.Sprintf("%s Service", config.Name))
	content = strings.ReplaceAll(content, "rh-trex template", fmt.Sprintf("%s template", config.Name))
	content = strings.ReplaceAll(content, "TRex service", fmt.Sprintf("%s service", config.Name))

	// Add clone section to CLAUDE.md
	if strings.HasSuffix(path, "CLAUDE.md") {
		content = addTRexCloneSection(content, config.Name)
	}

	return content
}

func addTRexCloneSection(content, projectName string) string {
	cloneSection := fmt.Sprintf(`
## TRex Clone Information
**Generated**: %s  
**Source**: github.com/openshift-online/rh-trex  
**Clone Name**: %s

This project was created using the TRex cloning system. The clone includes:
- ✅ Full CRUD operations for entities
- ✅ OpenAPI specifications
- ✅ Database migrations
- ✅ Service layer architecture
- ✅ Entity generator for instant API development

`, time.Now().Format("January 2, 2006"), projectName)

	return content + cloneSection
}

// Utility functions for name transformations
func toCamelCase(s string) string {
	if s == "" {
		return ""
	}

	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_'
	})

	result := ""
	for _, part := range parts {
		if len(part) > 0 {
			result += strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return result
}

func toSqlSafeName(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}

func toBinaryName(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "-", ""), "_", "")
}
