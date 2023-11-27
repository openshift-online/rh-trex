package presenters

import (
	"fmt"

	"github.com/openshift-online/rh-trex/pkg/api/openapi"

	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/errors"
)

const (
	BasePath = "/api/rh-trex/v1"
)

func ObjectPath(id string, obj interface{}) *string {
	return openapi.PtrString(fmt.Sprintf("%s/%s/%s", BasePath, path(obj), id))
}

func path(i interface{}) string {
	switch i.(type) {
	case api.Dinosaur, *api.Dinosaur:
		return "dinosaurs"
	// Add manually
	case api.Subscription, *api.Subscription:
		return "subscriptions"
	case errors.ServiceError, *errors.ServiceError:
		return "errors"
	default:
		return ""
	}
}
