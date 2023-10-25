package presenters

import (
	"fmt"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/api/openapi"

	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/api"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/errors"
)

const (
	BasePath = "/api/ocm-example-service/v1"
)

func ObjectPath(id string, obj interface{}) *string {
	return openapi.PtrString(fmt.Sprintf("%s/%s/%s", BasePath, path(obj), id))
}

func path(i interface{}) string {
	switch i.(type) {
	case api.Dinosaur, *api.Dinosaur:
		return "dinosaurs"
	case errors.ServiceError, *errors.ServiceError:
		return "errors"
	default:
		return ""
	}
}
