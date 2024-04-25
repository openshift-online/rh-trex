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

		k := myWriter{
			Kind:              kind,
			KindPlural:        fmt.Sprintf("%ss", kind),
			KindLowerPlural:   strings.ToLower(fmt.Sprintf("%ss", kind)),
			KindLowerSingular: strings.ToLower(kind),
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
			"generate-servicelocator": fmt.Sprintf("cmd/ensemble/environments/locator_%s.go", k.KindLowerSingular),
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
	}
}

func datePad(d int) string {
	if d < 10 {
		return fmt.Sprintf("0%d", d)
	}
	return fmt.Sprintf("%d", d)
}

type myWriter struct {
	Kind              string
	KindPlural        string
	KindLowerPlural   string
	KindLowerSingular string
	ID                string
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
