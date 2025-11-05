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

	"github.com/openshift-online/rh-trex/test"
)

func TestOpenAPIGet(t *testing.T) {
	h, _ := test.RegisterIntegration(t)

	protocol := "http"
	if h.AppConfig.Server.EnableHTTPS {
		protocol = "https"
	}
	openAPIURL := fmt.Sprintf("%s://%s/api/rh-trex/v1/openapi", protocol, h.AppConfig.Server.BindAddress)

	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		Get(openAPIURL)

	Expect(err).NotTo(HaveOccurred(), "Error getting OpenAPI spec: %v", err)
	Expect(resp.StatusCode()).To(Equal(http.StatusOK), "Expected status code 200")

	// Verify content type header
	contentType := resp.Header().Get("Content-Type")
	Expect(contentType).To(Equal("application/json"), "Expected Content-Type to be application/json")

	// Verify the response body is valid JSON
	var openAPISpec map[string]interface{}
	err = json.Unmarshal(resp.Body(), &openAPISpec)
	Expect(err).NotTo(HaveOccurred(), "Error parsing OpenAPI JSON response: %v", err)

	// Verify the OpenAPI spec has required fields
	Expect(openAPISpec).To(HaveKey("openapi"), "Expected OpenAPI spec to have 'openapi' field")
	Expect(openAPISpec).To(HaveKey("info"), "Expected OpenAPI spec to have 'info' field")
	Expect(openAPISpec).To(HaveKey("paths"), "Expected OpenAPI spec to have 'paths' field")

	// Verify the OpenAPI version
	Expect(openAPISpec["openapi"]).To(Equal("3.0.0"), "Expected OpenAPI version to be 3.0.0")

	// Verify info section
	info, ok := openAPISpec["info"].(map[string]interface{})
	Expect(ok).To(BeTrue(), "Expected 'info' to be an object")
	Expect(info).To(HaveKey("title"), "Expected 'info' to have 'title' field")
	Expect(info).To(HaveKey("version"), "Expected 'info' to have 'version' field")
}

func TestOpenAPIUIGet(t *testing.T) {
	h, _ := test.RegisterIntegration(t)

	protocol := "http"
	if h.AppConfig.Server.EnableHTTPS {
		protocol = "https"
	}
	openAPIUIURL := fmt.Sprintf("%s://%s/api/rh-trex/v1/openapi.html", protocol, h.AppConfig.Server.BindAddress)

	t.Logf("OpenAPI UI URL: %s", openAPIUIURL)

	resp, err := resty.R().
		SetHeader("Content-Type", "text/html").
		Get(openAPIUIURL)

	Expect(err).NotTo(HaveOccurred(), "Error getting OpenAPI UI: %v", err)
	Expect(resp.StatusCode()).To(Equal(http.StatusOK), "Expected status code 200")

	// Verify content type header
	contentType := resp.Header().Get("Content-Type")
	Expect(contentType).To(Equal("text/html"), "Expected Content-Type to be text/html")

	// Verify the response body is not empty
	body := resp.String()
	Expect(body).NotTo(BeEmpty(), "Expected OpenAPI UI HTML to not be empty")

	// Verify the HTML contains expected elements (basic checks)
	Expect(body).To(ContainSubstring("<html"), "Expected HTML to contain '<html' tag")
	Expect(body).To(ContainSubstring("</html>"), "Expected HTML to contain '</html>' tag")
}
