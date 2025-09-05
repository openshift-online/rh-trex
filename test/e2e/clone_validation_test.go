package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// CloneValidationTest represents the end-to-end test structure
type CloneValidationTest struct {
	t              *testing.T
	tempDir        string
	cloneDir       string
	projectName    string
	containers     []string
	processes      []*os.Process
	cleanupEnabled bool
}

// NewCloneValidationTest creates a new instance of the e2e test
func NewCloneValidationTest(t *testing.T) *CloneValidationTest {
	tempDir, err := os.MkdirTemp("", "trex-e2e-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	projectName := fmt.Sprintf("test-service-%d", time.Now().Unix())
	cloneDir := filepath.Join(tempDir, projectName)

	test := &CloneValidationTest{
		t:              t,
		tempDir:        tempDir,
		cloneDir:       cloneDir,
		projectName:    projectName,
		containers:     make([]string, 0),
		processes:      make([]*os.Process, 0),
		cleanupEnabled: true,
	}

	// Ensure cleanup always runs
	t.Cleanup(func() {
		test.cleanup()
	})

	return test
}

// TestCloneValidation is the main e2e test function
func TestCloneValidation(t *testing.T) {
	test := NewCloneValidationTest(t)

	// Step 1: Clone Process Validation
	if !test.runStep1CloneProcessValidation() {
		t.Fatal("Step 1: Clone Process Validation failed")
	}

	// TODO: Add remaining steps (2-5) in future iterations
	t.Log("E2E test completed successfully - Steps 2-5 will be implemented next")
}

// runStep1CloneProcessValidation implements Step 1: Clone Process Validation
func (cvt *CloneValidationTest) runStep1CloneProcessValidation() bool {
	cvt.t.Log("Starting Step 1: Clone Process Validation")

	// Execute clone command from project root
	projectRoot, err := cvt.findProjectRoot()
	if err != nil {
		cvt.t.Errorf("Failed to find project root: %v", err)
		return false
	}
	
	cloneCmd := exec.Command("go", "run", "./scripts/clone/main.go",
		"--name", cvt.projectName,
		"--destination", cvt.cloneDir,
	)
	cloneCmd.Dir = projectRoot
	
	cvt.t.Logf("Running clone command: %s", strings.Join(cloneCmd.Args, " "))
	output, err := cloneCmd.CombinedOutput()
	if err != nil {
		cvt.t.Errorf("Clone command failed: %v\nOutput: %s", err, string(output))
		return false
	}

	cvt.t.Logf("Clone command output: %s", string(output))

	// Verify all expected files/directories are created
	if !cvt.validateClonedStructure() {
		return false
	}

	// Validate file content transformations
	if !cvt.validateContentTransformations() {
		return false
	}

	cvt.t.Log("Step 1: Clone Process Validation completed successfully")
	return true
}

// validateClonedStructure verifies that all expected files and directories exist
func (cvt *CloneValidationTest) validateClonedStructure() bool {
	cvt.t.Log("Validating cloned directory structure")

	expectedPaths := []string{
		"go.mod",
		"go.sum", 
		"Makefile",
		"cmd",
		"pkg",
		"scripts",
		"test",
		"docs",
		"openapi",
		"plugins",
		"secrets",
		"templates",
	}

	for _, path := range expectedPaths {
		fullPath := filepath.Join(cvt.cloneDir, path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			cvt.t.Errorf("Expected path does not exist: %s", fullPath)
			return false
		}
	}

	cvt.t.Log("Directory structure validation passed")
	return true
}

// validateContentTransformations verifies that module names, import paths, and API paths are correctly transformed
func (cvt *CloneValidationTest) validateContentTransformations() bool {
	cvt.t.Log("Validating content transformations")

	// Check go.mod file for correct module name
	goModPath := filepath.Join(cvt.cloneDir, "go.mod")
	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		cvt.t.Errorf("Failed to read go.mod: %v", err)
		return false
	}

	expectedModuleLine := fmt.Sprintf("module github.com/openshift-online/%s", cvt.projectName)
	if !strings.Contains(string(goModContent), expectedModuleLine) {
		cvt.t.Errorf("go.mod does not contain expected module name. Expected: %s, Content: %s", 
			expectedModuleLine, string(goModContent))
		return false
	}

	// Check main.go for proper imports
	mainGoPath := filepath.Join(cvt.cloneDir, "cmd", "trex", "main.go")
	if _, err := os.Stat(mainGoPath); err == nil {
		mainGoContent, err := os.ReadFile(mainGoPath)
		if err != nil {
			cvt.t.Errorf("Failed to read main.go: %v", err)
			return false
		}

		// Verify that rh-trex import paths have been replaced with the new project name
		expectedImportPath := fmt.Sprintf("github.com/openshift-online/%s/", cvt.projectName)
		if !strings.Contains(string(mainGoContent), expectedImportPath) {
			cvt.t.Errorf("main.go does not contain expected import path %s: %s", expectedImportPath, string(mainGoContent))
			return false
		}
		
		// Verify that old rh-trex import paths are gone
		oldImportPath := "github.com/openshift-online/rh-trex/"
		if strings.Contains(string(mainGoContent), oldImportPath) {
			cvt.t.Errorf("main.go still contains old import path %s: %s", oldImportPath, string(mainGoContent))
			return false
		}
	}

	cvt.t.Log("Content transformations validation passed")
	return true
}

