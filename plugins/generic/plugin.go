package generic

import (
	"github.com/openshift-online/rh-trex/cmd/trex/environments"
	"github.com/openshift-online/rh-trex/cmd/trex/environments/registry"
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/services"
)

// ServiceLocator Service Locator
type ServiceLocator func() services.GenericService

func NewServiceLocator(env *environments.Env) ServiceLocator {
	return func() services.GenericService {
		return services.NewGenericService(dao.NewGenericDao(&env.Database.SessionFactory))
	}
}

// Service helper function to get the generic service from the registry
func Service(s *environments.Services) services.GenericService {
	if s == nil {
		return nil
	}
	if obj := s.GetService("Generic"); obj != nil {
		locator := obj.(ServiceLocator)
		return locator()
	}
	return nil
}

func init() {
	// Service registration
	registry.RegisterService("Generic", func(env interface{}) interface{} {
		return NewServiceLocator(env.(*environments.Env))
	})
}
