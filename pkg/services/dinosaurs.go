package services

import (
	"context"

	"github.com/openshift-online/rh-trex/pkg/dao"
	logger "github.com/openshift-online/rh-trex/pkg/logger"

	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/errors"
)

type DinosaurService interface {
	Get(ctx context.Context, id string) (*api.Dinosaur, *errors.ServiceError)
	Create(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, *errors.ServiceError)
	Replace(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, *errors.ServiceError)
	Delete(ctx context.Context, id string) *errors.ServiceError
	All(ctx context.Context) (api.DinosaurList, *errors.ServiceError)

	FindBySpecies(ctx context.Context, species string) (api.DinosaurList, *errors.ServiceError)
	FindByIDs(ctx context.Context, ids []string) (api.DinosaurList, *errors.ServiceError)

	// idempotent functions for the control plane, but can also be called synchronously by any actor
	OnUpsert(ctx context.Context, id string) error
	OnDelete(ctx context.Context, id string) error
}

func NewDinosaurService(dinosaurDao dao.DinosaurDao, events EventService) DinosaurService {
	return &sqlDinosaurService{
		dinosaurDao: dinosaurDao,
		events:      events,
	}
}

var _ DinosaurService = &sqlDinosaurService{}

type sqlDinosaurService struct {
	dinosaurDao dao.DinosaurDao
	events      EventService
}

func (s *sqlDinosaurService) OnUpsert(ctx context.Context, id string) error {
	logger := logger.NewOCMLogger(ctx)

	dinosaur, err := s.dinosaurDao.Get(ctx, id)
	if err != nil {
		return err
	}

	logger.Infof("Do idempotent somethings with this dinosaur: %s", dinosaur.ID)

	return nil
}

func (s *sqlDinosaurService) OnDelete(ctx context.Context, id string) error {
	logger := logger.NewOCMLogger(ctx)
	logger.Infof("This dino didn't make it to the asteroid: %s", id)
	return nil
}

func (s *sqlDinosaurService) Get(ctx context.Context, id string) (*api.Dinosaur, *errors.ServiceError) {
	dinosaur, err := s.dinosaurDao.Get(ctx, id)
	if err != nil {
		return nil, handleGetError("Dinosaur", "id", id, err)
	}
	return dinosaur, nil
}

func (s *sqlDinosaurService) Create(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, *errors.ServiceError) {
	dinosaur, err := s.dinosaurDao.Create(ctx, dinosaur)
	if err != nil {
		return nil, handleCreateError("Dinosaur", err)
	}

	_, eErr := s.events.Create(ctx, &api.Event{
		Source:    "Dinosaurs",
		SourceID:  dinosaur.ID,
		EventType: api.CreateEventType,
	})
	if eErr != nil {
		return nil, handleCreateError("Dinosaur", err)
	}

	return dinosaur, nil
}

func (s *sqlDinosaurService) Replace(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, *errors.ServiceError) {
	dinosaur, err := s.dinosaurDao.Replace(ctx, dinosaur)
	if err != nil {
		return nil, handleUpdateError("Dinosaur", err)
	}

	_, eErr := s.events.Create(ctx, &api.Event{
		Source:    "Dinosaurs",
		SourceID:  dinosaur.ID,
		EventType: api.UpdateEventType,
	})
	if eErr != nil {
		return nil, handleUpdateError("Dinosaur", err)
	}
	return dinosaur, nil
}

func (s *sqlDinosaurService) Delete(ctx context.Context, id string) *errors.ServiceError {
	if err := s.dinosaurDao.Delete(ctx, id); err != nil {
		return handleDeleteError("Dinosaur", errors.GeneralError("Unable to delete dinosaur: %s", err))
	}

	_, err := s.events.Create(ctx, &api.Event{
		Source:    "Dinosaurs",
		SourceID:  id,
		EventType: api.DeleteEventType,
	})
	if err != nil {
		return handleDeleteError("Dinosaur", err)
	}

	return nil
}

func (s *sqlDinosaurService) FindByIDs(ctx context.Context, ids []string) (api.DinosaurList, *errors.ServiceError) {
	dinosaurs, err := s.dinosaurDao.FindByIDs(ctx, ids)
	if err != nil {
		return nil, errors.GeneralError("Unable to get all dinosaurs: %s", err)
	}
	return dinosaurs, nil
}

func (s *sqlDinosaurService) FindBySpecies(ctx context.Context, species string) (api.DinosaurList, *errors.ServiceError) {
	dinosaurs, err := s.dinosaurDao.FindBySpecies(ctx, species)
	if err != nil {
		return nil, handleGetError("Dinosaur", "species", species, err)
	}
	return dinosaurs, nil
}

func (s *sqlDinosaurService) All(ctx context.Context) (api.DinosaurList, *errors.ServiceError) {
	dinosaurs, err := s.dinosaurDao.All(ctx)
	if err != nil {
		return nil, errors.GeneralError("Unable to get all dinosaurs: %s", err)
	}
	return dinosaurs, nil
}
