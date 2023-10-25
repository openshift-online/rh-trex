package presenters

import (
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/api/openapi"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/errors"
)

func PresentError(err *errors.ServiceError) openapi.Error {
	return err.AsOpenapiError("")
}
