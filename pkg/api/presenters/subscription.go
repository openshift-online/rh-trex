package presenters

import (
	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/api/openapi"
	"github.com/openshift-online/rh-trex/pkg/util"
)

func ConvertSubscription(subscription openapi.Subscription) *api.Subscription {
	return &api.Subscription{
		Meta: api.Meta{
			ID: util.NilToEmptyString(subscription.Id),
		},
	}
}

func PresentSubscription(subscription *api.Subscription) openapi.Subscription {
	reference := PresentReference(subscription.ID, subscription)
	return openapi.Subscription{
		Id:   reference.Id,
		Kind: reference.Kind,
		Href: reference.Href,
	}
}
