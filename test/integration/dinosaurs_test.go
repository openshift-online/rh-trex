package integration

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/services"

	. "github.com/onsi/gomega"
	"gopkg.in/resty.v1"

	"github.com/openshift-online/rh-trex/pkg/api/openapi"
	"github.com/openshift-online/rh-trex/test"
)

func TestDinosaurGet(t *testing.T) {
	h, client := test.RegisterIntegration(t)

	account := h.NewRandAccount()
	ctx := h.NewAuthenticatedContext(account)

	// 401 using no JWT token
	_, _, err := client.DefaultAPI.ApiRhTrexV1DinosaursIdGet(context.Background(), "foo").Execute()
	Expect(err).To(HaveOccurred(), "Expected 401 but got nil error")

	// GET responses per openapi spec: 200 and 404,
	_, resp, err := client.DefaultAPI.ApiRhTrexV1DinosaursIdGet(ctx, "foo").Execute()
	Expect(err).To(HaveOccurred(), "Expected 404")
	Expect(resp.StatusCode).To(Equal(http.StatusNotFound))

	dino, err := h.Factories.NewDinosaur(h.NewID())
	Expect(err).NotTo(HaveOccurred())

	dinosaur, resp, err := client.DefaultAPI.ApiRhTrexV1DinosaursIdGet(ctx, dino.ID).Execute()
	Expect(err).NotTo(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	Expect(*dinosaur.Id).To(Equal(dino.ID), "found object does not match test object")
	Expect(dinosaur.Species).To(Equal(dino.Species), "species mismatch")
	Expect(*dinosaur.Kind).To(Equal("Dinosaur"))
	Expect(*dinosaur.Href).To(Equal(fmt.Sprintf("/api/rh-trex/v1/dinosaurs/%s", dino.ID)))
	Expect(*dinosaur.CreatedAt).To(BeTemporally("~", dino.CreatedAt))
	Expect(*dinosaur.UpdatedAt).To(BeTemporally("~", dino.UpdatedAt))
}

func TestDinosaurPost(t *testing.T) {
	h, client := test.RegisterIntegration(t)

	account := h.NewRandAccount()
	ctx := h.NewAuthenticatedContext(account)

	// POST responses per openapi spec: 201, 409, 500
	dino := openapi.Dinosaur{
		Species: time.Now().String(),
	}

	// 201 Created
	dinosaur, resp, err := client.DefaultAPI.ApiRhTrexV1DinosaursPost(ctx).Dinosaur(dino).Execute()
	Expect(err).NotTo(HaveOccurred(), "Error posting object:  %v", err)
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))
	Expect(*dinosaur.Id).NotTo(BeEmpty(), "Expected ID assigned on creation")
	Expect(dinosaur.Species).To(Equal(dino.Species), "species mismatch")
	Expect(*dinosaur.Kind).To(Equal("Dinosaur"))
	Expect(*dinosaur.Href).To(Equal(fmt.Sprintf("/api/rh-trex/v1/dinosaurs/%s", *dinosaur.Id)))

	// 400 bad request. posting junk json is one way to trigger 400.
	jwtToken := ctx.Value(openapi.ContextAccessToken)
	restyResp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", jwtToken)).
		SetBody(`{ this is invalid }`).
		Post(h.RestURL("/dinosaurs"))

	Expect(err).NotTo(HaveOccurred(), "Error posting object:  %v", err)
	Expect(restyResp.StatusCode()).To(Equal(http.StatusBadRequest))
}

