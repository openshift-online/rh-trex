package server

import (
	"context"

	"github.com/openshift-online/rh-trex/cmd/trex/environments"
	"github.com/openshift-online/rh-trex/pkg/controllers"
	"github.com/openshift-online/rh-trex/pkg/db"

	"github.com/openshift-online/rh-trex/pkg/logger"
)

type ControllerRegistrationFunc func(manager *controllers.KindControllerManager, services *environments.Services)

var controllerRegistry = make(map[string]ControllerRegistrationFunc)

func RegisterController(name string, registrationFunc ControllerRegistrationFunc) {
	controllerRegistry[name] = registrationFunc
}

func LoadDiscoveredControllers(manager *controllers.KindControllerManager, services *environments.Services) {
	for name, registrationFunc := range controllerRegistry {
		registrationFunc(manager, services)
		_ = name // prevent unused variable warning
	}
}

func NewControllersServer() *ControllersServer {

	s := &ControllersServer{
		KindControllerManager: controllers.NewKindControllerManager(
			db.NewAdvisoryLockFactory(env().Database.SessionFactory),
			env().Services.Events(),
		),
	}

	// Auto-discovered controllers (no manual editing needed)
	LoadDiscoveredControllers(s.KindControllerManager, &env().Services)

	return s
}

type ControllersServer struct {
	KindControllerManager *controllers.KindControllerManager
	DB                    db.SessionFactory
}

// Start is a blocking call that starts this controller server
func (s ControllersServer) Start() {
	log := logger.NewOCMLogger(context.Background())

	log.Infof("Kind controller listening for events")

	// blocking call
	env().Database.SessionFactory.NewListener(context.Background(), "events", s.KindControllerManager.Handle)
}
