package environments

import (
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/services"
)

type GenericServiceLocator func() services.GenericService

func NewGenericServiceLocator(env *Env) GenericServiceLocator {
	return func() services.GenericService {
		return services.NewGenericService(dao.NewGenericDao(&env.Database.SessionFactory))
	}
}
