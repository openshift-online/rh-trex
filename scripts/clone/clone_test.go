package main

import (
	"os"
	"strings"
	"testing"
)

// TestCloneFileCategories tests the file categorization logic with real files
func TestCloneFileCategories(t *testing.T) {
	tests := []struct {
		path     string
		expected FileCategory
	}{
		{"../../go.mod", ModuleFile},
		{"../../pkg/services/dinosaurs.go", GoSourceFile},
		{"../../cmd/trex/main.go", GoSourceFile},
		{"../../openapi/openapi.yaml", OpenAPIFile},
		{"../../openapi/openapi.dinosaurs.yaml", OpenAPIFile},
		{"../../Makefile", InfrastructureFile},
		{"../../Dockerfile", InfrastructureFile},
		{"../../templates/service-template.yml", InfrastructureFile},
		{"../../README.md", DocumentationFile},
		{"../../CLAUDE.md", DocumentationFile},
		{"../../pkg/config/config.yaml", ConfigurationFile},
		{"../some-unknown-file.txt", SkipFile},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := categorizeFile(tt.path)
			if result != tt.expected {
				t.Errorf("CategorizeFile(%q) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

// TestProcessModuleFile tests go.mod processing with real go.mod
func TestProcessModuleFile(t *testing.T) {
	config := &CloneConfig{
		Name: "user-service",
		Repo: "github.com/openshift-online",
	}

	// Read real go.mod file
	content, err := os.ReadFile("../../go.mod")
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}
	
	originalContent := string(content)
	result := processModuleFile(originalContent, config)

	// Verify ONLY the module declaration was changed
	if !strings.Contains(result, "module github.com/openshift-online/user-service") {
		t.Error("Module declaration was not updated correctly")
	}

	// Verify rh-trex-core dependencies are preserved
	if strings.Contains(originalContent, "rh-trex-core") && !strings.Contains(result, "rh-trex-core") {
		t.Error("rh-trex-core dependency was incorrectly replaced")
	}

	// Verify no other unwanted replacements
	originalLines := strings.Split(originalContent, "\n")
	resultLines := strings.Split(result, "\n")
	
	for i, line := range originalLines {
		if i >= len(resultLines) {
			break
		}
		if strings.HasPrefix(line, "module github.com/openshift-online/rh-trex") {
			// This line should be changed
			continue
		}
		if line != resultLines[i] {
			t.Errorf("Unexpected change in line %d:\noriginal: %q\nresult: %q", i+1, line, resultLines[i])
		}
	}
}

// TestProcessGoSourceFile tests Go source file processing with real Go file
func TestProcessGoSourceFile(t *testing.T) {
	config := &CloneConfig{
		Name: "user-service",
		Repo: "github.com/openshift-online",
	}

	// Read real Go source file
	content, err := os.ReadFile("../../pkg/services/dinosaurs.go")
	if err != nil {
		t.Fatalf("Failed to read dinosaurs.go: %v", err)
	}
	
	originalContent := string(content)
	result := processGoSourceFile(originalContent, config)

	// Verify import paths were updated
	if strings.Contains(originalContent, "github.com/openshift-online/rh-trex/") {
		if !strings.Contains(result, "github.com/openshift-online/user-service/") {
			t.Error("Import paths were not updated correctly")
		}
		if strings.Contains(result, "github.com/openshift-online/rh-trex/") {
			t.Error("Some import paths were not replaced")
		}
	}

	// Verify rh-trex-core imports are preserved
	if strings.Contains(originalContent, "rh-trex-core") && !strings.Contains(result, "rh-trex-core") {
		t.Error("rh-trex-core imports were incorrectly replaced")
	}

	// Count occurrences to ensure complete replacement
	originalRhTrexCount := strings.Count(originalContent, "github.com/openshift-online/rh-trex/")
	resultRhTrexCount := strings.Count(result, "github.com/openshift-online/rh-trex/")
	resultUserServiceCount := strings.Count(result, "github.com/openshift-online/user-service/")
	
	if resultRhTrexCount > 0 {
		t.Errorf("Found %d unreplaced rh-trex import paths", resultRhTrexCount)
	}
	if originalRhTrexCount > 0 && resultUserServiceCount == 0 {
		t.Error("No user-service import paths found after replacement")
	}
}

// TestProcessOpenAPIFile tests OpenAPI file processing with real OpenAPI file
func TestProcessOpenAPIFile(t *testing.T) {
	config := &CloneConfig{
		Name: "user-service",
		Repo: "github.com/openshift-online",
	}

	// Read real OpenAPI file
	content, err := os.ReadFile("../../openapi/openapi.yaml")
	if err != nil {
		t.Fatalf("Failed to read openapi.yaml: %v", err)
	}
	
	originalContent := string(content)
	result := processOpenAPIFile(originalContent, config)

	// Verify API paths were updated
	if strings.Contains(originalContent, "/api/rh-trex/") {
		if !strings.Contains(result, "/api/user-service/") {
			t.Error("API paths were not updated correctly")
		}
		if strings.Contains(result, "/api/rh-trex/") {
			t.Error("Some API paths were not replaced")
		}
	}

	// Verify operation IDs were updated
	if strings.Contains(originalContent, "ApiRhTrexV1") {
		if !strings.Contains(result, "ApiUserServiceV1") {
			t.Error("Operation IDs were not updated correctly")
		}
		if strings.Contains(result, "ApiRhTrexV1") {
			t.Error("Some operation IDs were not replaced")
		}
	}

	// Verify service descriptions were updated
	if strings.Contains(originalContent, "rh-trex Service API") {
		if !strings.Contains(result, "user-service Service API") {
			t.Error("Service descriptions were not updated correctly")
		}
		if strings.Contains(result, "rh-trex Service API") {
			t.Error("Some service descriptions were not replaced")
		}
	}

	// Verify no unintended changes in structure
	originalPathCount := strings.Count(originalContent, "paths:")
	resultPathCount := strings.Count(result, "paths:")
	if originalPathCount != resultPathCount {
		t.Error("OpenAPI structure was modified unexpectedly")
	}
}

// TestProcessInfrastructureFile tests infrastructure file processing with real Makefile
func TestProcessInfrastructureFile(t *testing.T) {
	config := &CloneConfig{
		Name: "user-service",
		Repo: "github.com/openshift-online",
	}

	// Read real Makefile
	content, err := os.ReadFile("../../Makefile")
	if err != nil {
		t.Fatalf("Failed to read Makefile: %v", err)
	}
	
	originalContent := string(content)
	result := processInfrastructureFile(originalContent, config)

	// Verify binary name replacements (should be userservice, no hyphens)
	expectedBinaryName := "userservice"
	
	// Check for binary path replacements
	if strings.Contains(originalContent, "/usr/local/bin/trex") {
		if !strings.Contains(result, "/usr/local/bin/"+expectedBinaryName) {
			t.Errorf("Binary paths were not updated to %s", expectedBinaryName)
		}
	}

	// Check for service name replacements  
	if strings.Contains(originalContent, "trex-service") {
		if !strings.Contains(result, expectedBinaryName+"-service") {
			t.Errorf("Service names were not updated correctly")
		}
	}

	// Verify database name replacements (should be user_service, SQL-safe)
	expectedSqlName := "user_service"
	if strings.Contains(originalContent, "rhtrex") {
		if !strings.Contains(result, expectedSqlName) {
			t.Errorf("Database names were not updated to %s", expectedSqlName)
		}
	}

	// Verify rh-trex-core references are preserved
	if strings.Contains(originalContent, "rh-trex-core") && !strings.Contains(result, "rh-trex-core") {
		t.Error("rh-trex-core references were incorrectly replaced")
	}

	// Count lines to ensure structure is preserved
	originalLineCount := len(strings.Split(originalContent, "\n"))
	resultLineCount := len(strings.Split(result, "\n"))
	if originalLineCount != resultLineCount {
		t.Error("Makefile structure was modified (different line count)")
	}
}

// TestProcessDocumentationFile tests documentation processing with real README
func TestProcessDocumentationFile(t *testing.T) {
	config := &CloneConfig{
		Name: "user-service",
		Repo: "github.com/openshift-online",
	}

	// Read real README file
	content, err := os.ReadFile("../../README.md")
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}
	
	originalContent := string(content)
	result := processDocumentationFile(originalContent, "README.md", config)

	// Verify project name replacements
	if strings.Contains(originalContent, "rh-trex Service") {
		if !strings.Contains(result, "user-service Service") {
			t.Error("Project service names were not updated correctly")
		}
		if strings.Contains(result, "rh-trex Service") {
			t.Error("Some project service names were not replaced")
		}
	}

	if strings.Contains(originalContent, "TRex service") {
		if !strings.Contains(result, "user-service service") {
			t.Error("TRex service references were not updated correctly")
		}
		if strings.Contains(result, "TRex service") {
			t.Error("Some TRex service references were not replaced")
		}
	}

	// Verify template references were updated
	if strings.Contains(originalContent, "rh-trex template") {
		if !strings.Contains(result, "user-service template") {
			t.Error("Template references were not updated correctly")
		}
	}

	// Verify source references to rh-trex repo are preserved
	// (These should NOT be replaced as they refer to the template source)
	sourceReferences := []string{
		"github.com/openshift-online/rh-trex",
		"Template Source:",
	}
	
	for _, ref := range sourceReferences {
		if strings.Contains(originalContent, ref) && !strings.Contains(result, ref) {
			t.Errorf("Source reference '%s' was incorrectly replaced", ref)
		}
	}
}

// TestProcessCLAUDEFile tests CLAUDE.md gets clone section added
func TestProcessCLAUDEFile(t *testing.T) {
	config := &CloneConfig{
		Name: "user-service",
		Repo: "github.com/openshift-online",
	}

	// Read real CLAUDE.md file
	content, err := os.ReadFile("../../CLAUDE.md")
	if err != nil {
		t.Fatalf("Failed to read CLAUDE.md: %v", err)
	}
	
	originalContent := string(content)
	result := processDocumentationFile(originalContent, "CLAUDE.md", config)

	// Verify clone section was added
	if !strings.Contains(result, "## TRex Clone Information") {
		t.Error("Clone section was not added to CLAUDE.md")
	}

	if !strings.Contains(result, "Clone Name**: user-service") {
		t.Error("Clone name was not added correctly to CLAUDE.md")
	}

	if !strings.Contains(result, "github.com/openshift-online/rh-trex") {
		t.Error("Source repository reference was not preserved in clone section")
	}

	// Verify original content is still present
	if len(result) <= len(originalContent) {
		t.Error("Clone section does not appear to have been added (no length increase)")
	}
}

// TestNameTransformationFunctions tests utility functions with various inputs
func TestNameTransformationFunctions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		fn       func(string) string
	}{
		// ToCamelCase tests
		{"toCamelCase hyphenated", "user-service", "UserService", toCamelCase},
		{"toCamelCase underscored", "user_service", "UserService", toCamelCase},
		{"toCamelCase mixed", "ocm-ai_service", "OcmAiService", toCamelCase},
		{"toCamelCase single", "service", "Service", toCamelCase},
		{"toCamelCase empty", "", "", toCamelCase},
		
		// toSqlSafeName tests
		{"toSqlSafeName basic", "user-service", "user_service", toSqlSafeName},
		{"toSqlSafeName multiple", "my-test-api", "my_test_api", toSqlSafeName},
		{"toSqlSafeName unchanged", "service", "service", toSqlSafeName},
		{"toSqlSafeName empty", "", "", toSqlSafeName},
		
		// toBinaryName tests  
		{"toBinaryName hyphenated", "user-service", "userservice", toBinaryName},
		{"toBinaryName underscored", "user_service", "userservice", toBinaryName},
		{"toBinaryName mixed", "my-test_api", "mytestapi", toBinaryName},
		{"toBinaryName unchanged", "service", "service", toBinaryName},
		{"toBinaryName empty", "", "", toBinaryName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.input)
			if result != tt.expected {
				t.Errorf("%s(%q) = %q, expected %q", tt.name, tt.input, result, tt.expected)
			}
		})
	}
}

