package environments

import (
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/services"
)

type EventServiceLocator func() services.EventService

func NewEventServiceLocator(env *Env) EventServiceLocator {
	return func() services.EventService {
		return services.NewEventService(dao.NewEventDao(&env.Database.SessionFactory))
	}
}
