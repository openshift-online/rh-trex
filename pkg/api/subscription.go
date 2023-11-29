package api

import "gorm.io/gorm"

type Subscription struct {
	Meta
}

type SubscriptionList []*Subscription
type SubscriptionIndex map[string]*Subscription

func (l SubscriptionList) Index() SubscriptionIndex {
	index := SubscriptionIndex{}
	for _, o := range l {
		index[o.ID] = o
	}
	return index
}

func (d *Subscription) BeforeCreate(tx *gorm.DB) error {
	d.ID = NewID()
	return nil
}

type SubscriptionPatchRequest struct {
}
