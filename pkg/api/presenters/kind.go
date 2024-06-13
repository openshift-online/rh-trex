package presenters

import (
	"github.com/openshift-online/rh-trex/pkg/api/openapi"
	"github.com/openshift-online/rh-trex/pkg/util"
)

func ObjectKind(i interface{}) *string {
	result := util.GetBaseType(i)

	return openapi.PtrString(result)
}
