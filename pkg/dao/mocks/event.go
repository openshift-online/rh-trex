package mocks

import (
	"context"

	"gorm.io/gorm"

	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/api"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/dao"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/errors"
)

var _ dao.EventDao = &eventDaoMock{}

type eventDaoMock struct {
	events api.EventList
}

func NewEventDao() *eventDaoMock {
	return &eventDaoMock{}
}

func (d *eventDaoMock) Get(ctx context.Context, id string) (*api.Event, error) {
	for _, dino := range d.events {
		if dino.ID == id {
			return dino, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (d *eventDaoMock) Create(ctx context.Context, event *api.Event) (*api.Event, error) {
	d.events = append(d.events, event)
	return event, nil
}

func (d *eventDaoMock) Replace(ctx context.Context, event *api.Event) (*api.Event, error) {
	return nil, errors.NotImplemented("Event").AsError()
}

func (d *eventDaoMock) Delete(ctx context.Context, id string) error {
	newEvents := api.EventList{}
	for _, e := range d.events {
		if e.ID == id {
			// deleting this one
			// do not include in the new list
		} else {
			newEvents = append(newEvents, e)
		}
	}
	d.events = newEvents
	return nil
}

func (d *eventDaoMock) FindByIDs(ctx context.Context, ids []string) (api.EventList, error) {
	return nil, errors.NotImplemented("Event").AsError()
}

func (d *eventDaoMock) All(ctx context.Context) (api.EventList, error) {
	return d.events, nil
}
