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
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/golang/glog"
	"github.com/openshift-online/rh-trex/pkg/api"
)

type metadataHandler struct{}

func NewMetadataHandler() *metadataHandler {
	return &metadataHandler{}
}

// Get sends API documentation response.
func (h metadataHandler) Get(w http.ResponseWriter, r *http.Request) {
	// Set the content type:
	w.Header().Set("Content-Type", "application/json")

	// Prepare the body:
	versions := []api.VersionMetadata{
		{
			ID:   "v1",
			HREF: r.URL.Path + "/v1",
		},
	}
	body := api.Metadata{
		ID:        "trex",
		Kind:      "API",
		HREF:      r.URL.Path,
		Versions:  versions,
		Version:   api.Version,
		BuildTime: api.BuildTime,
	}
	data, err := json.Marshal(body)
	if err != nil {
		api.SendPanic(w, r)
		return
	}

	// Send the response:
	_, err = w.Write(data)
	if err != nil {
		err = fmt.Errorf("can't send response body for request '%s'", r.URL.Path)
		glog.Error(err)
		sentry.CaptureException(err)
		return
	}
}

// GetV1 sends API version v1 documentation response.
func (h metadataHandler) GetV1(w http.ResponseWriter, r *http.Request) {
	// Set the content type:
	w.Header().Set("Content-Type", "application/json")

	// Prepare the body:
	id := "v1"
	collections := []api.CollectionMetadata{
		{
			ID:   "dinosaurs",
			Kind: "DinosaurList",
			HREF: r.URL.Path + "/dinosaurs",
		},
	}
	body := api.VersionMetadata{
		ID:          id,
		Kind:        "APIVersion",
		HREF:        r.URL.Path,
		Collections: collections,
	}
	data, err := json.Marshal(body)
	if err != nil {
		api.SendPanic(w, r)
		return
	}

	// Send the response:
	_, err = w.Write(data)
	if err != nil {
		glog.Errorf("can't send response body for request '%s'", r.URL.Path)
		return
	}
}
