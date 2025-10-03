package services

import (
	"context"

	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/errors"
)

type EventService interface {
	Get(ctx context.Context, id string) (*api.Event, *errors.ServiceError)
	Create(ctx context.Context, event *api.Event) (*api.Event, *errors.ServiceError)
	Replace(ctx context.Context, event *api.Event) (*api.Event, *errors.ServiceError)
	Delete(ctx context.Context, id string) *errors.ServiceError
	All(ctx context.Context) (api.EventList, *errors.ServiceError)

	FindByIDs(ctx context.Context, ids []string) (api.EventList, *errors.ServiceError)
}

func NewEventService(eventDao dao.EventDao) EventService {
	return &sqlEventService{
		eventDao: eventDao,
	}
}

var _ EventService = &sqlEventService{}

type sqlEventService struct {
	eventDao dao.EventDao
}

func (s *sqlEventService) Get(ctx context.Context, id string) (*api.Event, *errors.ServiceError) {
	event, err := s.eventDao.Get(ctx, id)
	if err != nil {
		return nil, handleGetError("Event", "id", id, err)
	}
	return event, nil
}

func (s *sqlEventService) Create(ctx context.Context, event *api.Event) (*api.Event, *errors.ServiceError) {
	event, err := s.eventDao.Create(ctx, event)
	if err != nil {
		return nil, handleCreateError("Event", err)
	}
	return event, nil
}

func (s *sqlEventService) Replace(ctx context.Context, event *api.Event) (*api.Event, *errors.ServiceError) {
	event, err := s.eventDao.Replace(ctx, event)
	if err != nil {
		return nil, handleUpdateError("Event", err)
	}
	return event, nil
}

func (s *sqlEventService) Delete(ctx context.Context, id string) *errors.ServiceError {
	if err := s.eventDao.Delete(ctx, id); err != nil {
		return handleDeleteError("Event", err)
	}
	return nil
}

func (s *sqlEventService) FindByIDs(ctx context.Context, ids []string) (api.EventList, *errors.ServiceError) {
	events, err := s.eventDao.FindByIDs(ctx, ids)
	if err != nil {
		return nil, handleGetError("Event", "list", "", err)
	}
	return events, nil
}

func (s *sqlEventService) All(ctx context.Context) (api.EventList, *errors.ServiceError) {
	events, err := s.eventDao.All(ctx)
	if err != nil {
		return nil, handleGetError("Event", "list", "", err)
	}
	return events, nil
}
