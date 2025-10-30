package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/openshift-online/rh-trex/cmd/trex/server"
	"github.com/openshift-online/rh-trex/pkg/api"
	"github.com/openshift-online/rh-trex/pkg/controllers"
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/db"
	"github.com/openshift-online/rh-trex/plugins/events"
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
			proccessedEvent = append(proccessedEvent, id)
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
					events.Service(&h.Env().Services),
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

	// make some time for the controllers to start, or they will miss events
	time.Sleep(100 * time.Millisecond)

	const N = 50

	dinos, err := h.Factories.NewDinosaurList("bronto", N)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(dinos)).To(Equal(N))

	// This is to check only two create events is processed. It waits for 5 seconds to ensure all events have been
	// processed by the controllers.
	Eventually(func() error {
		if len(proccessedEvent) != N {
			return fmt.Errorf("should have %d create events but got %d", N, len(proccessedEvent))
		}
		return nil
	}, 10*time.Second, 1*time.Second).Should(Succeed())
}
