package config

const (
	// APIBasePath is the base path for all API endpoints
	APIBasePath = "/api/rh-trex"
	
	// APIV1Path is the v1 API path
	APIV1Path = APIBasePath + "/v1"
	
	// APIErrorsPath is the path for error endpoints
	APIErrorsPath = APIV1Path + "/errors"
	
	// ServiceName is the service identifier used in error codes and API responses
	ServiceName = "rh-trex"
)