func TestDinosaurPatch(t *testing.T) {
	h, client := test.RegisterIntegration(t)

	account := h.NewRandAccount()
	ctx := h.NewAuthenticatedContext(account)

	// POST responses per openapi spec: 201, 409, 500

	dino, err := h.Factories.NewDinosaur("Brontosaurus")
	Expect(err).NotTo(HaveOccurred())

	// 200 OK
	species := "Dodo"
	dinosaur, resp, err := client.DefaultAPI.ApiRhTrexV1DinosaursIdPatch(ctx, dino.ID).DinosaurPatchRequest(openapi.DinosaurPatchRequest{Species: &species}).Execute()
	Expect(err).NotTo(HaveOccurred(), "Error posting object:  %v", err)
	Expect(resp.StatusCode).To(Equal(http.StatusOK))
	Expect(*dinosaur.Id).To(Equal(dino.ID))
	Expect(dinosaur.Species).To(Equal(species), "species mismatch")
	Expect(*dinosaur.CreatedAt).To(BeTemporally("~", dino.CreatedAt))
	Expect(*dinosaur.Kind).To(Equal("Dinosaur"))
	Expect(*dinosaur.Href).To(Equal(fmt.Sprintf("/api/rh-trex/v1/dinosaurs/%s", *dinosaur.Id)))

	jwtToken := ctx.Value(openapi.ContextAccessToken)
	// 500 server error. posting junk json is one way to trigger 500.
	restyResp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", jwtToken)).
		SetBody(`{ this is invalid }`).
		Patch(h.RestURL("/dinosaurs/foo"))

	Expect(err).NotTo(HaveOccurred(), "Error posting object:  %v", err)
	Expect(restyResp.StatusCode()).To(Equal(http.StatusBadRequest))

	// species can not be empty in request body
	restyResp, err = resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", jwtToken)).
		SetBody(`{"species":""}`).
		Patch(h.RestURL(fmt.Sprintf("/dinosaurs/%s", *dinosaur.Id)))
	Expect(err).NotTo(HaveOccurred(), "Error posting object:  %v", err)
	Expect(restyResp.StatusCode()).To(Equal(http.StatusBadRequest))
	Expect(restyResp.String()).To(ContainSubstring("species cannot be empty"))

	Eventually(func() error {
		dao := dao.NewEventDao(&h.Env().Database.SessionFactory)
		events, err := dao.FindByIDs(ctx, []string{*dinosaur.Id})
		Expect(err).NotTo(HaveOccurred(), "Error getting events:  %v", err)
		Expect(len(events)).To(Equal(2), "expected Create and Update events")
		Expect(contains(api.CreateEventType, events)).To(BeTrue())
		Expect(contains(api.UpdateEventType, events)).To(BeTrue())
		return nil
	}, 5*time.Second, 1*time.Second)
}

func contains(et api.EventType, events api.EventList) bool {
	for _, e := range events {
		if e.EventType == et {
			return true
		}
	}
	return false
}

func TestDinosaurPaging(t *testing.T) {
	h, client := test.RegisterIntegration(t)

	account := h.NewRandAccount()
	ctx := h.NewAuthenticatedContext(account)

	// Paging
	_, err := h.Factories.NewDinosaurList("Bronto", 20)
	Expect(err).NotTo(HaveOccurred())

	list, _, err := client.DefaultAPI.ApiRhTrexV1DinosaursGet(ctx).Execute()
	Expect(err).NotTo(HaveOccurred(), "Error getting dinosaur list: %v", err)
	Expect(len(list.Items)).To(Equal(20))
	Expect(list.Size).To(Equal(int32(20)))
	Expect(list.Total).To(Equal(int32(20)))
	Expect(list.Page).To(Equal(int32(1)))

	list, _, err = client.DefaultAPI.ApiRhTrexV1DinosaursGet(ctx).Page(2).Size(5).Execute()
	Expect(err).NotTo(HaveOccurred(), "Error getting dinosaur list: %v", err)
	Expect(len(list.Items)).To(Equal(5))
	Expect(list.Size).To(Equal(int32(5)))
	Expect(list.Total).To(Equal(int32(20)))
	Expect(list.Page).To(Equal(int32(2)))
}

func TestDinosaurListSearch(t *testing.T) {
	h, client := test.RegisterIntegration(t)

	account := h.NewRandAccount()
	ctx := h.NewAuthenticatedContext(account)

	dinosaurs, err := h.Factories.NewDinosaurList("bronto", 20)
	Expect(err).NotTo(HaveOccurred())

	search := fmt.Sprintf("id in ('%s')", dinosaurs[0].ID)
	list, _, err := client.DefaultAPI.ApiRhTrexV1DinosaursGet(ctx).Search(search).Execute()
	Expect(err).NotTo(HaveOccurred(), "Error getting dinosaur list: %v", err)
	Expect(len(list.Items)).To(Equal(1))
	Expect(list.Total).To(Equal(int32(1)))
	Expect(*list.Items[0].Id).To(Equal(dinosaurs[0].ID))
}

