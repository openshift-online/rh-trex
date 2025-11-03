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

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
	"gopkg.in/resty.v1"

	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/test"
)

func TestMetadataGet(t *testing.T) {
	h, _ := test.RegisterIntegration(t)

	// Build the metadata URL (metadata endpoint is at /api/rh-trex, not /api/rh-trex/v1)
	protocol := "http"
	if h.AppConfig.Server.EnableHTTPS {
		protocol = "https"
	}
	metadataURL := fmt.Sprintf("%s://%s/api/rh-trex", protocol, h.AppConfig.Server.BindAddress)

	// Test GET /api/rh-trex - metadata endpoint doesn't require authentication
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		Get(metadataURL)

	Expect(err).NotTo(HaveOccurred(), "Error getting metadata: %v", err)
	Expect(resp.StatusCode()).To(Equal(http.StatusOK))

	// Parse the response body
	var metadata api.Metadata
	err = json.Unmarshal(resp.Body(), &metadata)
	Expect(err).NotTo(HaveOccurred(), "Error parsing metadata response: %v", err)
	//
	// Verify content type header
	contentType := resp.Header().Get("Content-Type")
	Expect(contentType).To(Equal("application/json"), "Expected Content-Type to be application/json")

	// Verify all metadata fields
	Expect(metadata.ID).To(Equal("trex"), "Expected ID to be 'trex'")
	Expect(metadata.Kind).To(Equal("API"), "Expected Kind to be 'API'")
	Expect(metadata.HREF).To(Equal("/api/rh-trex"), "Expected HREF to match the request path")
	Expect(metadata.Version).NotTo(BeEmpty(), "Expected Version to be set")
	Expect(metadata.BuildTime).NotTo(BeEmpty(), "Expected BuildTime to be set")
}