// addContainer tracks a container for cleanup
func (cvt *CloneValidationTest) addContainer(containerID string) {
	cvt.containers = append(cvt.containers, containerID)
}

// addProcess tracks a process for cleanup
func (cvt *CloneValidationTest) addProcess(process *os.Process) {
	cvt.processes = append(cvt.processes, process)
}

// disableCleanup disables automatic cleanup (useful for debugging)
func (cvt *CloneValidationTest) disableCleanup() {
	cvt.cleanupEnabled = false
}

// cleanup implements Step 6: Cleanup
func (cvt *CloneValidationTest) cleanup() {
	if !cvt.cleanupEnabled {
		cvt.t.Log("Cleanup disabled - resources left for debugging")
		cvt.t.Logf("Temp directory: %s", cvt.tempDir)
		return
	}

	cvt.t.Log("Starting Step 6: Cleanup")

	// Stop and remove all podman containers created during test
	for _, containerID := range cvt.containers {
		cvt.t.Logf("Stopping container: %s", containerID)
		
		// Stop container
		stopCmd := exec.Command("podman", "stop", containerID)
		if output, err := stopCmd.CombinedOutput(); err != nil {
			cvt.t.Logf("Warning: Failed to stop container %s: %v, output: %s", 
				containerID, err, string(output))
		}

		// Remove container
		rmCmd := exec.Command("podman", "rm", containerID)
		if output, err := rmCmd.CombinedOutput(); err != nil {
			cvt.t.Logf("Warning: Failed to remove container %s: %v, output: %s", 
				containerID, err, string(output))
		}
	}

	// Clean up any background processes
	for _, process := range cvt.processes {
		cvt.t.Logf("Terminating process: %d", process.Pid)
		
		// Try graceful termination first
		if err := process.Signal(os.Interrupt); err != nil {
			cvt.t.Logf("Warning: Failed to send interrupt to process %d: %v", process.Pid, err)
		}

		// Wait briefly for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		done := make(chan error, 1)
		go func() {
			_, err := process.Wait()
			done <- err
		}()

		select {
		case <-ctx.Done():
			// Force kill if graceful shutdown failed
			cvt.t.Logf("Force killing process: %d", process.Pid)
			if err := process.Kill(); err != nil {
				cvt.t.Logf("Warning: Failed to kill process %d: %v", process.Pid, err)
			}
		case err := <-done:
			if err != nil {
				cvt.t.Logf("Process %d exited with error: %v", process.Pid, err)
			}
		}
		cancel()
	}

	// Remove temporary directories and files with safety checks
	if cvt.tempDir != "" {
		if cvt.isSafeToDelete(cvt.tempDir) {
			cvt.t.Logf("Removing temp directory: %s", cvt.tempDir)
			if err := os.RemoveAll(cvt.tempDir); err != nil {
				cvt.t.Logf("Warning: Failed to remove temp directory %s: %v", cvt.tempDir, err)
			}
		} else {
			cvt.t.Errorf("SAFETY CHECK FAILED: Refusing to delete potentially dangerous path: %s", cvt.tempDir)
		}
	}

	cvt.t.Log("Step 6: Cleanup completed")
}