func TestUpdateDinosaurWithRacingRequests_BlockingAdvisoryLock(t *testing.T) {
	testUpdateDinosaurWithRacingRequests(t, true, true, 2)
}

func TestUpdateDinosaurWithRacingRequests_NonBlockingAdvisoryLock(t *testing.T) {
	testUpdateDinosaurWithRacingRequests(t, true, false, 1)
}

func TestUpdateDinosaurWithRacingRequests_WithoutLock(t *testing.T) {
	testUpdateDinosaurWithRacingRequests(t, false, false, 2)
}

func testUpdateDinosaurWithRacingRequests(t *testing.T, useAdvisoryLock, useBlockingAdvisoryLock bool, expectedUpdates int) {
	h, client := test.RegisterIntegration(t)

	services.DisableAdvisoryLock = !useAdvisoryLock
	services.UseBlockingAdvisoryLock = useBlockingAdvisoryLock

	defer func() {
		services.DisableAdvisoryLock = false
		services.UseBlockingAdvisoryLock = true
	}()

	account := h.NewRandAccount()
	ctx := h.NewAuthenticatedContext(account)

	dino, err := h.Factories.NewDinosaur("Stegosaurus")
	Expect(err).NotTo(HaveOccurred())

	firstDinoUpdate := "AdvisoryLockosaurus"
	secondDinoUpdate := "AdvisoryLockosaurusSecond"
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup
	wg1.Add(1)
	wg2.Add(2)

	// the call to update the dino with name==AdvisoryLockosaurus take 1second in the service, since it has hardcoded an special case
	go func() {
		wg1.Done()
		client.DefaultAPI.ApiRhTrexV1DinosaursIdPatch(ctx, dino.ID).DinosaurPatchRequest(openapi.DinosaurPatchRequest{Species: &firstDinoUpdate}).Execute()
		Expect(err).NotTo(HaveOccurred(), "Error posting object:  %v", err)
		wg2.Done()
	}()

	// we need the gorutine to have sent the first update
	// there may be some uncertainty, since the first API call should have been fired, but may or may not have reached the server yet
	// so we wait additional time
	wg1.Wait()
	time.Sleep(100 * time.Millisecond)

	go func() {
		client.DefaultAPI.ApiRhTrexV1DinosaursIdPatch(ctx, dino.ID).DinosaurPatchRequest(openapi.DinosaurPatchRequest{Species: &secondDinoUpdate}).Execute()
		wg2.Done()
	}()

	// waits for all goroutines above to complete
	wg2.Wait()

	eventdao := dao.NewEventDao(&h.Env().Database.SessionFactory)
	events, err := eventdao.All(ctx)
	Expect(err).NotTo(HaveOccurred(), "Error getting events:  %v", err)

	dinodao := dao.NewDinosaurDao(&h.Env().Database.SessionFactory)
	readDino, err := dinodao.Get(ctx, dino.ID)
	if useBlockingAdvisoryLock {
		Expect(readDino.Species).To(Equal(secondDinoUpdate))
	} else {
		Expect(readDino.Species).To(Equal(firstDinoUpdate))
	}

	updatedCount := 0
	for _, e := range events {
		if e.SourceID == dino.ID && e.EventType == api.UpdateEventType {
			updatedCount = updatedCount + 1
		}
	}

	// the dinosaur patch request is protected by the advisory lock, so there should only be one update
	Expect(updatedCount).To(Equal(expectedUpdates))

	// all the locks should be released finally
	Eventually(func() error {
		var count int
		err := h.DBFactory.DirectDB().
			QueryRow("select count(*) from pg_locks where locktype='advisory';").
			Scan(&count)
		Expect(err).NotTo(HaveOccurred(), "Error querying pg_locks:  %v", err)

		if count != 0 {
			return fmt.Errorf("there are %d unreleased advisory lock", count)
		}
		return nil
	}, 5*time.Second, 1*time.Second).Should(Succeed())
}
