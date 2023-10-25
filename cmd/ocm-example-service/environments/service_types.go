package environments

import (
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/dao"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/services"
)

type DinosaurServiceLocator func() services.DinosaurService

func NewDinosaurServiceLocator(env *Env) DinosaurServiceLocator {
	return func() services.DinosaurService {
		return services.NewDinosaurService(dao.NewDinosaurDao(&env.Database.SessionFactory), env.Services.Events())
	}
}

type GenericServiceLocator func() services.GenericService

func NewGenericServiceLocator(env *Env) GenericServiceLocator {
	return func() services.GenericService {
		return services.NewGenericService(dao.NewGenericDao(&env.Database.SessionFactory))
	}
}

type EventServiceLocator func() services.EventService

func NewEventServiceLocator(env *Env) EventServiceLocator {
	return func() services.EventService {
		return services.NewEventService(dao.NewEventDao(&env.Database.SessionFactory))
	}
}
