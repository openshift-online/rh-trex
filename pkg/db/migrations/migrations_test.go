package migrations

import (
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMigration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Migration Suite")
}

var _ = Describe("Migrate", func() {
	It("Expects same amount of files and migrationList", func() {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		files, err := os.ReadDir(cwd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		amountGoFiles := []string{}
		for _, file := range files {
			if !strings.Contains(file.Name(), ".go") {
				continue
			}
			if strings.Contains(file.Name(), "migration_structs") || strings.Contains(file.Name(), "migrations_test") {
				continue
			}
			amountGoFiles = append(amountGoFiles, file.Name())
		}
		// Disconsiders migration_structs.go and test files
		Expect(amountGoFiles).To(HaveLen(len(MigrationList)))
	})
})
