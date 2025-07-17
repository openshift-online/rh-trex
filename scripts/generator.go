package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/pflag"
)

/*

scripts/generator.go

This script generates basic CRUD functionality for a new Kind.

It's rude and crude, but it generates working code.

TODO: all of it can be better

*/

var (
	kind                        string = "Asteroid"
	repo                        string = "github.com/openshift-online"
	project                     string = "rh-trex"
	openapiEndpointStart        string = "# NEW ENDPOINT START"
	openapiEndpointEnd          string = "# NEW ENDPOINT END"
	openApiSchemaStart          string = "# NEW SCHEMA START"
	openApiSchemaEnd            string = "# NEW SCHEMA END"
	openApiEndpointMatchingLine string = "  # AUTO-ADD NEW PATHS"
	openApiSchemaMatchingLine   string = "    # AUTO-ADD NEW SCHEMAS"
)

func init() {
	_ = flag.Set("logtostderr", "true")
	flags := pflag.CommandLine
	flags.AddGoFlagSet(flag.CommandLine)

	flags.StringVar(&kind, "kind", kind, "the name of the kind.  e.g Account or User")
	flags.StringVar(&repo, "repo", repo, "the name of the repo.  e.g github.com/yourproject")
	flags.StringVar(&project, "project", project, "the name of the project.  e.g rh-trex")
}

func getCmdDir() string {
	entries, err := os.ReadDir("cmd")
	if err != nil {
		panic(err)
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			return entry.Name()
		}
	}
	
	panic("No command directory found in cmd/")
}

func main() {
	// Parse flags
	pflag.Parse()

	templates := []string{
		"api",
		"presenters",
		"dao",
		"services",
		"mock",
		"migration",
		"test",
		"test-factories",
		"handlers",
		"openapi-kind",
		"servicelocator",
	}

	for _, nm := range templates {
		path := fmt.Sprintf("templates/generate-%s.txt", nm)
		contents, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		kindTmpl, err := template.New(nm).Parse(string(contents))
		if err != nil {
			panic(err)
		}

		kindLowerCamel := strings.ToLower(string(kind[0])) + kind[1:]
		kindSnakeCase := toSnakeCase(kind)
		projectCamelCase := toCamelCase(project)
		k := myWriter{
			Project:             project,
			ProjectCamelCase:    projectCamelCase,
			Repo:                repo,
			Cmd:                 getCmdDir(),
			Kind:                kind,
			KindPlural:          fmt.Sprintf("%ss", kind),
			KindLowerPlural:     kindLowerCamel + "s",
			KindLowerSingular:   kindLowerCamel,
			KindSnakeCasePlural: kindSnakeCase + "s",
		}

		now := time.Now()
		k.ID = fmt.Sprintf("%d%s%s%s%s", now.Year(), datePad(int(now.Month())), datePad(now.Day()), datePad(now.Hour()), datePad(now.Minute()))

		outputPaths := map[string]string{
			"generate-api":            fmt.Sprintf("pkg/%s/%s.go", nm, k.KindLowerSingular),
			"generate-presenters":     fmt.Sprintf("pkg/api/presenters/%s.go", k.KindLowerSingular),
			"generate-dao":            fmt.Sprintf("pkg/%s/%s.go", nm, k.KindLowerSingular),
			"generate-handlers":       fmt.Sprintf("pkg/%s/%s.go", nm, k.KindLowerSingular),
			"generate-migration":      fmt.Sprintf("pkg/db/migrations/%s_add_%s.go", k.ID, k.KindLowerPlural),
			"generate-mock":           fmt.Sprintf("pkg/dao/mocks/%s.go", k.KindLowerSingular),
			"generate-openapi-kind":   fmt.Sprintf("openapi/openapi.%s.yaml", k.KindLowerPlural),
			"generate-test-factories": fmt.Sprintf("test/factories/%s.go", k.KindLowerPlural),
			"generate-test":           fmt.Sprintf("test/integration/%s_test.go", k.KindLowerPlural),
			"generate-services":       fmt.Sprintf("pkg/%s/%s.go", nm, k.KindLowerSingular),
			"generate-servicelocator": fmt.Sprintf("cmd/%s/environments/locator_%s.go", k.Cmd, k.KindLowerSingular),
		}

		outputPath, ok := outputPaths["generate-"+nm]
		if !ok {
			panic("expected to find outputPath for " + nm)
		}

		f, err := os.Create(outputPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		w := bufio.NewWriter(f)
		err = kindTmpl.Execute(w, k)
		if err != nil {
			panic(err)
		}
		w.Flush()
		f.Sync()

		if strings.EqualFold("generate-"+nm, "generate-openapi-kind") {
			modifyOpenapi("openapi/openapi.yaml", fmt.Sprintf("openapi/openapi.%s.yaml", k.KindLowerPlural))
		}
		
		// Add controller registration and presenter mappings after all templates are processed
		if nm == "services" {
			addControllerRegistration(k)
			addPresenterMappings(k)
			addServiceLocatorToTypes(k)
			addServiceLocatorToFramework(k)
			addRouteRegistration(k)
		}
		
		// Add migration registration after migration is generated
		if nm == "migration" {
			addMigrationRegistration(k)
		}
	}
}

func datePad(d int) string {
	if d < 10 {
		return fmt.Sprintf("0%d", d)
	}
	return fmt.Sprintf("%d", d)
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "-")
	var result strings.Builder
	
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(string(part[0])) + part[1:])
		}
	}
	
	return result.String()
}