// findProjectRoot finds the project root by looking for go.mod
func (cvt *CloneValidationTest) findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find go.mod in any parent directory")
}

// Helper method to run commands in the cloned directory
func (cvt *CloneValidationTest) runCommandInClone(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = cvt.cloneDir
	return cmd.CombinedOutput()
}

// Helper method to check if a service is running on a specific port
func (cvt *CloneValidationTest) isServiceRunning(port string) bool {
	cmd := exec.Command("curl", "-f", fmt.Sprintf("http://localhost:%s/api/trex/v1/openapi", port))
	err := cmd.Run()
	return err == nil
}

// isSafeToDelete implements multiple safety checks to prevent accidental deletion of important directories
func (cvt *CloneValidationTest) isSafeToDelete(path string) bool {
	// Convert to absolute path for consistent checking
	absPath, err := filepath.Abs(path)
	if err != nil {
		cvt.t.Logf("Safety check: Failed to get absolute path for %s: %v", path, err)
		return false
	}

	// Critical system directories that should never be deleted
	dangerousPaths := []string{
		"/",
		"/home",
		"/usr",
		"/etc",
		"/var",
		"/opt",
		"/bin",
		"/sbin",
		"/lib",
		"/lib64",
		"/boot",
		"/sys",
		"/proc",
		"/dev",
		"/run",
		"/tmp", // While /tmp is often safe, we want to be extra cautious
	}

	// Check if path is a dangerous system directory
	for _, dangerous := range dangerousPaths {
		if absPath == dangerous {
			cvt.t.Logf("Safety check: Refusing to delete system directory: %s", absPath)
			return false
		}
	}

	// Check if path is too short (likely a root or important directory)
	if len(absPath) < 10 {
		cvt.t.Logf("Safety check: Path too short, likely important: %s", absPath)
		return false
	}

	// Must be in the OS temp directory or have temp-related naming
	tempDir := os.TempDir()
	if !strings.HasPrefix(absPath, tempDir) {
		cvt.t.Logf("Safety check: Path not in OS temp directory (%s): %s", tempDir, absPath)
		return false
	}

	// Must contain our test identifier
	if !strings.Contains(absPath, "trex-e2e-test-") {
		cvt.t.Logf("Safety check: Path does not contain test identifier: %s", absPath)
		return false
	}

	// Verify the directory was created by os.MkdirTemp (has expected naming pattern)
	tempDirName := filepath.Base(absPath)
	if !strings.HasPrefix(tempDirName, "trex-e2e-test-") {
		cvt.t.Logf("Safety check: Directory name doesn't match expected pattern: %s", tempDirName)
		return false
	}

	// Check that it's actually a directory
	info, err := os.Stat(absPath)
	if err != nil {
		cvt.t.Logf("Safety check: Cannot stat path %s: %v", absPath, err)
		return false
	}
	if !info.IsDir() {
		cvt.t.Logf("Safety check: Path is not a directory: %s", absPath)
		return false
	}

	// Additional check: ensure we're not trying to delete the user's home directory
	userHome, err := os.UserHomeDir()
	if err == nil && strings.HasPrefix(userHome, absPath) {
		cvt.t.Logf("Safety check: Path would delete user home directory: %s", absPath)
		return false
	}

	// Check depth - temp test directories should be reasonably deep
	// Allow depth of 2 for /tmp/trex-e2e-test-* directories
	pathParts := strings.Split(strings.Trim(absPath, "/"), "/")
	minDepth := 2
	if !strings.HasPrefix(absPath, "/tmp/") {
		minDepth = 3 // Be more strict for non-/tmp paths
	}
	if len(pathParts) < minDepth {
		cvt.t.Logf("Safety check: Path not deep enough, likely important: %s (depth: %d, min: %d)", absPath, len(pathParts), minDepth)
		return false
	}

	cvt.t.Logf("Safety check: Path is safe to delete: %s", absPath)
	return true
}