package mocks

import (
	"context"

	"gorm.io/gorm"

	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/errors"
)

var _ dao.SubscriptionDao = &subscriptionDaoMock{}

type subscriptionDaoMock struct {
	subscriptions api.SubscriptionList
}

func NewSubscriptionDao() *subscriptionDaoMock {
	return &subscriptionDaoMock{}
}

func (d *subscriptionDaoMock) Get(ctx context.Context, id string) (*api.Subscription, error) {
	for _, dino := range d.subscriptions {
		if dino.ID == id {
			return dino, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (d *subscriptionDaoMock) Create(ctx context.Context, subscription *api.Subscription) (*api.Subscription, error) {
	d.subscriptions = append(d.subscriptions, subscription)
	return subscription, nil
}

func (d *subscriptionDaoMock) Replace(ctx context.Context, subscription *api.Subscription) (*api.Subscription, error) {
	return nil, errors.NotImplemented("Subscription").AsError()
}

func (d *subscriptionDaoMock) Delete(ctx context.Context, id string) error {
	return errors.NotImplemented("Subscription").AsError()
}

func (d *subscriptionDaoMock) FindByIDs(ctx context.Context, ids []string) (api.SubscriptionList, error) {
	return nil, errors.NotImplemented("Subscription").AsError()
}

func (d *subscriptionDaoMock) All(ctx context.Context) (api.SubscriptionList, error) {
	return d.subscriptions, nil
}
