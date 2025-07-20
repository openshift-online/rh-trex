package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/openshift-online/rh-trex/pkg/config"
)

type CloneConfig struct {
	Name        string
	Repo        string
	Destination string
}

var rw os.FileMode = 0777

func main() {
	// Parse command line flags
	cloneCfg := &CloneConfig{}
	flag.StringVar(&cloneCfg.Name, "name", "rh-trex", "Name of the new service being provisioned")
	flag.StringVar(&cloneCfg.Destination, "destination", "/tmp/clone-test", "Target directory for the newly provisioned instance")
	flag.StringVar(&cloneCfg.Repo, "repo", "github.com/openshift-online", "git repo of project")
	flag.Parse()

	// Always log to stderr by default
	if err := flag.Set("logtostderr", "true"); err != nil {
		glog.Infof("Unable to set logtostderr to true")
	}

	if cloneCfg.Name == "" {
		glog.Fatalf("--name is required")
	}
	if cloneCfg.Destination == "" {
		glog.Fatalf("--destination is required")
	}

	if err := cloneProject(cloneCfg); err != nil {
		glog.Fatalf("Clone failed: %v", err)
	}

	glog.Infof("Clone completed successfully!")
}

func cloneProject(cloneCfg *CloneConfig) error {
	glog.Infof("creating new TRex instance as %s in directory %s", cloneCfg.Name, cloneCfg.Destination)

	// Ensure the destination base directory exists
	if err := os.MkdirAll(cloneCfg.Destination, rw); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %v", cloneCfg.Destination, err)
	}

	// walk the filesystem, starting at the root of the project
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ignore git subdirectories
		if path == ".git" || strings.Contains(path, ".git/") {
			return nil
		}

		// Skip clone command to prevent self-corruption
		if strings.Contains(path, "cmd/trex/clone/") || strings.Contains(path, "/clone/cmd.go") {
			return nil
		}

		// Skip cloner.go but keep generator.go for entity generation in clones
		if path == "scripts/cloner.go" {
			return nil
		}

		dest := cloneCfg.Destination + "/" + path
		if strings.Contains(dest, "trex") {
			dest = strings.Replace(dest, "trex", strings.ToLower(cloneCfg.Name), -1)
		}

		if info.IsDir() {
			// does this path exist in the destination?
			if _, err := os.Stat(dest); os.IsNotExist(err) {
				glog.Infof("Directory does not exist, creating: %s", dest)
			}

			err := os.MkdirAll(dest, rw)
			if err != nil {
				return err
			}

		} else {
			content, err := config.ReadFile(path)
			if err != nil {
				return err
			}

			// Special handling for CLAUDE.md to add TRex clone information
			if strings.HasSuffix(path, "CLAUDE.md") {
				content = addTRexCloneSection(content, cloneCfg.Name)
			}

			if strings.Contains(content, "github.com/openshift-online/rh-trex/pkg/") {
				glog.Infof("find/replace required for file: %s", path)
				replacement := fmt.Sprintf("%s/%s", cloneCfg.Repo, strings.ToLower(cloneCfg.Name))
				// Replace specific rh-trex package imports, preserving rh-trex-core
				content = strings.Replace(content, "github.com/openshift-online/rh-trex/pkg/", replacement+"/pkg/", -1)
				content = strings.Replace(content, "github.com/openshift-online/rh-trex/cmd/", replacement+"/cmd/", -1)
			}

			if strings.Contains(content, "RHTrex") {
				glog.Infof("find/replace required for file: %s", path)
				content = strings.Replace(content, "RHTrex", cloneCfg.Name, -1)
			}

			if strings.Contains(content, "rh-trex") && !strings.Contains(content, "github.com/openshift-online/rh-trex-core") {
				glog.Infof("find/replace required for file: %s", path)
				// Use line-by-line replacement to preserve rh-trex-core dependencies
				lines := strings.Split(content, "\n")
				for i, line := range lines {
					if strings.Contains(line, "rh-trex") && !strings.Contains(line, "rh-trex-core") {
						lines[i] = strings.Replace(line, "rh-trex", strings.ToLower(cloneCfg.Name), -1)
					}
				}
				content = strings.Join(lines, "\n")
			}

			if strings.Contains(content, "rhtrex") {
				glog.Infof("find/replace required for file: %s", path)
				// Use line-by-line replacement to handle database variables with SQL-safe names
				lines := strings.Split(content, "\n")
				for i, line := range lines {
					if strings.Contains(line, "rhtrex") {
						if strings.Contains(line, "db_name:=") || strings.Contains(line, "db_user:=") || 
						   strings.Contains(line, "POSTGRES_DB=") || strings.Contains(line, "POSTGRES_USER=") {
							// Use SQL-safe name for database variables
							lines[i] = strings.Replace(line, "rhtrex", toSqlSafeName(strings.ToLower(cloneCfg.Name)), -1)
						} else {
							// Use regular name for other variables
							lines[i] = strings.Replace(line, "rhtrex", strings.ToLower(cloneCfg.Name), -1)
						}
					}
				}
				content = strings.Join(lines, "\n")
			}

			if strings.Contains(content, "trex") && !strings.Contains(content, "rh-trex-core") {
				glog.Infof("find/replace required for file: %s", path)
				// Use line-by-line replacement to preserve rh-trex-core dependencies
				lines := strings.Split(content, "\n")
				for i, line := range lines {
					if strings.Contains(line, "trex") && !strings.Contains(line, "rh-trex-core") {
						if strings.Contains(line, "db_user:=") {
							// Use SQL-safe name for database variables
							lines[i] = strings.Replace(line, "trex", toSqlSafeName(strings.ToLower(cloneCfg.Name)), -1)
						} else {
							// Use regular name for other variables
							lines[i] = strings.Replace(line, "trex", strings.ToLower(cloneCfg.Name), -1)
						}
					}
				}
				content = strings.Join(lines, "\n")
			}

			if strings.Contains(content, "TRex") {
				glog.Infof("find/replace required for file: %s", path)
				content = strings.Replace(content, "TRex", cloneCfg.Name, -1)
			}

			// 1. Go module declaration replacement
			if strings.HasPrefix(content, "module github.com/openshift-online/rh-trex") {
				glog.Infof("find/replace required for file: %s", path)
				content = strings.Replace(content, "module github.com/openshift-online/rh-trex", 
					fmt.Sprintf("module %s/%s", cloneCfg.Repo, strings.ToLower(cloneCfg.Name)), 1)
			}

			// 2. API URL paths replacement
			if strings.Contains(content, "/api/rh-trex/v1/") {
				glog.Infof("find/replace required for file: %s", path)
				content = strings.ReplaceAll(content, "/api/rh-trex/v1/", 
					fmt.Sprintf("/api/%s/v1/", strings.ToLower(cloneCfg.Name)))
			}

			// 3. OpenAPI client method names replacement
			if strings.Contains(content, "ApiRhTrexV1") {
				glog.Infof("find/replace required for file: %s", path)
				content = strings.ReplaceAll(content, "ApiRhTrexV1", 
					fmt.Sprintf("Api%sV1", toCamelCase(cloneCfg.Name)))
			}

			// 4. Database container names replacement
			if strings.Contains(content, "psql-rhtrex") {
				glog.Infof("find/replace required for file: %s", path)
				content = strings.ReplaceAll(content, "psql-rhtrex", 
					fmt.Sprintf("psql-%s", toSqlSafeName(strings.ToLower(cloneCfg.Name))))
			}


			// 6. Error code prefixes replacement
			if strings.Contains(content, "ERROR_CODE_PREFIX = \"rh-trex\"") {
				glog.Infof("find/replace required for file: %s", path)
				content = strings.ReplaceAll(content, "ERROR_CODE_PREFIX = \"rh-trex\"", 
					fmt.Sprintf("ERROR_CODE_PREFIX = \"%s\"", strings.ToLower(cloneCfg.Name)))
			}

			if exists(dest) {
				e := os.Remove(dest)
				if e != nil {
					return err
				}
			}

			file, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_RDWR, rw)
			if err != nil {
				return err
			}

			written, fErr := file.WriteString(content)
			if fErr != nil {
				return fErr
			}

			glog.Infof("wrote %d bytes for file %s", written, dest)
			file.Sync()
			file.Close()
		}

		return nil
	})

	return err
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

