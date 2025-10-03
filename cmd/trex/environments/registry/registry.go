package registry

import (
	"sync"
)

// ServiceLocatorFunc is a function that creates a service locator
type ServiceLocatorFunc func(env interface{}) interface{}

// ServiceRegistry holds registered services
type ServiceRegistry struct {
	mu       sync.RWMutex
	services map[string]ServiceLocatorFunc
}

var globalRegistry = &ServiceRegistry{
	services: make(map[string]ServiceLocatorFunc),
}

// RegisterService registers a service with the global registry
func RegisterService(name string, locatorFunc ServiceLocatorFunc) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.services[name] = locatorFunc
}

// ServicesInterface defines the interface for the Services struct
type ServicesInterface interface {
	SetService(name string, service interface{})
}

// LoadDiscoveredServices loads all registered services into the Services struct
func LoadDiscoveredServices(services ServicesInterface, env interface{}) {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	for name, locatorFunc := range globalRegistry.services {
		// Call the locator function to create the service and store it in the registry
		serviceLocator := locatorFunc(env)
		services.SetService(name, serviceLocator)
	}
}