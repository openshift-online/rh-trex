package generic

import (
	"github.com/openshift-online/rh-trex/cmd/trex/environments"
	"github.com/openshift-online/rh-trex/cmd/trex/environments/registry"
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/services"
)

// Service Locator
type GenericServiceLocator func() services.GenericService

func NewGenericServiceLocator(env *environments.Env) GenericServiceLocator {
	return func() services.GenericService {
		return services.NewGenericService(dao.NewGenericDao(&env.Database.SessionFactory))
	}
}

// GenericService helper function to get the generic service from the registry
func GenericService(s *environments.Services) services.GenericService {
	if s == nil {
		return nil
	}
	if obj := s.GetService("Generic"); obj != nil {
		locator := obj.(GenericServiceLocator)
		return locator()
	}
	return nil
}

func init() {
	// Service registration
	registry.RegisterService("Generic", func(env interface{}) interface{} {
		return NewGenericServiceLocator(env.(*environments.Env))
	})
}
