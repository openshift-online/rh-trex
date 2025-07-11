package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// ProjectConfig defines the configuration for a new project
type ProjectConfig struct {
	Name           string
	Module         string
	Repository     string
	Author         string
	Description    string
	License        string
	TRexCoreVersion string
	
	// Initial resources to create
	Resources []ResourceConfig
}

// ResourceConfig defines a resource to create in the new project
type ResourceConfig struct {
	Name       string
	Attributes []AttributeConfig
}

// AttributeConfig defines an attribute for a resource
type AttributeConfig struct {
	Name     string
	Type     string
	Required bool
	Index    bool
}

// ProjectTemplate handles creating new projects from templates
type ProjectTemplate struct {
	config    ProjectConfig
	templates map[string]*template.Template
}

// NewProjectTemplate creates a new project template
func NewProjectTemplate(config ProjectConfig) *ProjectTemplate {
	return &ProjectTemplate{
		config:    config,
		templates: make(map[string]*template.Template),
	}
}

// Generate creates a new project from the template
func (pt *ProjectTemplate) Generate(outputDir string) error {
	// Create project directory structure
	if err := pt.createDirectoryStructure(outputDir); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}
	
	// Generate core files
	if err := pt.generateCoreFiles(outputDir); err != nil {
		return fmt.Errorf("failed to generate core files: %w", err)
	}
	
	// Generate resources
	if err := pt.generateResources(outputDir); err != nil {
		return fmt.Errorf("failed to generate resources: %w", err)
	}
	
	// Generate deployment files
	if err := pt.generateDeploymentFiles(outputDir); err != nil {
		return fmt.Errorf("failed to generate deployment files: %w", err)
	}
	
	return nil
}

// createDirectoryStructure creates the basic directory structure
func (pt *ProjectTemplate) createDirectoryStructure(outputDir string) error {
	dirs := []string{
		"cmd/" + pt.config.Name,
		"pkg/config",
		"pkg/resources",
		"pkg/handlers",
		"test/integration",
		"test/factories",
		"deployments",
		"openapi",
	}
	
	for _, dir := range dirs {
		fullPath := filepath.Join(outputDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}
	
	return nil
}

// generateCoreFiles generates the main application files
func (pt *ProjectTemplate) generateCoreFiles(outputDir string) error {
	// Generate go.mod
	if err := pt.generateGoMod(outputDir); err != nil {
		return err
	}
	
	// Generate main.go
	if err := pt.generateMainGo(outputDir); err != nil {
		return err
	}
	
	// Generate Dockerfile
	if err := pt.generateDockerfile(outputDir); err != nil {
		return err
	}
	
	// Generate Makefile
	if err := pt.generateMakefile(outputDir); err != nil {
		return err
	}
	
	return nil
}

// generateGoMod creates the go.mod file
func (pt *ProjectTemplate) generateGoMod(outputDir string) error {
	content := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/openshift-online/rh-trex %s
	github.com/spf13/cobra v1.7.0
	github.com/golang/glog v1.1.0
	gorm.io/gorm v1.25.0
	gorm.io/driver/postgres v1.5.0
)
`, pt.config.Module, pt.config.TRexCoreVersion)
	
	return os.WriteFile(filepath.Join(outputDir, "go.mod"), []byte(content), 0644)
}

// generateMainGo creates the main.go file
func (pt *ProjectTemplate) generateMainGo(outputDir string) error {
	tmpl := `package main

import (
	"context"
	"flag"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/openshift-online/rh-trex/pkg/core/controllers"
	"github.com/openshift-online/rh-trex/pkg/core/generator"
	"github.com/openshift-online/rh-trex/pkg/db"
	"{{.Module}}/pkg/config"
	"{{.Module}}/pkg/resources"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "{{.Name}}",
		Short: "{{.Description}}",
		Long:  "{{.Description}}",
	}

	// Add subcommands
	rootCmd.AddCommand(newServeCommand())
	rootCmd.AddCommand(newMigrateCommand())

	flag.CommandLine.Parse([]string{})
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	if err := rootCmd.Execute(); err != nil {
		glog.Fatalf("Error executing command: %v", err)
	}
}

func newServeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the {{.Name}} server",
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	}
}

func newMigrateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			migrate()
		},
	}
}

func serve() {
	cfg := config.Load()
	
	// Initialize database
	database, err := db.NewDatabase(cfg.Database.ConnectionString)
	if err != nil {
		glog.Fatalf("Failed to initialize database: %v", err)
	}

	// Create controller manager
	lockFactory := db.NewAdvisoryLockFactory(database)
	controllerMgr := controllers.NewControllerManager(lockFactory, nil)

	// Create resource factory
	eventEmitter := generator.NewEventEmitterAdapter(nil) // Wire up your event service
	factory := generator.NewResourceFactory(database.DB, controllerMgr, eventEmitter)

	// Register resources
	{{range .Resources}}
	generator.RegisterResourceType(factory, "{{.Name}}", resources.{{.Name}}{})
	{{end}}

	// Start server
	glog.Infof("Starting {{.Name}} server...")
	
	// Start controller manager
	ctx := context.Background()
	controllerMgr.Start(ctx)
}

