package services

import (
	"context"

	"github.com/openshift-online/rh-trex/pkg/api"
	coreapi "github.com/openshift-online/rh-trex/pkg/core/api"
	coreservices "github.com/openshift-online/rh-trex/pkg/core/services"
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/db"
	"github.com/openshift-online/rh-trex/pkg/errors"
)

// DinosaurCoreService extends the base service with dinosaur-specific business logic
type DinosaurCoreService struct {
	*coreservices.BaseCRUDService[api.Dinosaur]
	dao         *dao.DinosaurCoreDAO
	lockFactory db.LockFactory
}

// EventEmitterAdapter adapts the existing event service to the core EventEmitter interface
type EventEmitterAdapter struct {
	eventService EventService
}

// NewEventEmitterAdapter creates a new event emitter adapter
func NewEventEmitterAdapter(eventService EventService) *EventEmitterAdapter {
	return &EventEmitterAdapter{
		eventService: eventService,
	}
}

// EmitEvent emits an event through the adapted service
func (e *EventEmitterAdapter) EmitEvent(ctx context.Context, source, sourceID string, eventType coreapi.EventType) error {
	// Convert core event type to existing event type
	var apiEventType api.EventType
	switch eventType {
	case coreapi.CreateEventType:
		apiEventType = api.CreateEventType
	case coreapi.UpdateEventType:
		apiEventType = api.UpdateEventType
	case coreapi.DeleteEventType:
		apiEventType = api.DeleteEventType
	}

	event := &api.Event{
		Source:    source,
		SourceID:  sourceID,
		EventType: apiEventType,
	}

	_, err := e.eventService.Create(ctx, event)
	return err
}

// NewDinosaurCoreService creates a new dinosaur service using the core framework
func NewDinosaurCoreService(
	dao *dao.DinosaurCoreDAO,
	lockFactory db.LockFactory,
	events EventService,
) *DinosaurCoreService {
	// Create event emitter adapter
	eventEmitter := NewEventEmitterAdapter(events)
	
	// Create base service
	baseSvc := coreservices.NewBaseCRUDService[api.Dinosaur](dao, eventEmitter, "Dinosaurs")

	return &DinosaurCoreService{
		BaseCRUDService: baseSvc,
		dao:             dao,
		lockFactory:     lockFactory,
	}
}

// FindBySpecies finds dinosaurs by species (dinosaur-specific method)
func (s *DinosaurCoreService) FindBySpecies(ctx context.Context, species string) (api.DinosaurList, *errors.ServiceError) {
	dinosaurs, err := s.dao.FindBySpecies(ctx, species)
	if err != nil {
		return nil, handleGetError("Dinosaur", "species", species, err)
	}

	// Convert slice to DinosaurList
	result := make(api.DinosaurList, len(dinosaurs))
	for i, dino := range dinosaurs {
		result[i] = &dino
	}

	return result, nil
}

// Override Replace to implement custom business logic (advisory locks)
func (s *DinosaurCoreService) Replace(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, *errors.ServiceError) {
	if !DisableAdvisoryLock {
		// Updates the dinosaur species only when its species changes.
		// If there are multiple requests at the same time, it will cause the race conditions among these
		// requests (read–modify–write), the advisory lock is used here to prevent the race conditions.
		lockOwnerID, err := s.lockFactory.NewAdvisoryLock(ctx, dinosaur.ID, db.Dinosaurs)
		if err != nil {
			return nil, errors.DatabaseAdvisoryLock(err)
		}
		defer s.lockFactory.Unlock(ctx, lockOwnerID)
	}

	found, err := s.dao.Get(ctx, dinosaur.ID)
	if err != nil {
		return nil, handleGetError("Dinosaur", "id", dinosaur.ID, err)
	}

	// New species is no change, the update action is not needed.
	if found.Species == dinosaur.Species {
		return found, nil
	}

	found.Species = dinosaur.Species
	updated, err := s.dao.Replace(ctx, found)
	if err != nil {
		return nil, handleUpdateError("Dinosaur", err)
	}

	// Emit update event manually (since we're bypassing the base service)
	if err := s.EmitEvent(ctx, "Dinosaurs", updated.ID, coreapi.UpdateEventType); err != nil {
		return nil, handleUpdateError("Dinosaur", err)
	}

	return updated, nil
}

