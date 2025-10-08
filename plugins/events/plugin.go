package events

import (
	"github.com/openshift-online/rh-trex/cmd/trex/environments"
	"github.com/openshift-online/rh-trex/cmd/trex/environments/registry"
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/services"
)

// Service Locator
type EventServiceLocator func() services.EventService

func NewEventServiceLocator(env *environments.Env) EventServiceLocator {
	return func() services.EventService {
		return services.NewEventService(dao.NewEventDao(&env.Database.SessionFactory))
	}
}

// EventService helper function to get the event service from the registry
func EventService(s *environments.Services) services.EventService {
	if s == nil {
		return nil
	}
	if obj := s.GetService("Events"); obj != nil {
		locator := obj.(EventServiceLocator)
		return locator()
	}
	return nil
}

func init() {
	// Service registration
	registry.RegisterService("Events", func(env interface{}) interface{} {
		return NewEventServiceLocator(env.(*environments.Env))
	})
}
