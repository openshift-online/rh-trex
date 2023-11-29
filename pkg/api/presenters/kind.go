package presenters

import (
	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/api/openapi"
	"github.com/openshift-online/rh-trex/pkg/errors"
)

func ObjectKind(i interface{}) *string {
	result := ""
	switch i.(type) {
	case api.Dinosaur, *api.Dinosaur:
		result = "Dinosaur"
	//add manually
	case api.Subscription, *api.Subscription:
		result = "Subscription"
	case errors.ServiceError, *errors.ServiceError:
		result = "Error"
	}

	return openapi.PtrString(result)
}
