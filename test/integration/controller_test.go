package integration

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/openshift-online/rh-trex/cmd/trex/server"
	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/controllers"
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/db"
	"github.com/openshift-online/rh-trex/test"
)

func TestControllerRacing(t *testing.T) {
	h, _ := test.RegisterIntegration(t)

	account := h.NewRandAccount()
	authCtx := h.NewAuthenticatedContext(account)
	dao := dao.NewEventDao(&h.Env().Database.SessionFactory)

	// The handler filters the events by source id/type/reconciled, and only record
	// the event with create type. Due to the event lock, each create event
	// should be only processed once.
	var proccessedEvent []string
	var mu sync.Mutex
	onUpsert := func(ctx context.Context, id string) error {
		events, err := dao.All(authCtx)
		if err != nil {
			return err
		}

		for _, evt := range events {
			if evt.SourceID != id {
				continue
			}
			if evt.EventType != api.CreateEventType {
				continue
			}
			// the event has been reconciled by others, ignore.
			if evt.ReconciledDate != nil {
				continue
			}
			mu.Lock()
			proccessedEvent = append(proccessedEvent, id)
			mu.Unlock()
		}

		return nil
	}

	// Start 3 controllers concurrently
	threads := 3
	for i := 0; i < threads; i++ {
		go func() {
			s := &server.ControllersServer{
				KindControllerManager: controllers.NewKindControllerManager(
					db.NewAdvisoryLockFactory(h.Env().Database.SessionFactory),
					h.Env().Services.Events(),
				),
			}

			s.KindControllerManager.Add(&controllers.ControllerConfig{
				Source: "Dinosaurs",
				Handlers: map[api.EventType][]controllers.ControllerHandlerFunc{
					api.CreateEventType: {onUpsert},
					api.UpdateEventType: {onUpsert},
				},
			})

			s.Start()
		}()
	}

	_, err := h.Factories.NewDinosaurList("bronto", 50)
	Expect(err).NotTo(HaveOccurred())

	// This is to check only two create events is processed. It waits for 5 seconds to ensure all events have been
	// processed by the controllers.
	Eventually(func() error {
		mu.Lock()
		eventCount := len(proccessedEvent)
		mu.Unlock()
		if eventCount != 50 {
			return fmt.Errorf("should have only 50 create events but got %d", eventCount)
		}
		return nil
	}, 5*time.Second, 1*time.Second).Should(Succeed())
}
