package clone

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/openshift-online/rh-trex/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
	"strings"
)

type provisionCfgFlags struct {
	Name        string
	Repo        string
	Destination string
}

func (c *provisionCfgFlags) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.Name, "name", c.Name, "Name of the new service being provisioned")
	fs.StringVar(&c.Destination, "destination", c.Destination, "Target directory for the newly provisioned instance")
	fs.StringVar(&c.Repo, "repo", c.Repo, "git repo of project")
}

var provisionCfg = &provisionCfgFlags{
	Name:        "rh-trex",
	Repo:        "github.com/openshift-online",
	Destination: "/tmp/clone-test",
}

// migrate sub-command handles running migrations
func NewCloneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone a new TRex instance",
		Long:  "Clone a new TRex instance",
		Run:   clone,
	}

	provisionCfg.AddFlags(cmd.PersistentFlags())
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	return cmd
}

var rw os.FileMode = 0777

func clone(_ *cobra.Command, _ []string) {

	glog.Infof("creating new TRex instance as %s in directory %s", provisionCfg.Name, provisionCfg.Destination)

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

		dest := provisionCfg.Destination + "/" + path
		if strings.Contains(dest, "trex") {
			dest = strings.Replace(dest, "trex", strings.ToLower(provisionCfg.Name), -1)
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
				content = addTRexCloneSection(content, provisionCfg.Name)
			}

			if strings.Contains(content, "github.com/openshift-online/rh-trex/pkg/") {
				glog.Infof("find/replace required for file: %s", path)
				replacement := fmt.Sprintf("%s/%s", provisionCfg.Repo, strings.ToLower(provisionCfg.Name))
				// Replace specific rh-trex package imports, preserving rh-trex-core
				content = strings.Replace(content, "github.com/openshift-online/rh-trex/pkg/", replacement+"/pkg/", -1)
				content = strings.Replace(content, "github.com/openshift-online/rh-trex/cmd/", replacement+"/cmd/", -1)
			}

			if strings.Contains(content, "RHTrex") {
				glog.Infof("find/replace required for file: %s", path)
				content = strings.Replace(content, "RHTrex", provisionCfg.Name, -1)
			}

			if strings.Contains(content, "rh-trex") && !strings.Contains(content, "github.com/openshift-online/rh-trex-core") {
				glog.Infof("find/replace required for file: %s", path)
				// Use line-by-line replacement to preserve rh-trex-core dependencies
				lines := strings.Split(content, "\n")
				for i, line := range lines {
					if strings.Contains(line, "rh-trex") && !strings.Contains(line, "rh-trex-core") {
						lines[i] = strings.Replace(line, "rh-trex", strings.ToLower(provisionCfg.Name), -1)
					}
				}
				content = strings.Join(lines, "\n")
			}

			if strings.Contains(content, "rhtrex") {
				glog.Infof("find/replace required for file: %s", path)
				content = strings.Replace(content, "rhtrex", strings.ToLower(provisionCfg.Name), -1)
			}

			if strings.Contains(content, "trex") && !strings.Contains(content, "rh-trex-core") {
				glog.Infof("find/replace required for file: %s", path)
				// Use line-by-line replacement to preserve rh-trex-core dependencies
				lines := strings.Split(content, "\n")
				for i, line := range lines {
					if strings.Contains(line, "trex") && !strings.Contains(line, "rh-trex-core") {
						lines[i] = strings.Replace(line, "trex", strings.ToLower(provisionCfg.Name), -1)
					}
				}
				content = strings.Join(lines, "\n")
			}

			if strings.Contains(content, "TRex") {
				glog.Infof("find/replace required for file: %s", path)
				content = strings.Replace(content, "TRex", provisionCfg.Name, -1)
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

	if err != nil {
		fmt.Println(err)
	}

}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
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
