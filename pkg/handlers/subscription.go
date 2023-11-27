package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/api/openapi"
	"github.com/openshift-online/rh-trex/pkg/api/presenters"
	"github.com/openshift-online/rh-trex/pkg/errors"
	"github.com/openshift-online/rh-trex/pkg/services"
)

var _ RestHandler = subscriptionHandler{}

type subscriptionHandler struct {
	subscription services.SubscriptionService
	generic      services.GenericService
}

func NewSubscriptionHandler(subscription services.SubscriptionService, generic services.GenericService) *subscriptionHandler {
	return &subscriptionHandler{
		subscription: subscription,
		generic:      generic,
	}
}

func (h subscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var subscription openapi.Subscription
	cfg := &handlerConfig{
		&subscription,
		[]validate{
			validateEmpty(&subscription, "Id", "id"),
		},
		func() (interface{}, *errors.ServiceError) {
			ctx := r.Context()
			dino := presenters.ConvertSubscription(subscription)
			dino, err := h.subscription.Create(ctx, dino)
			if err != nil {
				return nil, err
			}
			return presenters.PresentSubscription(dino), nil
		},
		handleError,
	}

	handle(w, r, cfg, http.StatusCreated)
}

func (h subscriptionHandler) Patch(w http.ResponseWriter, r *http.Request) {
	var patch openapi.SubscriptionPatchRequest

	cfg := &handlerConfig{
		&patch,
		[]validate{},
		func() (interface{}, *errors.ServiceError) {
			ctx := r.Context()
			id := mux.Vars(r)["id"]
			found, err := h.subscription.Get(ctx, id)
			if err != nil {
				return nil, err
			}

			//patch a field

			dino, err := h.subscription.Replace(ctx, found)
			if err != nil {
				return nil, err
			}
			return presenters.PresentSubscription(dino), nil
		},
		handleError,
	}

	handle(w, r, cfg, http.StatusOK)
}

func (h subscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	cfg := &handlerConfig{
		Action: func() (interface{}, *errors.ServiceError) {
			ctx := r.Context()

			listArgs := services.NewListArguments(r.URL.Query())
			var subscriptions = []api.Subscription{}
			paging, err := h.generic.List(ctx, "username", listArgs, &subscriptions)
			if err != nil {
				return nil, err
			}
			dinoList := openapi.SubscriptionList{
				Kind:  "SubscriptionList",
				Page:  int32(paging.Page),
				Size:  int32(paging.Size),
				Total: int32(paging.Total),
				Items: []openapi.Subscription{},
			}

			for _, dino := range subscriptions {
				converted := presenters.PresentSubscription(&dino)
				dinoList.Items = append(dinoList.Items, converted)
			}
			if listArgs.Fields != nil {
				filteredItems, err := presenters.SliceFilter(listArgs.Fields, dinoList.Items)
				if err != nil {
					return nil, err
				}
				return filteredItems, nil
			}
			return dinoList, nil
		},
	}

	handleList(w, r, cfg)
}

func (h subscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	cfg := &handlerConfig{
		Action: func() (interface{}, *errors.ServiceError) {
			id := mux.Vars(r)["id"]
			ctx := r.Context()
			subscription, err := h.subscription.Get(ctx, id)
			if err != nil {
				return nil, err
			}

			return presenters.PresentSubscription(subscription), nil
		},
	}

	handleGet(w, r, cfg)
}

func (h subscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	cfg := &handlerConfig{
		Action: func() (interface{}, *errors.ServiceError) {
			return nil, errors.NotImplemented("delete")
		},
	}
	handleDelete(w, r, cfg, http.StatusNoContent)
}
