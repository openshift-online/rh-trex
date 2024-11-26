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

type GenericServiceLocator func() services.GenericService

func NewGenericServiceLocator(env *Env) GenericServiceLocator {
	return func() services.GenericService {
		return services.NewGenericService(
			dao.NewGenericDao(&env.Database.SessionFactory),
			env.Clients.OCM)
	}
}

type EventServiceLocator func() services.EventService

func NewEventServiceLocator(env *Env) EventServiceLocator {
	return func() services.EventService {
		return services.NewEventService(dao.NewEventDao(&env.Database.SessionFactory))
	}
}
