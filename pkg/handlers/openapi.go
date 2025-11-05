package handlers

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/errors"
)

//go:embed openapi-ui.html
var openapiui embed.FS

type openAPIHandler struct {
	openAPIDefinitions []byte
	uiContent          []byte
}

func NewOpenAPIHandler() (*openAPIHandler, error) {
	// Load the fully resolved OpenAPI spec from embedded filesystem
	resolvedData, err := api.GetOpenAPISpec()
	if err != nil {
		return nil, errors.GeneralError(
			"can't load OpenAPI specification from embedded file: %v",
			err,
		)
	}

	// Convert YAML to JSON
	data, err := yaml.YAMLToJSON(resolvedData)
	if err != nil {
		return nil, errors.GeneralError(
			"can't convert OpenAPI specification from YAML to JSON: %v",
			err,
		)
	}
	glog.Info("Loaded fully resolved OpenAPI specification from embedded pkg/api/openapi/api/openapi.yaml")

	// Load the OpenAPI UI HTML content
	uiContent, err := fs.ReadFile(openapiui, "openapi-ui.html")
	if err != nil {
		return nil, errors.GeneralError(
			"can't load OpenAPI UI HTML from embedded file: %v",
			err,
		)
	}
	glog.Info("Loaded OpenAPI UI HTML from embedded file")

	return &openAPIHandler{
		openAPIDefinitions: data,
		uiContent:          uiContent,
	}, nil
}

func (h *openAPIHandler) GetOpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.openAPIDefinitions)
}

func (h *openAPIHandler) GetOpenAPIUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.uiContent)
}
