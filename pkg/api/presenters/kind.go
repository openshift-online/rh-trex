package presenters

import (
	"fmt"

	"github.com/openshift-online/rh-trex/pkg/api/openapi"
	"github.com/openshift-online/rh-trex/pkg/errors"
)

type KindMappingFunc func(interface{}) string

var kindRegistry = make(map[string]KindMappingFunc)

func RegisterKind(objType interface{}, kindValue string) {
	typeName := fmt.Sprintf("%T", objType)
	kindRegistry[typeName] = func(interface{}) string {
		return kindValue
	}
}

func LoadDiscoveredKinds(i interface{}) string {
	typeName := fmt.Sprintf("%T", i)
	if mappingFunc, found := kindRegistry[typeName]; found {
		return mappingFunc(i)
	}
	return ""
}

func ObjectKind(i interface{}) *string {
	result := ""

	// Check auto-discovered kinds first
	if discoveredKind := LoadDiscoveredKinds(i); discoveredKind != "" {
		result = discoveredKind
	} else {
		// Built-in mappings
		switch i.(type) {
		case errors.ServiceError, *errors.ServiceError:
			result = "Error"
		}
	}

	return openapi.PtrString(result)
}
