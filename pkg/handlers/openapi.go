package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/openshift-online/rh-trex/data/generated/openapi"
	"github.com/openshift-online/rh-trex/pkg/errors"
)

type openAPIHandler struct {
	OpenAPIDefinitions []byte
}

func NewOpenAPIHandler() (*openAPIHandler, error) {
	// Try to load the fully resolved OpenAPI spec first (no external $refs)
	// This file is generated and may not always be present, so we fallback to embedded asset
	resolvedPath := filepath.Join("pkg", "api", "openapi", "api", "openapi.yaml")
	resolvedData, resolvedErr := os.ReadFile(resolvedPath)
	if resolvedErr == nil {
		// Successfully loaded resolved spec, convert to JSON
		data, err := yaml.YAMLToJSON(resolvedData)
		if err != nil {
			return nil, errors.GeneralError(
				"can't convert resolved OpenAPI specification from YAML to JSON: %v",
				err,
			)
		}
		glog.Info("Loaded fully resolved OpenAPI specification from pkg/api/openapi/api/openapi.yaml")
		return &openAPIHandler{OpenAPIDefinitions: data}, nil
	}

	// Fallback to embedded asset (may have external $refs)
	data, err := openapi.Asset("openapi.yaml")
	if err != nil {
		return nil, errors.GeneralError(
			"can't load OpenAPI specification from asset 'openapi.yaml'",
		)
	}
	data, err = yaml.YAMLToJSON(data)
	if err != nil {
		return nil, errors.GeneralError(
			"can't convert OpenAPI specification loaded from asset 'openapi.yaml' from YAML to JSON",
		)
	}
	glog.Info("Loaded OpenAPI specification from embedded asset (may contain external $refs)")
	return &openAPIHandler{OpenAPIDefinitions: data}, nil
}

func (h *openAPIHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.OpenAPIDefinitions)
}
