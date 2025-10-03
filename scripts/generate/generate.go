package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"
)

type GeneratorConfig struct {
	Kind    string
	Repo    string
	Project string
}

func generateEntity(config *GeneratorConfig) error {
	fmt.Printf("🚀 Generating %s entity...\n", config.Kind)
	
	// Auto-detect project if not provided
	if config.Project == "" {
		project, err := detectProjectName()
		if err != nil {
			return fmt.Errorf("failed to detect project name: %v", err)
		}
		config.Project = project
		config.Repo = "github.com/openshift-online"
	}
	
	// Generate entity files
	if err := generateEntityFiles(config); err != nil {
		return fmt.Errorf("failed to generate entity files: %v", err)
	}
	
	// Update service locator
	if err := updateServiceLocator(config); err != nil {
		return fmt.Errorf("failed to update service locator: %v", err)
	}
	
	// Update OpenAPI spec
	if err := updateOpenAPISpec(config); err != nil {
		return fmt.Errorf("failed to update OpenAPI spec: %v", err)
	}
	
	fmt.Printf("✅ Successfully generated %s entity!\n", config.Kind)
	fmt.Printf("📋 Next step: Run 'make generate'\n")
	
	return nil
}

func detectProjectName() (string, error) {
	// Look for cmd directory
	cmdEntries, err := os.ReadDir("cmd")
	if err != nil {
		return "", err
	}
	
	for _, entry := range cmdEntries {
		if entry.IsDir() && entry.Name() != "." && entry.Name() != ".." {
			return entry.Name(), nil
		}
	}
	
	return "", fmt.Errorf("no project directory found in cmd/")
}

