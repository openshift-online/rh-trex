package dao

import (
	"context"

	"github.com/openshift-online/rh-trex/pkg/api"
	coreapi "github.com/openshift-online/rh-trex-core/api"
	coredao "github.com/openshift-online/rh-trex-core/dao"
	"github.com/openshift-online/rh-trex/pkg/db"
)

// DinosaurCoreDAO extends the base DAO with dinosaur-specific methods
type DinosaurCoreDAO struct {
	*coredao.BaseDAO[api.Dinosaur]
	sessionFactory *db.SessionFactory
}

// NewDinosaurCoreDAO creates a new dinosaur DAO using the core framework
func NewDinosaurCoreDAO(sessionFactory *db.SessionFactory) *DinosaurCoreDAO {
	// Create a GORM DB instance from the session factory
	db := (*sessionFactory).New(context.Background())
	
	return &DinosaurCoreDAO{
		BaseDAO:        coredao.NewBaseDAO[api.Dinosaur](db),
		sessionFactory: sessionFactory,
	}
}

// FindBySpecies finds dinosaurs by species using the core framework
func (d *DinosaurCoreDAO) FindBySpecies(ctx context.Context, species string) ([]api.Dinosaur, error) {
	db := (*d.sessionFactory).New(ctx)
	
	var dinosaurs []api.Dinosaur
	err := db.Where("species = ?", species).Find(&dinosaurs).Error
	return dinosaurs, err
}

// Override methods to use session factory for transaction support
func (d *DinosaurCoreDAO) Get(ctx context.Context, id string) (*api.Dinosaur, error) {
	db := (*d.sessionFactory).New(ctx)
	var dinosaur api.Dinosaur
	err := db.Where("id = ?", id).First(&dinosaur).Error
	if err != nil {
		return nil, err
	}
	return &dinosaur, nil
}

func (d *DinosaurCoreDAO) Create(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, error) {
	db := (*d.sessionFactory).New(ctx)
	err := db.Create(dinosaur).Error
	if err != nil {
		return nil, err
	}
	return dinosaur, nil
}

func (d *DinosaurCoreDAO) Replace(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, error) {
	db := (*d.sessionFactory).New(ctx)
	err := db.Save(dinosaur).Error
	if err != nil {
		return nil, err
	}
	return dinosaur, nil
}

func (d *DinosaurCoreDAO) Delete(ctx context.Context, id string) error {
	db := (*d.sessionFactory).New(ctx)
	return db.Where("id = ?", id).Delete(&api.Dinosaur{}).Error
}

func (d *DinosaurCoreDAO) FindByIDs(ctx context.Context, ids []string) ([]api.Dinosaur, error) {
	db := (*d.sessionFactory).New(ctx)
	var dinosaurs []api.Dinosaur
	err := db.Where("id IN ?", ids).Find(&dinosaurs).Error
	return dinosaurs, err
}

func (d *DinosaurCoreDAO) All(ctx context.Context) ([]api.Dinosaur, error) {
	db := (*d.sessionFactory).New(ctx)
	var dinosaurs []api.Dinosaur
	err := db.Find(&dinosaurs).Error
	return dinosaurs, err
}

// Implement the core DAO interface methods that return the expected types
func (d *DinosaurCoreDAO) List(ctx context.Context, query coreapi.ListQuery) ([]api.Dinosaur, error) {
	db := (*d.sessionFactory).New(ctx)
	var dinosaurs []api.Dinosaur
	
	// Apply pagination if specified
	if query.Size > 0 {
		offset := (query.Page - 1) * query.Size
		db = db.Offset(offset).Limit(query.Size)
	}
	
	// Apply search if specified
	if query.Search != "" {
		db = db.Where("species ILIKE ?", "%"+query.Search+"%")
	}
	
	err := db.Find(&dinosaurs).Error
	return dinosaurs, err
}

func (d *DinosaurCoreDAO) Count(ctx context.Context, query coreapi.ListQuery) (int, error) {
	db := (*d.sessionFactory).New(ctx)
	var count int64
	
	// Apply search if specified
	if query.Search != "" {
		db = db.Where("species ILIKE ?", "%"+query.Search+"%")
	}
	
	err := db.Model(&api.Dinosaur{}).Count(&count).Error
	return int(count), err
}

// Adapter method to convert core API ListQuery to our DAO ListQuery
func (d *DinosaurCoreDAO) ListWithCoreQuery(ctx context.Context, query coreapi.ListQuery) ([]api.Dinosaur, error) {
	return d.List(ctx, query)
}

func (d *DinosaurCoreDAO) CountWithCoreQuery(ctx context.Context, query coreapi.ListQuery) (int, error) {
	return d.Count(ctx, query)
}