type myWriter struct {
	Repo                     string
	Project                  string
	ProjectCamelCase         string
	Cmd                      string
	Kind                     string
	KindPlural               string
	KindLowerPlural          string
	KindLowerSingular        string
	KindSnakeCasePlural      string
	ID                       string
}

func modifyOpenapi(mainPath string, kindPath string) {
	endpointStrings := readBetweenLines(kindPath, openapiEndpointStart, openapiEndpointEnd)
	kindFileName := strings.Split(kindPath, "/")[1]
	for _, line := range endpointStrings {
		endpointStr := strings.TrimSpace(line)
		endpointStr = strings.Replace(endpointStr, "/", "~1", -1)
		endpointStr = strings.Replace(endpointStr, ":", "", -1)
		refPath := fmt.Sprintf(`    $ref: '%s#/paths/%s'`, kindFileName, endpointStr)
		pathsLine := fmt.Sprintf("%s%s", line, refPath)
		writeAfterLine(mainPath, openApiEndpointMatchingLine, pathsLine)
	}
	schemaStrings := readBetweenLines(kindPath, openApiSchemaStart, openApiSchemaEnd)
	for _, line := range schemaStrings {
		schemaStr := strings.TrimSpace(line)
		schemaStr = strings.Replace(schemaStr, ":", "", -1)
		refPath := fmt.Sprintf(`      $ref: '%s#/components/schemas/%s'`, kindFileName, schemaStr)
		pathsLine := fmt.Sprintf("%s%s", line, refPath)
		writeAfterLine(mainPath, openApiSchemaMatchingLine, pathsLine)
	}
}

func readBetweenLines(path string, startLine string, endLine string) []string {
	readFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	readFlag := false
	var totalMatches []string
	var matchedString strings.Builder
	for fileScanner.Scan() {
		trimmed := strings.TrimSpace(fileScanner.Text())
		if trimmed == startLine {
			readFlag = true
		} else if trimmed == endLine {
			readFlag = false
			totalMatches = append(totalMatches, matchedString.String())
			matchedString.Reset()
		} else if readFlag {
			matchedString.WriteString(fileScanner.Text() + "\n")
		}
	}
	readFile.Close()
	return totalMatches
}

