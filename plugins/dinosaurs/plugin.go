package dinosaurs

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/openshift-online/rh-trex/cmd/trex/environments"
	"github.com/openshift-online/rh-trex/cmd/trex/environments/registry"
	"github.com/openshift-online/rh-trex/cmd/trex/server"
	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/api/presenters"
	"github.com/openshift-online/rh-trex/pkg/auth"
	"github.com/openshift-online/rh-trex/pkg/controllers"
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/db"
	"github.com/openshift-online/rh-trex/pkg/handlers"
	"github.com/openshift-online/rh-trex/pkg/services"
	"github.com/openshift-online/rh-trex/plugins/events"
	"github.com/openshift-online/rh-trex/plugins/generic"
)

// Service Locator
type DinosaurServiceLocator func() services.DinosaurService

func NewDinosaurServiceLocator(env *environments.Env) DinosaurServiceLocator {
	return func() services.DinosaurService {
		return services.NewDinosaurService(
			db.NewAdvisoryLockFactory(env.Database.SessionFactory),
			dao.NewDinosaurDao(&env.Database.SessionFactory),
			events.EventService(&env.Services),
		)
	}
}

// DinosaurService helper function to get the dinosaur service from the registry
func DinosaurService(s *environments.Services) services.DinosaurService {
	if s == nil {
		return nil
	}
	if obj := s.GetService("Dinosaurs"); obj != nil {
		locator := obj.(DinosaurServiceLocator)
		return locator()
	}
	return nil
}

func init() {
	// Service registration
	registry.RegisterService("Dinosaurs", func(env interface{}) interface{} {
		return NewDinosaurServiceLocator(env.(*environments.Env))
	})

	// Routes registration
	server.RegisterRoutes("dinosaurs", func(apiV1Router *mux.Router, services server.ServicesInterface, authMiddleware auth.JWTMiddleware, authzMiddleware auth.AuthorizationMiddleware) {
		envServices := services.(*environments.Services)
		dinosaurHandler := handlers.NewDinosaurHandler(DinosaurService(envServices), generic.GenericService(envServices))

		dinosaursRouter := apiV1Router.PathPrefix("/dinosaurs").Subrouter()
		dinosaursRouter.HandleFunc("", dinosaurHandler.List).Methods(http.MethodGet)
		dinosaursRouter.HandleFunc("/{id}", dinosaurHandler.Get).Methods(http.MethodGet)
		dinosaursRouter.HandleFunc("", dinosaurHandler.Create).Methods(http.MethodPost)
		dinosaursRouter.HandleFunc("/{id}", dinosaurHandler.Patch).Methods(http.MethodPatch)
		dinosaursRouter.HandleFunc("/{id}", dinosaurHandler.Delete).Methods(http.MethodDelete)
		dinosaursRouter.Use(authMiddleware.AuthenticateAccountJWT)
		dinosaursRouter.Use(authzMiddleware.AuthorizeApi)
	})

	// Controller registration
	server.RegisterController("Dinosaurs", func(manager *controllers.KindControllerManager, services *environments.Services) {
		dinoServices := DinosaurService(services)

		manager.Add(&controllers.ControllerConfig{
			Source: "Dinosaurs",
			Handlers: map[api.EventType][]controllers.ControllerHandlerFunc{
				api.CreateEventType: {dinoServices.OnUpsert},
				api.UpdateEventType: {dinoServices.OnUpsert},
				api.DeleteEventType: {dinoServices.OnDelete},
			},
		})
	})

	// Presenter registration
	presenters.RegisterPath(api.Dinosaur{}, "dinosaurs")
	presenters.RegisterPath(&api.Dinosaur{}, "dinosaurs")
	presenters.RegisterKind(api.Dinosaur{}, "Dinosaur")
	presenters.RegisterKind(&api.Dinosaur{}, "Dinosaur")
}
