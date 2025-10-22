package events

import (
	"github.com/openshift-online/rh-trex/cmd/trex/environments"
	"github.com/openshift-online/rh-trex/cmd/trex/environments/registry"
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/services"
)

// ServiceLocator Service Locator
type ServiceLocator func() services.EventService

func NewServiceLocator(env *environments.Env) ServiceLocator {
	return func() services.EventService {
		return services.NewEventService(dao.NewEventDao(&env.Database.SessionFactory))
	}
}

// Service helper function to get the event service from the registry
func Service(s *environments.Services) services.EventService {
	if s == nil {
		return nil
	}
	if obj := s.GetService("Events"); obj != nil {
		locator := obj.(ServiceLocator)
		return locator()
	}
	return nil
}

func init() {
	// Service registration
	registry.RegisterService("Events", func(env interface{}) interface{} {
		return NewServiceLocator(env.(*environments.Env))
	})
}