func writeAfterLine(path string, matchingLine string, lineToWrite string) {
	input, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	_ = strings.Replace(string(input), matchingLine, lineToWrite+"\n"+matchingLine, -1)
	output := bytes.Replace(input, []byte(matchingLine), []byte(lineToWrite+"\n"+matchingLine), -1)
	if err = os.WriteFile(path, output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addControllerRegistration(k myWriter) {
	controllerFile := fmt.Sprintf("cmd/%s/server/controllers.go", k.Cmd)
	
	// Add controller registration
	controllerRegistration := fmt.Sprintf(`
	%sServices := env().Services.%s()

	s.KindControllerManager.Add(&controllers.ControllerConfig{
		Source: "%s",
		Handlers: map[api.EventType][]controllers.ControllerHandlerFunc{
			api.CreateEventType: {%sServices.OnUpsert},
			api.UpdateEventType: {%sServices.OnUpsert},
			api.DeleteEventType: {%sServices.OnDelete},
		},
	})
`, k.KindLowerSingular, k.KindPlural, k.KindPlural, k.KindLowerSingular, k.KindLowerSingular, k.KindLowerSingular)
	
	// Insert before the return statement
	matchingLine := `	return s`
	writeBeforePattern(controllerFile, matchingLine, controllerRegistration)
}

func addPresenterMappings(k myWriter) {
	// Add Kind mapping to presenters/kind.go
	kindFile := "pkg/api/presenters/kind.go"
	kindMapping := fmt.Sprintf(`	case api.%s, *api.%s:
		result = "%s"`, k.Kind, k.Kind, k.Kind)
	
	// Insert before the errors case
	kindMatchingLine := `	case errors.ServiceError, *errors.ServiceError:
		result = "Error"`
	writeBeforePattern(kindFile, kindMatchingLine, kindMapping)
	
	// Add path mapping to presenters/path.go  
	pathFile := "pkg/api/presenters/path.go"
	pathMapping := fmt.Sprintf(`	case api.%s, *api.%s:
		return "%s"`, k.Kind, k.Kind, k.KindSnakeCasePlural)
	
	// Insert before the errors case
	pathMatchingLine := `	case errors.ServiceError, *errors.ServiceError:
		return "errors"`
	writeBeforePattern(pathFile, pathMatchingLine, pathMapping)
}

func writeBeforePattern(path string, matchingLine string, lineToWrite string) {
	input, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	output := bytes.Replace(input, []byte(matchingLine), []byte(lineToWrite+"\n"+matchingLine), -1)
	if err = os.WriteFile(path, output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addServiceLocatorToTypes(k myWriter) {
	typesFile := fmt.Sprintf("cmd/%s/environments/types.go", k.Cmd)
	
	// Add service locator field to Services struct
	serviceField := fmt.Sprintf("\t%s %sServiceLocator", k.KindPlural, k.Kind)
	
	// Insert after the "// ADD LOCATORS HERE" comment
	matchingLine := "\t// ADD LOCATORS HERE"
	writeAfterLine(typesFile, matchingLine, serviceField)
}

func addServiceLocatorToFramework(k myWriter) {
	frameworkFile := fmt.Sprintf("cmd/%s/environments/framework.go", k.Cmd)
	
	// Add service locator initialization to LoadServices method
	serviceInitialization := fmt.Sprintf("\te.Services.%s = New%sServiceLocator(e)", k.KindPlural, k.Kind)
	
	// Insert after the "// ADD SERVICES HERE" comment
	matchingLine := "\t// ADD SERVICES HERE"
	writeAfterLine(frameworkFile, matchingLine, serviceInitialization)
}

func addRouteRegistration(k myWriter) {
	routesFile := fmt.Sprintf("cmd/%s/server/routes.go", k.Cmd)
	
	// Add handler creation and route registration
	routeRegistration := fmt.Sprintf(`
	%sHandler := handlers.New%sHandler(services.%s(), services.Generic())

	//  /api/rh-trex/v1/%s
	apiV1%sRouter := apiV1Router.PathPrefix("/%s").Subrouter()
	apiV1%sRouter.HandleFunc("", %sHandler.List).Methods(http.MethodGet)
	apiV1%sRouter.HandleFunc("/{id}", %sHandler.Get).Methods(http.MethodGet)
	apiV1%sRouter.HandleFunc("", %sHandler.Create).Methods(http.MethodPost)
	apiV1%sRouter.HandleFunc("/{id}", %sHandler.Patch).Methods(http.MethodPatch)
	apiV1%sRouter.HandleFunc("/{id}", %sHandler.Delete).Methods(http.MethodDelete)
	apiV1%sRouter.Use(authMiddleware.AuthenticateAccountJWT)
	apiV1%sRouter.Use(authzMiddleware.AuthorizeApi)`, 
		k.KindLowerSingular, k.Kind, k.KindPlural,
		k.KindLowerPlural, 
		k.KindPlural, k.KindLowerPlural,
		k.KindPlural, k.KindLowerSingular,
		k.KindPlural, k.KindLowerSingular,
		k.KindPlural, k.KindLowerSingular,
		k.KindPlural, k.KindLowerSingular,
		k.KindPlural, k.KindLowerSingular,
		k.KindPlural, k.KindLowerSingular,
		k.KindPlural)
	
	// Insert after the "// ADD ROUTES HERE" comment
	matchingLine := "\t// ADD ROUTES HERE"
	writeAfterLine(routesFile, matchingLine, routeRegistration)
}

func addMigrationRegistration(k myWriter) {
	migrationFile := "pkg/db/migrations/migration_structs.go"
	
	// Add migration function call
	migrationCall := fmt.Sprintf("\tadd%s(),", k.KindPlural)
	
	// Insert after the "// ADD MIGRATIONS HERE" comment
	matchingLine := "\t// ADD MIGRATIONS HERE"
	writeAfterLine(migrationFile, matchingLine, migrationCall)
}
