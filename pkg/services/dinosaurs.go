package services

import (
	"context"

	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/db"
	"github.com/openshift-online/rh-trex/pkg/logger"

	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/errors"
)

// This flag will only be used in integration test to prove that the advisory lock works
var DisableAdvisoryLock = false

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

func NewDinosaurService(lockFactory db.LockFactory, dinosaurDao dao.DinosaurDao, events EventService, sessionFactory db.SessionFactory) DinosaurService {
	return &sqlDinosaurService{
		lockFactory:    lockFactory,
		dinosaurDao:    dinosaurDao,
		events:         events,
		sessionFactory: sessionFactory,
	}
}

var _ DinosaurService = &sqlDinosaurService{}

type sqlDinosaurService struct {
	lockFactory    db.LockFactory
	dinosaurDao    dao.DinosaurDao
	events         EventService
	sessionFactory db.SessionFactory
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
	var createdDinosaur *api.Dinosaur
	var serviceErr *errors.ServiceError

	// Execute within a transaction
	txErr := db.WithTransaction(ctx, s.sessionFactory, func(txCtx context.Context) error {
		var err error
		createdDinosaur, serviceErr = s.create(txCtx, dinosaur)
		if serviceErr != nil {
			return serviceErr
		}
		return err
	})

	if txErr != nil {
		if serviceErr != nil {
			return nil, serviceErr
		}
		return nil, handleCreateError("Dinosaur", txErr)
	}

	return createdDinosaur, nil
}

// create contains the business logic for creating a dinosaur and its associated event.
// This internal function is called within a transaction context.
func (s *sqlDinosaurService) create(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, *errors.ServiceError) {
	createdDinosaur, err := s.dinosaurDao.Create(ctx, dinosaur)
	if err != nil {
		return nil, handleCreateError("Dinosaur", err)
	}

	if createdDinosaur != nil {
	}

	_, eventErr := s.events.Create(ctx, &api.Event{
		Source:    "Dinosaurs",
		SourceID:  createdDinosaur.ID,
		EventType: api.CreateEventType,
	})
	if eventErr != nil {
		return nil, eventErr
	}

	return createdDinosaur, nil
}

func (s *sqlDinosaurService) Replace(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, *errors.ServiceError) {
	var updatedDinosaur *api.Dinosaur
	var serviceErr *errors.ServiceError

	// Execute within a transaction
	txErr := db.WithTransaction(ctx, s.sessionFactory, func(txCtx context.Context) error {
		var err error
		updatedDinosaur, serviceErr = s.replace(txCtx, dinosaur)
		if serviceErr != nil {
			return serviceErr
		}
		return err
	})

	if txErr != nil {
		if serviceErr != nil {
			return nil, serviceErr
		}
		return nil, handleUpdateError("Dinosaur", txErr)
	}

	return updatedDinosaur, nil
}

// replace contains the business logic for updating a dinosaur and its associated event.
// This internal function is called within a transaction context.
func (s *sqlDinosaurService) replace(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, *errors.ServiceError) {
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

	found, err := s.dinosaurDao.Get(ctx, dinosaur.ID)
	if err != nil {
		return nil, handleGetError("Dinosaur", "id", dinosaur.ID, err)
	}

	// New species is no change, the update action is not needed.
	if found.Species == dinosaur.Species {
		return found, nil
	}

	found.Species = dinosaur.Species
	updatedDinosaur, err := s.dinosaurDao.Replace(ctx, found)
	if err != nil {
		return nil, handleUpdateError("Dinosaur", err)
	}

	_, eventErr := s.events.Create(ctx, &api.Event{
		Source:    "Dinosaurs",
		SourceID:  updatedDinosaur.ID,
		EventType: api.UpdateEventType,
	})
	if eventErr != nil {
		return nil, eventErr
	}

	return updatedDinosaur, nil
}

func (s *sqlDinosaurService) Delete(ctx context.Context, id string) *errors.ServiceError {
	var serviceErr *errors.ServiceError

	// Execute within a transaction
	txErr := db.WithTransaction(ctx, s.sessionFactory, func(txCtx context.Context) error {
		serviceErr = s.delete(txCtx, id)
		if serviceErr != nil {
			return serviceErr
		}
		return nil
	})

	if txErr != nil {
		if serviceErr != nil {
			return serviceErr
		}
		return handleDeleteError("Dinosaur", errors.GeneralError("Unable to delete dinosaur: %s", txErr))
	}

	return nil
}

// delete contains the business logic for deleting a dinosaur and its associated event.
// This internal function is called within a transaction context.
func (s *sqlDinosaurService) delete(ctx context.Context, id string) *errors.ServiceError {
	err := s.dinosaurDao.Delete(ctx, id)
	if err != nil {
		return handleDeleteError("Dinosaur", errors.GeneralError("Unable to delete dinosaur: %s", err))
	}

	_, eventErr := s.events.Create(ctx, &api.Event{
		Source:    "Dinosaurs",
		SourceID:  id,
		EventType: api.DeleteEventType,
	})
	if eventErr != nil {
		return eventErr
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
