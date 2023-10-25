package presenters

import (
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/api"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/api/openapi"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/errors"
)

func ObjectKind(i interface{}) *string {
	result := ""
	switch i.(type) {
	case api.Dinosaur, *api.Dinosaur:
		result = "Dinosaur"
	case errors.ServiceError, *errors.ServiceError:
		result = "Error"
	}

	return openapi.PtrString(result)
}
