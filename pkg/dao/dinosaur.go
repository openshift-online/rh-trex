package dao

import (
	"context"

	"gorm.io/gorm/clause"

	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/db"
	"github.com/openshift-online/rh-trex/pkg/util"
)

var (
	dinosaurTableName = util.ToSnakeCase(api.DinosaurTypeName) + "s"
	dinosaurColumns   = []string{
		"id",
		"created_at",
		"updated_at",
		"species",
	}
)

func DinosaurApiToModel() TableMappingRelation {
	result := map[string]string{}
	applyBaseMapping(result, dinosaurColumns, dinosaurTableName)
	return TableMappingRelation{
		Mapping:           result,
		relationTableName: dinosaurTableName,
	}
}

type DinosaurDao interface {
	Get(ctx context.Context, id string) (*api.Dinosaur, error)
	Create(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, error)
	Replace(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, error)
	Delete(ctx context.Context, id string) error
	FindByIDs(ctx context.Context, ids []string) (api.DinosaurList, error)
	FindBySpecies(ctx context.Context, species string) (api.DinosaurList, error)
	All(ctx context.Context) (api.DinosaurList, error)
}

var _ DinosaurDao = &sqlDinosaurDao{}

type sqlDinosaurDao struct {
	sessionFactory *db.SessionFactory
}

func NewDinosaurDao(sessionFactory *db.SessionFactory) DinosaurDao {
	return &sqlDinosaurDao{sessionFactory: sessionFactory}
}

func (d *sqlDinosaurDao) Get(ctx context.Context, id string) (*api.Dinosaur, error) {
	g2 := (*d.sessionFactory).New(ctx)
	var dinosaur api.Dinosaur
	if err := g2.Take(&dinosaur, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &dinosaur, nil
}

func (d *sqlDinosaurDao) Create(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, error) {
	g2 := (*d.sessionFactory).New(ctx)
	if err := g2.Omit(clause.Associations).Create(dinosaur).Error; err != nil {
		db.MarkForRollback(ctx, err)
		return nil, err
	}
	return dinosaur, nil
}

func (d *sqlDinosaurDao) Replace(ctx context.Context, dinosaur *api.Dinosaur) (*api.Dinosaur, error) {
	g2 := (*d.sessionFactory).New(ctx)
	if err := g2.Omit(clause.Associations).Save(dinosaur).Error; err != nil {
		db.MarkForRollback(ctx, err)
		return nil, err
	}
	return dinosaur, nil
}

func (d *sqlDinosaurDao) Delete(ctx context.Context, id string) error {
	g2 := (*d.sessionFactory).New(ctx)
	if err := g2.Omit(clause.Associations).Delete(&api.Dinosaur{Meta: api.Meta{ID: id}}).Error; err != nil {
		db.MarkForRollback(ctx, err)
		return err
	}
	return nil
}

func (d *sqlDinosaurDao) FindByIDs(ctx context.Context, ids []string) (api.DinosaurList, error) {
	g2 := (*d.sessionFactory).New(ctx)
	dinosaurs := api.DinosaurList{}
	if err := g2.Where("id in (?)", ids).Find(&dinosaurs).Error; err != nil {
		return nil, err
	}
	return dinosaurs, nil
}

func (d *sqlDinosaurDao) FindBySpecies(ctx context.Context, species string) (api.DinosaurList, error) {
	g2 := (*d.sessionFactory).New(ctx)
	dinosaurs := api.DinosaurList{}
	if err := g2.Where("species = ?", species).Find(&dinosaurs).Error; err != nil {
		return nil, err
	}
	return dinosaurs, nil
}

func (d *sqlDinosaurDao) All(ctx context.Context) (api.DinosaurList, error) {
	g2 := (*d.sessionFactory).New(ctx)
	dinosaurs := api.DinosaurList{}
	if err := g2.Find(&dinosaurs).Error; err != nil {
		return nil, err
	}
	return dinosaurs, nil
}