func generateEntityFiles(config *GeneratorConfig) error {
	lowerKind := strings.ToLower(config.Kind)
	
	// Create directories
	dirs := []string{
		fmt.Sprintf("plugins/%s", lowerKind),
		fmt.Sprintf("pkg/api/openapi/%s", lowerKind),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	
	// Generate files from templates
	files := map[string]string{
		fmt.Sprintf("plugins/%s/handler.go", lowerKind):                    entityHandlerTemplate,
		fmt.Sprintf("plugins/%s/service.go", lowerKind):                    entityServiceTemplate,
		fmt.Sprintf("plugins/%s/dao.go", lowerKind):                        entityDAOTemplate,
		fmt.Sprintf("plugins/%s/presenter.go", lowerKind):                  entityPresenterTemplate,
		fmt.Sprintf("pkg/api/openapi/%s/model_%s.go", lowerKind, lowerKind): entityModelTemplate,
		fmt.Sprintf("openapi/openapi.%s.yaml", lowerKind):                  openAPISpecTemplate,
	}
	
	for filename, tmplContent := range files {
		if err := generateFileFromTemplate(filename, tmplContent, config); err != nil {
			return fmt.Errorf("failed to generate %s: %v", filename, err)
		}
	}
	
	return nil
}

func generateFileFromTemplate(filename, tmplContent string, config *GeneratorConfig) error {
	tmpl, err := template.New("entity").Parse(tmplContent)
	if err != nil {
		return err
	}
	
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	data := struct {
		Kind      string
		LowerKind string
		Repo      string
		Project   string
		Date      string
	}{
		Kind:      config.Kind,
		LowerKind: strings.ToLower(config.Kind),
		Repo:      config.Repo,
		Project:   config.Project,
		Date:      time.Now().Format("2006-01-02"),
	}
	
	return tmpl.Execute(file, data)
}

func updateServiceLocator(config *GeneratorConfig) error {
	// This is a simplified version - the full implementation would need
	// to handle the service locator registration properly
	fmt.Printf("⚠️  Manual step required: Add %sServiceLocator to types.go and framework.go\n", config.Kind)
	return nil
}

func updateOpenAPISpec(config *GeneratorConfig) error {
	// Read main openapi.yaml
	content, err := os.ReadFile("openapi/openapi.yaml")
	if err != nil {
		return err
	}
	
	lowerKind := strings.ToLower(config.Kind)
	
	// Add path reference
	pathRef := fmt.Sprintf("  /api/%s/v1/%ss:\n    $ref: 'openapi.%s.yaml#/paths/~1api~1%s~1v1~1%ss'\n  /api/%s/v1/%ss/{id}:\n    $ref: 'openapi.%s.yaml#/paths/~1api~1%s~1v1~1%ss~1{id}'\n", 
		config.Project, lowerKind, lowerKind, config.Project, lowerKind, config.Project, lowerKind, lowerKind, config.Project, lowerKind)
	
	// Add schema reference
	schemaRef := fmt.Sprintf("    %s:\n      $ref: 'openapi.%s.yaml#/components/schemas/%s'\n    %sList:\n      $ref: 'openapi.%s.yaml#/components/schemas/%sList'\n    %sPatchRequest:\n      $ref: 'openapi.%s.yaml#/components/schemas/%sPatchRequest'\n",
		config.Kind, lowerKind, config.Kind, config.Kind, lowerKind, config.Kind, config.Kind, lowerKind, config.Kind)
	
	// Simple string replacement for now
	updatedContent := strings.Replace(string(content), "  # AUTO-ADD NEW PATHS", pathRef+"  # AUTO-ADD NEW PATHS", 1)
	updatedContent = strings.Replace(updatedContent, "    # AUTO-ADD NEW SCHEMAS", schemaRef+"    # AUTO-ADD NEW SCHEMAS", 1)
	
	return os.WriteFile("openapi/openapi.yaml", []byte(updatedContent), 0644)
}

// Templates for entity generation
const entityHandlerTemplate = `package {{.LowerKind}}

import (
	"net/http"
	"{{.Repo}}/{{.Project}}/pkg/api"
	"{{.Repo}}/{{.Project}}/pkg/errors"
	"github.com/gorilla/mux"
)

type {{.Kind}}Handler struct {
	service {{.Kind}}Service
}

func New{{.Kind}}Handler(service {{.Kind}}Service) *{{.Kind}}Handler {
	return &{{.Kind}}Handler{service: service}
}

func (h *{{.Kind}}Handler) List(w http.ResponseWriter, r *http.Request) {
	cfg := &api.ListArguments{}
	{{.LowerKind}}s, paging, err := h.service.List(cfg)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}
	
	{{.LowerKind}}List := api.{{.Kind}}List{
		Kind:  "{{.Kind}}List",
		Page:  paging.Page,
		Size:  paging.Size,
		Total: paging.Total,
		Items: {{.LowerKind}}s,
	}
	
	api.WriteJSONResponse(w, http.StatusOK, {{.LowerKind}}List)
}`

const entityServiceTemplate = `package {{.LowerKind}}

import (
	"{{.Repo}}/{{.Project}}/pkg/api"
	"{{.Repo}}/{{.Project}}/pkg/errors"
)

type {{.Kind}}Service interface {
	List(cfg *api.ListArguments) ([]*api.{{.Kind}}, *api.PaginationConfig, error)
	Get(id string) (*api.{{.Kind}}, error)
	Create({{.LowerKind}} *api.{{.Kind}}) (*api.{{.Kind}}, error)
	Update(id string, patch *api.{{.Kind}}PatchRequest) (*api.{{.Kind}}, error)
	Delete(id string) error
}

type {{.LowerKind}}Service struct {
	dao {{.Kind}}Dao
}

func New{{.Kind}}Service(dao {{.Kind}}Dao) {{.Kind}}Service {
	return &{{.LowerKind}}Service{dao: dao}
}`

const entityDAOTemplate = `package {{.LowerKind}}

import (
	"{{.Repo}}/{{.Project}}/pkg/api"
	"{{.Repo}}/{{.Project}}/pkg/db"
)

type {{.Kind}}Dao interface {
	All() ([]*api.{{.Kind}}, error)
	FindById(id string) (*api.{{.Kind}}, error)
	Create({{.LowerKind}} *api.{{.Kind}}) error
	Update({{.LowerKind}} *api.{{.Kind}}) error
	Delete({{.LowerKind}} *api.{{.Kind}}) error
}

type sql{{.Kind}}Dao struct {
	sessionFactory *db.SessionFactory
}

func New{{.Kind}}Dao(sessionFactory *db.SessionFactory) {{.Kind}}Dao {
	return &sql{{.Kind}}Dao{sessionFactory: sessionFactory}
}`

const entityPresenterTemplate = `package {{.LowerKind}}

import "{{.Repo}}/{{.Project}}/pkg/api"

func Present{{.Kind}}({{.LowerKind}} *api.{{.Kind}}) *api.{{.Kind}} {
	return &api.{{.Kind}}{
		ID:        {{.LowerKind}}.ID,
		Kind:      "{{.Kind}}",
		Href:      "/api/{{.Project}}/v1/{{.LowerKind}}s/" + {{.LowerKind}}.ID,
		CreatedAt: {{.LowerKind}}.CreatedAt,
		UpdatedAt: {{.LowerKind}}.UpdatedAt,
		Name:      {{.LowerKind}}.Name,
	}
}`

const entityModelTemplate = `package openapi

import "time"

type {{.Kind}} struct {
	ID        string    ` + "`json:\"id\"`" + `
	Kind      string    ` + "`json:\"kind\"`" + `
	Href      string    ` + "`json:\"href\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
	Name      string    ` + "`json:\"name\"`" + `
}`

const openAPISpecTemplate = `openapi: 3.0.0
info:
  title: {{.Kind}} API
  version: 1.0.0
paths:
  /api/{{.Project}}/v1/{{.LowerKind}}s:
    get:
      summary: List {{.LowerKind}}s
      operationId: list{{.Kind}}s
components:
  schemas:
    {{.Kind}}:
      type: object
      properties:
        id:
          type: string
        name:
          type: string`