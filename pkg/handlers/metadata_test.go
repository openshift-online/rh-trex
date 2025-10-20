/*
Copyright (c) 2018 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openshift-online/rh-trex/pkg/api"
)

func TestMetadataHandler_Get_IncludesVersion(t *testing.T) {
	// Set a test version
	api.Version = "test-version-123"

	// Create a test request
	req := httptest.NewRequest("GET", "/api/rh-trex", nil)
	w := httptest.NewRecorder()

	// Create handler and call it
	handler := NewMetadataHandler()
	handler.Get(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Parse the response body
	var metadata api.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify the version is included
	if metadata.Version != "test-version-123" {
		t.Errorf("Expected version 'test-version-123', got '%s'", metadata.Version)
	}

	// Verify other fields are still present
	if metadata.ID != "trex" {
		t.Errorf("Expected ID 'trex', got '%s'", metadata.ID)
	}

	if metadata.Kind != "API" {
		t.Errorf("Expected Kind 'API', got '%s'", metadata.Kind)
	}

	if len(metadata.Versions) != 1 {
		t.Errorf("Expected 1 version, got %d", len(metadata.Versions))
	}
}

func TestMetadataHandler_GetV1(t *testing.T) {
	// Create a test request
	req := httptest.NewRequest("GET", "/api/rh-trex/v1", nil)
	w := httptest.NewRecorder()

	// Create handler and call it
	handler := NewMetadataHandler()
	handler.GetV1(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Parse the response body
	var versionMetadata api.VersionMetadata
	if err := json.NewDecoder(resp.Body).Decode(&versionMetadata); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify fields
	if versionMetadata.ID != "v1" {
		t.Errorf("Expected ID 'v1', got '%s'", versionMetadata.ID)
	}

	if versionMetadata.Kind != "APIVersion" {
		t.Errorf("Expected Kind 'APIVersion', got '%s'", versionMetadata.Kind)
	}

	if len(versionMetadata.Collections) < 1 {
		t.Errorf("Expected at least 1 collection, got %d", len(versionMetadata.Collections))
	}
}
