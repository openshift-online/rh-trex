package server

import (
	"context"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/api"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/controllers"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/db"

	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/logger"
)

func NewControllersServer() *ControllersServer {

	s := &ControllersServer{
		KindControllerManager: controllers.NewKindControllerManager(env().Services.Events()),
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
