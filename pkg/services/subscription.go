package services

import (
	"context"
	"github.com/openshift-online/rh-trex/pkg/dao"

	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/errors"
)

type SubscriptionService interface {
	Get(ctx context.Context, id string) (*api.Subscription, *errors.ServiceError)
	Create(ctx context.Context, subscription *api.Subscription) (*api.Subscription, *errors.ServiceError)
	Replace(ctx context.Context, subscription *api.Subscription) (*api.Subscription, *errors.ServiceError)
	Delete(ctx context.Context, id string) *errors.ServiceError
	All(ctx context.Context) (api.SubscriptionList, *errors.ServiceError)

	FindByIDs(ctx context.Context, ids []string) (api.SubscriptionList, *errors.ServiceError)
}

func NewSubscriptionService(subscriptionDao dao.SubscriptionDao) SubscriptionService {
	return &sqlSubscriptionService{
		subscriptionDao: subscriptionDao,
	}
}

var _ SubscriptionService = &sqlSubscriptionService{}

type sqlSubscriptionService struct {
	subscriptionDao dao.SubscriptionDao
}

func (s *sqlSubscriptionService) Get(ctx context.Context, id string) (*api.Subscription, *errors.ServiceError) {
	subscription, err := s.subscriptionDao.Get(ctx, id)
	if err != nil {
		return nil, handleGetError("Subscription", "id", id, err)
	}
	return subscription, nil
}

func (s *sqlSubscriptionService) Create(ctx context.Context, subscription *api.Subscription) (*api.Subscription, *errors.ServiceError) {
	subscription, err := s.subscriptionDao.Create(ctx, subscription)
	if err != nil {
		return nil, handleCreateError("Subscription", err)
	}
	return subscription, nil
}

func (s *sqlSubscriptionService) Replace(ctx context.Context, subscription *api.Subscription) (*api.Subscription, *errors.ServiceError) {
	subscription, err := s.subscriptionDao.Replace(ctx, subscription)
	if err != nil {
		return nil, handleUpdateError("Subscription", err)
	}
	return subscription, nil
}

func (s *sqlSubscriptionService) Delete(ctx context.Context, id string) *errors.ServiceError {
	if err := s.subscriptionDao.Delete(ctx, id); err != nil {
		return handleDeleteError("Subscription", errors.GeneralError("Unable to delete subscription: %s", err))
	}
	return nil
}

func (s *sqlSubscriptionService) FindByIDs(ctx context.Context, ids []string) (api.SubscriptionList, *errors.ServiceError) {
	subscriptions, err := s.subscriptionDao.FindByIDs(ctx, ids)
	if err != nil {
		return nil, errors.GeneralError("Unable to get all subscriptions: %s", err)
	}
	return subscriptions, nil
}

func (s *sqlSubscriptionService) All(ctx context.Context) (api.SubscriptionList, *errors.ServiceError) {
	subscriptions, err := s.subscriptionDao.All(ctx)
	if err != nil {
		return nil, errors.GeneralError("Unable to get all subscriptions: %s", err)
	}
	return subscriptions, nil
}
