package presenters

import (
	"fmt"

	"github.com/openshift-online/rh-trex/pkg/api/openapi"
	"github.com/openshift-online/rh-trex/pkg/util"
)

const (
	BasePath = "/api/rh-trex/v1"
)

func ObjectPath(id string, obj interface{}) *string {
	return openapi.PtrString(fmt.Sprintf("%s/%s/%s", BasePath, path(obj), id))
}

func path(i interface{}) string {
	return fmt.Sprintf("%ss", util.ToSnakeCase(util.GetBaseType(i)))
}
