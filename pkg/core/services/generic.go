package services

import (
	"context"

	"github.com/openshift-online/rh-trex/pkg/core/api"
	"github.com/openshift-online/rh-trex/pkg/errors"
)

// CRUDService defines the standard CRUD operations for any resource type
type CRUDService[T any] interface {
	Get(ctx context.Context, id string) (*T, *errors.ServiceError)
	Create(ctx context.Context, obj *T) (*T, *errors.ServiceError)
	Replace(ctx context.Context, obj *T) (*T, *errors.ServiceError)
	Delete(ctx context.Context, id string) *errors.ServiceError
	List(ctx context.Context, query api.ListQuery) (*api.List, *errors.ServiceError)
	FindByIDs(ctx context.Context, ids []string) ([]T, *errors.ServiceError)

	// Event handlers for controller framework
	OnUpsert(ctx context.Context, id string) error
	OnDelete(ctx context.Context, id string) error
}

// DAO defines the data access interface
type DAO[T any] interface {
	Get(ctx context.Context, id string) (*T, error)
	Create(ctx context.Context, obj *T) (*T, error)
	Replace(ctx context.Context, obj *T) (*T, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, query api.ListQuery) ([]T, error)
	Count(ctx context.Context, query api.ListQuery) (int, error)
	FindByIDs(ctx context.Context, ids []string) ([]T, error)
}

// EventEmitter handles event creation
type EventEmitter interface {
	EmitEvent(ctx context.Context, source, sourceID string, eventType api.EventType) error
}

// BaseCRUDService provides common CRUD functionality
type BaseCRUDService[T any] struct {
	dao    DAO[T]
	events EventEmitter
	source string
}

// NewBaseCRUDService creates a new base CRUD service
func NewBaseCRUDService[T any](dao DAO[T], events EventEmitter, source string) *BaseCRUDService[T] {
	return &BaseCRUDService[T]{
		dao:    dao,
		events: events,
		source: source,
	}
}

// Get retrieves a resource by ID
func (s *BaseCRUDService[T]) Get(ctx context.Context, id string) (*T, *errors.ServiceError) {
	obj, err := s.dao.Get(ctx, id)
	if err != nil {
		return nil, handleGetError(s.source, "id", id, err)
	}
	return obj, nil
}

// Create creates a new resource
func (s *BaseCRUDService[T]) Create(ctx context.Context, obj *T) (*T, *errors.ServiceError) {
	created, err := s.dao.Create(ctx, obj)
	if err != nil {
		return nil, handleCreateError(s.source, err)
	}

	// Emit creation event
	if s.events != nil {
		if meta := extractMeta(created); meta != nil {
			s.events.EmitEvent(ctx, s.source, meta.ID, api.CreateEventType)
		}
	}

	return created, nil
}

// Replace updates an existing resource
func (s *BaseCRUDService[T]) Replace(ctx context.Context, obj *T) (*T, *errors.ServiceError) {
	updated, err := s.dao.Replace(ctx, obj)
	if err != nil {
		return nil, handleUpdateError(s.source, err)
	}

	// Emit update event
	if s.events != nil {
		if meta := extractMeta(updated); meta != nil {
			s.events.EmitEvent(ctx, s.source, meta.ID, api.UpdateEventType)
		}
	}

	return updated, nil
}

// Delete removes a resource
func (s *BaseCRUDService[T]) Delete(ctx context.Context, id string) *errors.ServiceError {
	if err := s.dao.Delete(ctx, id); err != nil {
		return handleDeleteError(s.source, err)
	}

	// Emit deletion event
	if s.events != nil {
		s.events.EmitEvent(ctx, s.source, id, api.DeleteEventType)
	}

	return nil
}

// List retrieves resources with pagination
func (s *BaseCRUDService[T]) List(ctx context.Context, query api.ListQuery) (*api.List, *errors.ServiceError) {
	items, err := s.dao.List(ctx, query)
	if err != nil {
		return nil, errors.GeneralError("Unable to list %s: %s", s.source, err)
	}

	count, err := s.dao.Count(ctx, query)
	if err != nil {
		return nil, errors.GeneralError("Unable to count %s: %s", s.source, err)
	}

	return &api.List{
		Kind:  s.source + "List",
		Page:  query.Page,
		Size:  query.Size,
		Total: count,
		Items: items,
	}, nil
}

// FindByIDs finds resources by their IDs
func (s *BaseCRUDService[T]) FindByIDs(ctx context.Context, ids []string) ([]T, *errors.ServiceError) {
	items, err := s.dao.FindByIDs(ctx, ids)
	if err != nil {
		return nil, errors.GeneralError("Unable to find %s by IDs: %s", s.source, err)
	}
	return items, nil
}

// OnUpsert handles CREATE and UPDATE events (override in concrete implementations)
func (s *BaseCRUDService[T]) OnUpsert(ctx context.Context, id string) error {
	obj, err := s.dao.Get(ctx, id)
	if err != nil {
		return err
	}

	// Log processing (would use injected logger in real implementation)
	return s.processUpsert(ctx, obj)
}

// OnDelete handles DELETE events (override in concrete implementations)
func (s *BaseCRUDService[T]) OnDelete(ctx context.Context, id string) error {
	// Log processing (would use injected logger in real implementation)
	return s.processDelete(ctx, id)
}

// processUpsert can be overridden in concrete implementations for custom business logic
func (s *BaseCRUDService[T]) processUpsert(ctx context.Context, obj *T) error {
	// Default: no-op
	return nil
}

// processDelete can be overridden in concrete implementations for custom business logic
func (s *BaseCRUDService[T]) processDelete(ctx context.Context, id string) error {
	// Default: no-op
	return nil
}

// Helper functions
func extractMeta(obj interface{}) *api.Meta {
	// Use reflection or type assertion to extract Meta field
	// This is a simplified version - in practice, you'd use reflection
	type MetaProvider interface {
		GetMeta() *api.Meta
	}
	
	if provider, ok := obj.(MetaProvider); ok {
		return provider.GetMeta()
	}
	return nil
}

// Error handling helpers (these would be moved to the errors package)
func handleGetError(kind, field, value string, err error) *errors.ServiceError {
	return errors.NotFound("%s with %s='%s' not found", kind, field, value)
}

func handleCreateError(kind string, err error) *errors.ServiceError {
	return errors.GeneralError("Unable to create %s: %s", kind, err)
}

func handleUpdateError(kind string, err error) *errors.ServiceError {
	return errors.GeneralError("Unable to update %s: %s", kind, err)
}

func handleDeleteError(kind string, err error) *errors.ServiceError {
	return errors.GeneralError("Unable to delete %s: %s", kind, err)
}