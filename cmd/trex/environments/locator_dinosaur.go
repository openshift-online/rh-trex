package environments

import (
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/db"
	"github.com/openshift-online/rh-trex/pkg/services"
)

type DinosaurServiceLocator func() services.DinosaurService

func NewDinosaurServiceLocator(env *Env) DinosaurServiceLocator {
	return func() services.DinosaurService {
		return services.NewDinosaurService(
			db.NewAdvisoryLockFactory(env.Database.SessionFactory),
			dao.NewDinosaurDao(&env.Database.SessionFactory),
			env.Services.Events(),
		)
	}
}
