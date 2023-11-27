package integration

import (
	"context"
	"fmt"
	"github.com/openshift-online/rh-trex/pkg/api"
	"net/http"
	"testing"

	. "github.com/onsi/gomega"
	"gopkg.in/resty.v1"

	"github.com/openshift-online/rh-trex/pkg/api/openapi"
	"github.com/openshift-online/rh-trex/test"
)

func TestSubscriptionGet(t *testing.T) {
	h, client := test.RegisterIntegration(t)

	account := h.NewRandAccount()
	ctx := h.NewAuthenticatedContext(account)

	// 401 using no JWT token
	_, _, err := client.DefaultApi.ApiRhTrexV1SubscriptionsIdGet(context.Background(), "foo").Execute()
	Expect(err).To(HaveOccurred(), "Expected 401 but got nil error")

	// GET responses per openapi spec: 200 and 404,
	_, resp, err := client.DefaultApi.ApiRhTrexV1SubscriptionsIdGet(ctx, "foo").Execute()
	Expect(err).To(HaveOccurred(), "Expected 404")
	Expect(resp.StatusCode).To(Equal(http.StatusNotFound))

	dino := h.NewSubscription(api.NewID())

	subscription, resp, err := client.DefaultApi.ApiRhTrexV1SubscriptionsIdGet(ctx, dino.ID).Execute()
	Expect(err).NotTo(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	Expect(*subscription.Id).To(Equal(dino.ID), "found object does not match test object")
	Expect(*subscription.Kind).To(Equal("Subscription"))
	Expect(*subscription.Href).To(Equal(fmt.Sprintf("/api/rh-trex/v1/subscriptions/%s", dino.ID)))
}

func TestSubscriptionPost(t *testing.T) {
	h, client := test.RegisterIntegration(t)

	account := h.NewRandAccount()
	ctx := h.NewAuthenticatedContext(account)

	// POST responses per openapi spec: 201, 409, 500
	dino := openapi.Subscription{}

	// 201 Created
	subscription, resp, err := client.DefaultApi.ApiRhTrexV1SubscriptionsPost(ctx).Subscription(dino).Execute()
	Expect(err).NotTo(HaveOccurred(), "Error posting object:  %v", err)
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))
	Expect(*subscription.Id).NotTo(BeEmpty(), "Expected ID assigned on creation")
	Expect(*subscription.Kind).To(Equal("Subscription"))
	Expect(*subscription.Href).To(Equal(fmt.Sprintf("/api/rh-trex/v1/subscriptions/%s", *subscription.Id)))

	// 400 bad request. posting junk json is one way to trigger 400.
	jwtToken := ctx.Value(openapi.ContextAccessToken)
	restyResp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", jwtToken)).
		SetBody(`{ this is invalid }`).
		Post(h.RestURL("/subscriptions"))

	Expect(restyResp.StatusCode()).To(Equal(http.StatusBadRequest))
}

func TestSubscriptionPatch(t *testing.T) {
	h, client := test.RegisterIntegration(t)

	account := h.NewRandAccount()
	ctx := h.NewAuthenticatedContext(account)

	// POST responses per openapi spec: 201, 409, 500

	dino := h.NewSubscription(api.NewID())

	// 200 OK
	subscription, resp, err := client.DefaultApi.ApiRhTrexV1SubscriptionsIdPatch(ctx, dino.ID).SubscriptionPatchRequest(openapi.SubscriptionPatchRequest{}).Execute()
	Expect(err).NotTo(HaveOccurred(), "Error posting object:  %v", err)
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	Expect(*subscription.Id).To(Equal(dino.ID))
	Expect(*subscription.Kind).To(Equal("Subscription"))
	Expect(*subscription.Href).To(Equal(fmt.Sprintf("/api/rh-trex/v1/subscriptions/%s", *subscription.Id)))

	jwtToken := ctx.Value(openapi.ContextAccessToken)
	// 500 server error. posting junk json is one way to trigger 500.
	restyResp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", jwtToken)).
		SetBody(`{ this is invalid }`).
		Patch(h.RestURL("/subscriptions/foo"))

	Expect(restyResp.StatusCode()).To(Equal(http.StatusBadRequest))
}

func TestSubscriptionPaging(t *testing.T) {
	h, client := test.RegisterIntegration(t)

	account := h.NewRandAccount()
	ctx := h.NewAuthenticatedContext(account)

	// Paging
	_ = h.NewSubscriptionList(api.NewID(), 20)

	list, _, err := client.DefaultApi.ApiRhTrexV1SubscriptionsGet(ctx).Execute()
	Expect(err).NotTo(HaveOccurred(), "Error getting subscription list: %v", err)
	Expect(len(list.Items)).To(Equal(20))
	Expect(list.Size).To(Equal(int32(20)))
	Expect(list.Total).To(Equal(int32(20)))
	Expect(list.Page).To(Equal(int32(1)))

	list, _, err = client.DefaultApi.ApiRhTrexV1SubscriptionsGet(ctx).Page(2).Size(5).Execute()
	Expect(err).NotTo(HaveOccurred(), "Error getting subscription list: %v", err)
	Expect(len(list.Items)).To(Equal(5))
	Expect(list.Size).To(Equal(int32(5)))
	Expect(list.Total).To(Equal(int32(20)))
	Expect(list.Page).To(Equal(int32(2)))
}

func TestSubscriptionListSearch(t *testing.T) {
	h, client := test.RegisterIntegration(t)

	account := h.NewRandAccount()
	ctx := h.NewAuthenticatedContext(account)

	subscriptions := h.NewSubscriptionList(api.NewID(), 20)

	search := fmt.Sprintf("id in ('%s')", subscriptions[0].ID)
	list, _, err := client.DefaultApi.ApiRhTrexV1SubscriptionsGet(ctx).Search(search).Execute()
	Expect(err).NotTo(HaveOccurred(), "Error getting subscription list: %v", err)
	Expect(len(list.Items)).To(Equal(1))
	Expect(list.Total).To(Equal(int32(20)))
	Expect(*list.Items[0].Id).To(Equal(subscriptions[0].ID))
}