func migrate() {
	cfg := config.Load()
	
	// Initialize database
	database, err := db.NewDatabase(cfg.Database.ConnectionString)
	if err != nil {
		glog.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations
	if err := database.AutoMigrate(
		{{range .Resources}}
		&resources.{{.Name}}{},
		{{end}}
	); err != nil {
		glog.Fatalf("Failed to run migrations: %v", err)
	}

	glog.Infof("Migrations completed successfully")
}
`

	t, err := template.New("main").Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(outputDir, "cmd", pt.config.Name, "main.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, pt.config)
}

// generateResources creates resource files
func (pt *ProjectTemplate) generateResources(outputDir string) error {
	for _, resource := range pt.config.Resources {
		if err := pt.generateResource(outputDir, resource); err != nil {
			return err
		}
	}
	return nil
}

// generateResource creates a single resource file
func (pt *ProjectTemplate) generateResource(outputDir string, resource ResourceConfig) error {
	tmpl := `package resources

import (
	"context"

	"github.com/openshift-online/rh-trex/pkg/core/api"
	"github.com/openshift-online/rh-trex/pkg/core/dao"
	"github.com/openshift-online/rh-trex/pkg/core/services"
	"github.com/openshift-online/rh-trex/pkg/logger"
	"gorm.io/gorm"
)

// {{.Name}} represents a {{.Name}} resource
type {{.Name}} struct {
	api.Meta
	{{range .Attributes}}
	{{.Name}} {{.Type}} ` + "`json:\"{{.JSONName}}\"{{if .Index}} gorm:\"index\"{{end}}`" + `
	{{end}}
}

// GetMeta returns the metadata for the {{.Name}}
func (r *{{.Name}}) GetMeta() *api.Meta {
	return &r.Meta
}

// {{.Name}}DAO extends the base DAO with {{.Name}}-specific methods
type {{.Name}}DAO struct {
	*dao.BaseDAO[{{.Name}}]
}

// New{{.Name}}DAO creates a new {{.Name}} DAO
func New{{.Name}}DAO(db *gorm.DB) *{{.Name}}DAO {
	return &{{.Name}}DAO{
		BaseDAO: dao.NewBaseDAO[{{.Name}}](db),
	}
}

// {{.Name}}Service extends the base service with {{.Name}}-specific business logic
type {{.Name}}Service struct {
	*services.BaseCRUDService[{{.Name}}]
	dao    *{{.Name}}DAO
	logger logger.Logger
}

// New{{.Name}}Service creates a new {{.Name}} service
func New{{.Name}}Service(
	dao *{{.Name}}DAO,
	events services.EventEmitter,
) *{{.Name}}Service {
	baseSvc := services.NewBaseCRUDService[{{.Name}}](dao, events, "{{.Name}}s")
	
	return &{{.Name}}Service{
		BaseCRUDService: baseSvc,
		dao:            dao,
		logger:         logger.NewOCMLogger(context.Background()),
	}
}

// processUpsert implements custom business logic for {{.Name}} upsert events
func (s *{{.Name}}Service) processUpsert(ctx context.Context, resource *{{.Name}}) error {
	s.logger.Infof("Processing {{.Name}} upsert: %s", resource.ID)
	
	// Add any {{.Name}}-specific business logic here
	
	return nil
}

// processDelete implements custom business logic for {{.Name}} delete events
func (s *{{.Name}}Service) processDelete(ctx context.Context, id string) error {
	s.logger.Infof("{{.Name}} has been deleted: %s", id)
	
	// Add any {{.Name}}-specific cleanup logic here
	
	return nil
}
`

	// Prepare template data
	templateData := struct {
		Name       string
		Attributes []struct {
			Name     string
			Type     string
			JSONName string
			Index    bool
		}
	}{
		Name: resource.Name,
	}

	for _, attr := range resource.Attributes {
		templateData.Attributes = append(templateData.Attributes, struct {
			Name     string
			Type     string
			JSONName string
			Index    bool
		}{
			Name:     strings.Title(attr.Name),
			Type:     attr.Type,
			JSONName: strings.ToLower(attr.Name),
			Index:    attr.Index,
		})
	}

	t, err := template.New("resource").Parse(tmpl)
	if err != nil {
		return err
	}

	filename := strings.ToLower(resource.Name) + ".go"
	file, err := os.Create(filepath.Join(outputDir, "pkg", "resources", filename))
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, templateData)
}

// generateDockerfile creates a Dockerfile
func (pt *ProjectTemplate) generateDockerfile(outputDir string) error {
	content := fmt.Sprintf(`FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o %s ./cmd/%s

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/%s .

CMD ["./%s", "serve"]
`, pt.config.Name, pt.config.Name, pt.config.Name, pt.config.Name)
	
	return os.WriteFile(filepath.Join(outputDir, "Dockerfile"), []byte(content), 0644)
}

// generateMakefile creates a Makefile
func (pt *ProjectTemplate) generateMakefile(outputDir string) error {
	content := fmt.Sprintf(`.PHONY: build test run migrate

build:
	go build -o %s ./cmd/%s

test:
	go test ./...

run:
	go run ./cmd/%s serve

migrate:
	go run ./cmd/%s migrate

clean:
	rm -f %s
`, pt.config.Name, pt.config.Name, pt.config.Name, pt.config.Name, pt.config.Name)
	
	return os.WriteFile(filepath.Join(outputDir, "Makefile"), []byte(content), 0644)
}

// generateDeploymentFiles creates deployment manifests
func (pt *ProjectTemplate) generateDeploymentFiles(outputDir string) error {
	// Generate docker-compose.yml for local development
	content := fmt.Sprintf(`version: '3.8'
services:
  %s:
    build: .
    ports:
      - "8000:8000"
    depends_on:
      - postgres
    environment:
      - DATABASE_URL=postgres://user:password@postgres:5432/%s?sslmode=disable
  
  postgres:
    image: postgres:14
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=%s
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
`, pt.config.Name, pt.config.Name, pt.config.Name)
	
	return os.WriteFile(filepath.Join(outputDir, "docker-compose.yml"), []byte(content), 0644)
}