// toCamelCase converts a string to CamelCase (e.g., "ocm-ai" -> "OcmAi")
func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	
	// Split by hyphens and underscores
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_'
	})
	
	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(part[:1]))
			if len(part) > 1 {
				result.WriteString(strings.ToLower(part[1:]))
			}
		}
	}
	
	return result.String()
}

// toSqlSafeName converts a string to SQL-safe name by replacing hyphens with underscores
func toSqlSafeName(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}

// addTRexCloneSection adds TRex clone information to CLAUDE.md
func addTRexCloneSection(content, projectName string) string {
	cloneSection := fmt.Sprintf(`
## TRex Clone Information

**THIS IS A TREX CLONE** - This project was created from the TRex template framework.

- **Template Source**: [github.com/openshift-online/rh-trex](https://github.com/openshift-online/rh-trex)
- **Clone Name**: %s
- **Template Framework**: TRex provides REST API microservice templates with CRUD operations, authentication, database management, and code generation

### Applying TRex Updates

To apply bug fixes and improvements from the main TRex repository to this clone:

**1. Identify the fixes** in the main TRex repository (typically in generator templates or core functionality)

**2. Apply the same changes** to this clone by comparing files:
   - Generator: scripts/generator.go and templates/ directory
   - Core functionality: Follow TRex patterns for error handling, database operations, etc.

**3. Common update scenarios:**
   - **Generator fixes**: Compare scripts/generator.go with main TRex and apply missing functions/features
   - **Template updates**: Compare templates/ directory files and apply template improvements
   - **Core library integration**: Update to newer versions of rh-trex-core dependency
   - **Build system improvements**: Apply Makefile, CI/CD, or container updates

**4. Testing after updates:**
   - make test                 (Run unit tests)
   - make test-integration     (Run integration tests)
   - go run ./scripts/generator.go --kind TestKind  (Test generator functionality)

**5. Example update process** (like applied to ABE clone):
   - Compare generator files between main TRex and this clone
   - diff /path/to/main/trex/scripts/generator.go ./scripts/generator.go
   - Apply missing functions (e.g., toCamelCase, ProjectCamelCase field)
   - Update templates with dynamic ProjectCamelCase variables
   - Test the changes: make test && make test-integration

For systematic updates, use this checklist:
- [ ] Compare scripts/generator.go with main TRex
- [ ] Compare templates/ directory contents 
- [ ] Check for new rh-trex-core library versions
- [ ] Verify all tests pass after applying changes
- [ ] Test code generation with a sample Kind

`, strings.ToUpper(projectName))

	// Insert the clone section after the first # header (after "# CLAUDE.md")
	lines := strings.Split(content, "\n")
	var result []string

	headerFound := false
	sectionInserted := false

	for _, line := range lines {
		result = append(result, line)

		// Insert clone section after the first header and its description
		if !headerFound && strings.HasPrefix(line, "# ") {
			headerFound = true
		} else if headerFound && !sectionInserted && (strings.HasPrefix(line, "## ") || (strings.TrimSpace(line) == "" && len(result) > 3)) {
			// Insert before the next section or after a blank line following the description
			if strings.HasPrefix(line, "## ") {
				// Insert before this section
				result = result[:len(result)-1] // Remove the current line
				result = append(result, strings.Split(cloneSection, "\n")...)
				result = append(result, line) // Add back the current line
			} else if strings.TrimSpace(line) == "" {
				// Insert after blank line
				result = append(result, strings.Split(cloneSection, "\n")...)
			}
			sectionInserted = true
		}
	}

	// If we haven't inserted yet, append at the end
	if !sectionInserted {
		result = append(result, strings.Split(cloneSection, "\n")...)
	}

	return strings.Join(result, "\n")
}