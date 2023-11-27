package test

import (
	"context"
	"fmt"

	"github.com/openshift-online/rh-trex/pkg/api"
)

func (helper *Helper) NewDinosaur(species string) *api.Dinosaur {
	dinoService := helper.Env().Services.Dinosaurs()

	dinosaur := &api.Dinosaur{
		Species: species,
	}

	dino, err := dinoService.Create(context.Background(), dinosaur)
	if err != nil {
		helper.T.Errorf("error creating dinosaur: %q", err)
	}

	return dino
}

func (helper *Helper) NewDinosaurList(namePrefix string, count int) (dinosaurs []*api.Dinosaur) {
	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("%s_%d", namePrefix, i)
		dinosaurs = append(dinosaurs, helper.NewDinosaur(name))
	}
	return dinosaurs
}

func (helper *Helper) NewSubscription(id string) *api.Subscription {
	subService := helper.Env().Services.Subscriptions()

	subscription := &api.Subscription{
		Meta: api.Meta{ID: id},
	}

	sub, err := subService.Create(context.Background(), subscription)
	if err != nil {
		helper.T.Errorf("error creating dinosaur: %q", err)
	}

	return sub
}

func (helper *Helper) NewSubscriptionList(id string, count int) (subscriptions []*api.Subscription) {
	for i := 1; i <= count; i++ {
		subscriptions = append(subscriptions, helper.NewSubscription(id))
	}
	return subscriptions
}