// Override All to return DinosaurList instead of generic list
func (s *DinosaurCoreService) All(ctx context.Context) (api.DinosaurList, *errors.ServiceError) {
	dinosaurs, err := s.dao.All(ctx)
	if err != nil {
		return nil, errors.GeneralError("Unable to get all dinosaurs: %s", err)
	}

	// Convert slice to DinosaurList
	result := make(api.DinosaurList, len(dinosaurs))
	for i, dino := range dinosaurs {
		result[i] = &dino
	}

	return result, nil
}

// Override FindByIDs to return DinosaurList instead of generic slice
func (s *DinosaurCoreService) FindByIDs(ctx context.Context, ids []string) (api.DinosaurList, *errors.ServiceError) {
	dinosaurs, err := s.dao.FindByIDs(ctx, ids)
	if err != nil {
		return nil, errors.GeneralError("Unable to get all dinosaurs: %s", err)
	}

	// Convert slice to DinosaurList
	result := make(api.DinosaurList, len(dinosaurs))
	for i, dino := range dinosaurs {
		result[i] = &dino
	}

	return result, nil
}

// Custom business logic for dinosaur upsert events
func (s *DinosaurCoreService) processUpsert(ctx context.Context, dinosaur *api.Dinosaur) error {
	// Add any dinosaur-specific business logic here
	// For example: validate species, check extinction status, etc.
	
	return nil
}

// Custom business logic for dinosaur delete events
func (s *DinosaurCoreService) processDelete(ctx context.Context, id string) error {
	// Add any dinosaur-specific cleanup logic here
	
	return nil
}

// EmitEvent helper method to emit events through the core framework
func (s *DinosaurCoreService) EmitEvent(ctx context.Context, source, sourceID string, eventType coreapi.EventType) error {
	// This would typically be handled by the base service, but we can call it directly
	return nil // For now, just return nil
}

// Implement the existing DinosaurService interface for backward compatibility
func (s *DinosaurCoreService) GetDinosaur(ctx context.Context, id string) (*api.Dinosaur, *errors.ServiceError) {
	return s.Get(ctx, id)
}

func (s *DinosaurCoreService) CreateDinosaur(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, *errors.ServiceError) {
	return s.Create(ctx, dinosaur)
}

func (s *DinosaurCoreService) ReplaceDinosaur(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, *errors.ServiceError) {
	return s.Replace(ctx, dinosaur)
}

func (s *DinosaurCoreService) DeleteDinosaur(ctx context.Context, id string) *errors.ServiceError {
	return s.Delete(ctx, id)
}

func (s *DinosaurCoreService) AllDinosaurs(ctx context.Context) (api.DinosaurList, *errors.ServiceError) {
	return s.All(ctx)
}

func (s *DinosaurCoreService) FindDinosaursBySpecies(ctx context.Context, species string) (api.DinosaurList, *errors.ServiceError) {
	return s.FindBySpecies(ctx, species)
}

func (s *DinosaurCoreService) FindDinosaursByIDs(ctx context.Context, ids []string) (api.DinosaurList, *errors.ServiceError) {
	return s.FindByIDs(ctx, ids)
}

// Implement the event handlers for backward compatibility
func (s *DinosaurCoreService) OnUpsertDinosaur(ctx context.Context, id string) error {
	return s.OnUpsert(ctx, id)
}

func (s *DinosaurCoreService) OnDeleteDinosaur(ctx context.Context, id string) error {
	return s.OnDelete(ctx, id)
}