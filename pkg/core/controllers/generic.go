package controllers

import (
	"context"
	"sync"

	"github.com/openshift-online/rh-trex/pkg/core/api"
	"github.com/openshift-online/rh-trex/pkg/core/services"
	"github.com/openshift-online/rh-trex/pkg/db"
)

// ControllerHandlerFunc defines the signature for controller event handlers
type ControllerHandlerFunc func(ctx context.Context, id string) error

// ControllerConfig defines configuration for a controller
type ControllerConfig struct {
	Source   string
	Handlers map[api.EventType][]ControllerHandlerFunc
}

// Controller handles events for a specific resource type
type Controller struct {
	Source   string
	Handlers map[api.EventType][]ControllerHandlerFunc
}

// NewController creates a new controller
func NewController(config *ControllerConfig) *Controller {
	return &Controller{
		Source:   config.Source,
		Handlers: config.Handlers,
	}
}

// HandleEvent processes an event through registered handlers
func (c *Controller) HandleEvent(ctx context.Context, event *api.Event) error {
	handlers, exists := c.Handlers[event.EventType]
	if !exists {
		// No handlers registered - could log this
		return nil
	}

	for _, handler := range handlers {
		if err := handler(ctx, event.SourceID); err != nil {
			// Handler failed - could log this
			return err
		}
	}

	return nil
}

// ControllerManager manages multiple controllers
type ControllerManager struct {
	controllers map[string]*Controller
	lockFactory db.LockFactory
	eventBus    EventBus
	mu          sync.RWMutex
}

// EventBus defines the interface for event handling
type EventBus interface {
	Subscribe(source string, handler func(ctx context.Context, event *api.Event) error)
	Publish(ctx context.Context, event *api.Event) error
}

// NewControllerManager creates a new controller manager
func NewControllerManager(lockFactory db.LockFactory, eventBus EventBus) *ControllerManager {
	return &ControllerManager{
		controllers: make(map[string]*Controller),
		lockFactory: lockFactory,
		eventBus:    eventBus,
	}
}

// RegisterController registers a new controller
func (cm *ControllerManager) RegisterController(config *ControllerConfig) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	controller := NewController(config)
	cm.controllers[config.Source] = controller

	// Subscribe to events for this source
	if cm.eventBus != nil {
		cm.eventBus.Subscribe(config.Source, controller.HandleEvent)
	}

	// Controller registered - could log this
}

// AutoRegisterCRUDController automatically registers CRUD handlers for a service
func AutoRegisterCRUDController[T any](
	cm *ControllerManager,
	service services.CRUDService[T],
	source string,
) {
	config := &ControllerConfig{
		Source: source,
		Handlers: map[api.EventType][]ControllerHandlerFunc{
			api.CreateEventType: {service.OnUpsert},
			api.UpdateEventType: {service.OnUpsert},
			api.DeleteEventType: {service.OnDelete},
		},
	}

	cm.RegisterController(config)
}

// HandleEvent processes an event through the appropriate controller
func (cm *ControllerManager) HandleEvent(ctx context.Context, event *api.Event) error {
	cm.mu.RLock()
	controller, exists := cm.controllers[event.Source]
	cm.mu.RUnlock()

	if !exists {
		// No controller registered - could log this
		return nil
	}

	return controller.HandleEvent(ctx, event)
}

// ListControllers returns all registered controllers
func (cm *ControllerManager) ListControllers() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	sources := make([]string, 0, len(cm.controllers))
	for source := range cm.controllers {
		sources = append(sources, source)
	}
	return sources
}

// GetController returns a controller by source
func (cm *ControllerManager) GetController(source string) *Controller {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return cm.controllers[source]
}

// Start begins processing events (blocking call)
func (cm *ControllerManager) Start(ctx context.Context) {
	// Starting controller manager - could log this
	
	// This would typically start listening for database events
	// Implementation depends on the specific event bus used
	select {
	case <-ctx.Done():
		// Controller manager stopping - could log this
	}
}