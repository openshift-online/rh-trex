package server

import (
	"context"

	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/controllers"
	"github.com/openshift-online/rh-trex/pkg/db"

	"github.com/openshift-online/rh-trex/pkg/logger"
)

func NewControllersServer() *ControllersServer {

	s := &ControllersServer{
		KindControllerManager: controllers.NewKindControllerManager(
			db.NewAdvisoryLockFactory(env().Database.SessionFactory),
			env().Services.Events(),
		),
	}

	dinoServices := env().Services.Dinosaurs()

	s.KindControllerManager.Add(&controllers.ControllerConfig{
		Source: "Dinosaurs",
		Handlers: map[api.EventType][]controllers.ControllerHandlerFunc{
			api.CreateEventType: {dinoServices.OnUpsert},
			api.UpdateEventType: {dinoServices.OnUpsert},
			api.DeleteEventType: {dinoServices.OnDelete},
		},
	})





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
