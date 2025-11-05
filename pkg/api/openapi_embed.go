package api

import (
	"embed"
	"io/fs"
)

//go:embed openapi/api/openapi.yaml
var openapiFS embed.FS

// GetOpenAPISpec returns the embedded OpenAPI YAML file contents
func GetOpenAPISpec() ([]byte, error) {
	return fs.ReadFile(openapiFS, "openapi/api/openapi.yaml")
}

