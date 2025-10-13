package clone

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/openshift-online/rh-trex/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

		dest := provisionCfg.Destination + "/" + path
		dest = strings.Replace(dest, "trex", strings.ToLower(provisionCfg.Name), -1)

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

			replacement := fmt.Sprintf("%s/%s", provisionCfg.Repo, strings.ToLower(provisionCfg.Name))
			content = strings.Replace(content, "github.com/openshift-online/rh-trex", replacement, -1)
			content = strings.Replace(content, "RHTrex", provisionCfg.Name, -1)
			content = strings.Replace(content, "rh-trex", strings.ToLower(provisionCfg.Name), -1)
			content = strings.Replace(content, "rhtrex", strings.ToLower(provisionCfg.Name), -1)
			content = strings.Replace(content, "trex", strings.ToLower(provisionCfg.Name), -1)
			content = strings.Replace(content, "TRex", provisionCfg.Name, -1)

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
		return
	}

	// Print next steps for the customer
	serviceName := strings.ToLower(provisionCfg.Name)
	msg := fmt.Sprintf(`
✅ Clone completed successfully!

📋 Next steps to run your new service:

1. Navigate to your new service directory:
	cd %s

2. Install dependencies:
	go mod tidy

3. Build the project:
	go install gotest.tools/gotestsum@latest
	make binary

4. Set up the database:
	make db/setup

5. Run database migrations:
	./%s migrate

6. Test the application:
	make test
	make test-integration

7. Run your service (choose one option):

	Option A: Without authentication (recommended for local development):
	make run-no-auth

	Option B: With authentication (production-like):
	make run

8. Verify the service is running:

	If using Option A (no auth):
	curl http://localhost:8000/api/%s/v1/dinosaurs | jq

	If using Option B (with auth):
	ocm login --token=${OCM_ACCESS_TOKEN} --url=http://localhost:8000
	ocm get /api/%s/v1/dinosaurs

For more detailed information, refer to the README.md in your new service directory.
`, provisionCfg.Destination, serviceName, serviceName, serviceName)

	fmt.Println(msg)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
