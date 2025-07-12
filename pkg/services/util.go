package services

import (
	"fmt"
	"strings"

	"github.com/openshift-online/rh-trex/pkg/errors"
	coreerrors "github.com/openshift-online/rh-trex-core/errors"
)

// Field names suspected to contain personally identifiable information
var piiFields []string = []string{
	"username",
	"first_name",
	"last_name",
	"email",
	"address",
}

// convertCoreError converts a core library ServiceError to TRex's local ServiceError type
func convertCoreError(coreErr *coreerrors.ServiceError) *errors.ServiceError {
	if coreErr == nil {
		return nil
	}
	
	// Map core error codes to local error codes based on the error reason
	reason := coreErr.Reason
	if strings.Contains(reason, "not found") {
		return errors.NotFound(reason)
	}
	if strings.Contains(reason, "already exists") || strings.Contains(reason, "conflict") {
		return errors.Conflict(reason)
	}
	return errors.GeneralError(reason)
}

// Enhanced error handlers that delegate to core library but return local error types
func handleGetError(resourceType, field string, value interface{}, err error) *errors.ServiceError {
	valueStr := fmt.Sprintf("%v", value)
	coreErr := coreerrors.HandleGetError(resourceType, field, valueStr, err)
	return convertCoreError(coreErr)
}

func handleCreateError(resourceType string, err error) *errors.ServiceError {
	coreErr := coreerrors.HandleCreateError(resourceType, err)
	return convertCoreError(coreErr)
}

func handleUpdateError(resourceType string, err error) *errors.ServiceError {
	coreErr := coreerrors.HandleUpdateError(resourceType, err)
	return convertCoreError(coreErr)
}

func handleDeleteError(resourceType string, err error) *errors.ServiceError {
	coreErr := coreerrors.HandleDeleteError(resourceType, err)
	return convertCoreError(coreErr)
}
