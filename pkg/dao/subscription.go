package dao

import (
	"context"

	"gorm.io/gorm/clause"

	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/db"
)

type SubscriptionDao interface {
	Get(ctx context.Context, id string) (*api.Subscription, error)
	Create(ctx context.Context, subscription *api.Subscription) (*api.Subscription, error)
	Replace(ctx context.Context, subscription *api.Subscription) (*api.Subscription, error)
	Delete(ctx context.Context, id string) error
	FindByIDs(ctx context.Context, ids []string) (api.SubscriptionList, error)
	All(ctx context.Context) (api.SubscriptionList, error)
}

var _ SubscriptionDao = &sqlSubscriptionDao{}

type sqlSubscriptionDao struct {
	sessionFactory *db.SessionFactory
}

func NewSubscriptionDao(sessionFactory *db.SessionFactory) SubscriptionDao {
	return &sqlSubscriptionDao{sessionFactory: sessionFactory}
}

func (d *sqlSubscriptionDao) Get(ctx context.Context, id string) (*api.Subscription, error) {
	g2 := (*d.sessionFactory).New(ctx)
	var subscription api.Subscription
	if err := g2.Take(&subscription, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (d *sqlSubscriptionDao) Create(ctx context.Context, subscription *api.Subscription) (*api.Subscription, error) {
	g2 := (*d.sessionFactory).New(ctx)
	if err := g2.Omit(clause.Associations).Create(subscription).Error; err != nil {
		db.MarkForRollback(ctx, err)
		return nil, err
	}
	return subscription, nil
}

func (d *sqlSubscriptionDao) Replace(ctx context.Context, subscription *api.Subscription) (*api.Subscription, error) {
	g2 := (*d.sessionFactory).New(ctx)
	if err := g2.Omit(clause.Associations).Save(subscription).Error; err != nil {
		db.MarkForRollback(ctx, err)
		return nil, err
	}
	return subscription, nil
}

func (d *sqlSubscriptionDao) Delete(ctx context.Context, id string) error {
	g2 := (*d.sessionFactory).New(ctx)
	if err := g2.Omit(clause.Associations).Delete(&api.Subscription{Meta: api.Meta{ID: id}}).Error; err != nil {
		db.MarkForRollback(ctx, err)
		return err
	}
	return nil
}

func (d *sqlSubscriptionDao) FindByIDs(ctx context.Context, ids []string) (api.SubscriptionList, error) {
	g2 := (*d.sessionFactory).New(ctx)
	subscriptions := api.SubscriptionList{}
	if err := g2.Where("id in (?)", ids).Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (d *sqlSubscriptionDao) All(ctx context.Context) (api.SubscriptionList, error) {
	g2 := (*d.sessionFactory).New(ctx)
	subscriptions := api.SubscriptionList{}
	if err := g2.Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}