// TestNoUnintendedReplacements ensures framework functions are not replaced
func TestNoUnintendedReplacements(t *testing.T) {
	config := &CloneConfig{
		Name: "user-service",
		Repo: "github.com/openshift-online",
	}

	// Test content with framework elements that should NOT be replaced
	testContent := `
	import "github.com/openshift-online/rh-trex-core/pkg/auth"
	import "github.com/openshift-online/rh-trex/pkg/api"
	
	func addTRexCloneSection() {
		// This function name should not change
	}
	
	// TRex-generated code marker
	type TRexFrameworkType struct{}
	`

	result := processGoSourceFile(testContent, config)

	// Verify rh-trex-core is preserved
	if !strings.Contains(result, "rh-trex-core") {
		t.Error("rh-trex-core import was incorrectly replaced")
	}

	// Verify project import was replaced
	if !strings.Contains(result, "github.com/openshift-online/user-service/pkg/api") {
		t.Error("Project import was not replaced")
	}

	// Verify function names are preserved
	if !strings.Contains(result, "addTRexCloneSection") {
		t.Error("Framework function name was incorrectly replaced")
	}

	// Verify generation markers are preserved
	if !strings.Contains(result, "TRex-generated") {
		t.Error("Generation marker was incorrectly replaced")
	}
}