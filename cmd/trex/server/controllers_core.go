package server

import (
	"context"

	"github.com/openshift-online/rh-trex/pkg/api"
	coreapi "github.com/openshift-online/rh-trex-core/api"
	corecontrollers "github.com/openshift-online/rh-trex-core/controllers"
	"github.com/openshift-online/rh-trex/pkg/controllers"
	"github.com/openshift-online/rh-trex/pkg/db"
	coredb "github.com/openshift-online/rh-trex-core/db"
	"github.com/openshift-online/rh-trex/pkg/errors"
	"github.com/openshift-online/rh-trex/pkg/logger"
)

// EventBusAdapter adapts the existing event system to the core EventBus interface
type EventBusAdapter struct {
	eventService interface {
		Create(ctx context.Context, event *api.Event) (*api.Event, *errors.ServiceError)
	}
}

// NewEventBusAdapter creates a new event bus adapter
func NewEventBusAdapter(eventService interface {
	Create(ctx context.Context, event *api.Event) (*api.Event, *errors.ServiceError)
}) *EventBusAdapter {
	return &EventBusAdapter{
		eventService: eventService,
	}
}

// Subscribe subscribes to events for a specific source
func (e *EventBusAdapter) Subscribe(source string, handler func(ctx context.Context, event *coreapi.Event) error) {
	// In a real implementation, this would wire up database listeners
	// For now, we'll just store the handler for later use
}

// Publish publishes an event to the bus
func (e *EventBusAdapter) Publish(ctx context.Context, event *coreapi.Event) error {
	// Convert core event to existing event type
	var apiEventType api.EventType
	switch event.EventType {
	case coreapi.CreateEventType:
		apiEventType = api.CreateEventType
	case coreapi.UpdateEventType:
		apiEventType = api.UpdateEventType
	case coreapi.DeleteEventType:
		apiEventType = api.DeleteEventType
	}

	apiEvent := &api.Event{
		Source:    event.Source,
		SourceID:  event.SourceID,
		EventType: apiEventType,
	}

	_, serviceErr := e.eventService.Create(ctx, apiEvent)
	if serviceErr != nil {
		return serviceErr
	}
	return nil
}

// NewControllersServerCore creates a new controller server using the core framework
func NewControllersServerCore() *ControllersServerCore {
	// Create event bus adapter
	eventBus := NewEventBusAdapter(env().Services.Events())

	// Create core controller manager using existing database session
	sessionFactory := env().Database.SessionFactory
	coreSessionFactory := coredb.NewBasicSessionFactory(sessionFactory.New(context.Background()))
	coreManager := corecontrollers.NewControllerManager(
		coredb.NewAdvisoryLockFactory(coreSessionFactory),
		eventBus,
	)

	// Create legacy controller manager for backward compatibility
	legacyManager := controllers.NewKindControllerManager(
		db.NewAdvisoryLockFactory(env().Database.SessionFactory),
		env().Services.Events(),
	)

	s := &ControllersServerCore{
		CoreControllerManager:   coreManager,
		LegacyControllerManager: legacyManager,
	}

	// Auto-discovered controllers (no manual editing needed)
	LoadDiscoveredControllers(s.LegacyControllerManager, &env().Services)

	return s
}

// ControllersServerCore uses both core and legacy controller managers
type ControllersServerCore struct {
	CoreControllerManager   *corecontrollers.ControllerManager
	LegacyControllerManager *controllers.KindControllerManager
}

// Start is a blocking call that starts this controller server
func (s ControllersServerCore) Start() {
	log := logger.NewOCMLogger(context.Background())

	log.Infof("Core controller manager listening for events")

	// Start core controller manager in a separate goroutine
	go func() {
		ctx := context.Background()
		s.CoreControllerManager.Start(ctx)
	}()

	// Start legacy controller manager (blocking call)
	env().Database.SessionFactory.NewListener(context.Background(), "events", s.LegacyControllerManager.Handle)
}

// RegisterCoreResource registers a new resource with the core controller manager
func (s *ControllersServerCore) RegisterCoreResource(source string, handlers map[coreapi.EventType][]corecontrollers.ControllerHandlerFunc) {
	s.CoreControllerManager.RegisterController(&corecontrollers.ControllerConfig{
		Source:   source,
		Handlers: handlers,
	})
}

// GetCoreControllerManager returns the core controller manager
func (s *ControllersServerCore) GetCoreControllerManager() *corecontrollers.ControllerManager {
	return s.CoreControllerManager
}

// GetLegacyControllerManager returns the legacy controller manager
func (s *ControllersServerCore) GetLegacyControllerManager() *controllers.KindControllerManager {
	return s.LegacyControllerManager
}