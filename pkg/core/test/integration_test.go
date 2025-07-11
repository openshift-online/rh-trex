package test

import (
	"context"
	"testing"

	"github.com/openshift-online/rh-trex/pkg/api"
	coreapi "github.com/openshift-online/rh-trex/pkg/core/api"
	coredao "github.com/openshift-online/rh-trex/pkg/core/dao"
	coreservices "github.com/openshift-online/rh-trex/pkg/core/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestCoreDinosaurResource tests the core framework with a dinosaur resource
func TestCoreDinosaurResource(t *testing.T) {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&api.Dinosaur{}); err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	// Create core DAO
	dao := coredao.NewBaseDAO[api.Dinosaur](db)

	// Create mock event emitter
	eventEmitter := &MockEventEmitter{}

	// Create core service
	service := coreservices.NewBaseCRUDService[api.Dinosaur](dao, eventEmitter, "Dinosaurs")

	ctx := context.Background()

	// Test Create
	newDinosaur := &api.Dinosaur{
		Meta: coreapi.Meta{
			ID: "test-dino-1",
		},
		Species: "T-Rex",
	}

	created, serviceErr := service.Create(ctx, newDinosaur)
	if serviceErr != nil {
		t.Fatalf("Failed to create dinosaur: %v", serviceErr)
	}

	if created.Species != "T-Rex" {
		t.Errorf("Expected species 'T-Rex', got '%s'", created.Species)
	}

	// Test Get
	retrieved, serviceErr := service.Get(ctx, created.ID)
	if serviceErr != nil {
		t.Fatalf("Failed to get dinosaur: %v", serviceErr)
	}

	if retrieved.Species != "T-Rex" {
		t.Errorf("Expected species 'T-Rex', got '%s'", retrieved.Species)
	}

	// Test Update
	retrieved.Species = "Velociraptor"
	updated, serviceErr := service.Replace(ctx, retrieved)
	if serviceErr != nil {
		t.Fatalf("Failed to update dinosaur: %v", serviceErr)
	}

	if updated.Species != "Velociraptor" {
		t.Errorf("Expected species 'Velociraptor', got '%s'", updated.Species)
	}

	// Test List
	query := coreapi.ListQuery{
		Page: 1,
		Size: 10,
	}

	list, serviceErr := service.List(ctx, query)
	if serviceErr != nil {
		t.Fatalf("Failed to list dinosaurs: %v", serviceErr)
	}

	if list.Total != 1 {
		t.Errorf("Expected total 1, got %d", list.Total)
	}

	// Test Delete
	serviceErr = service.Delete(ctx, created.ID)
	if serviceErr != nil {
		t.Fatalf("Failed to delete dinosaur: %v", serviceErr)
	}

	// Verify deletion
	_, serviceErr = service.Get(ctx, created.ID)
	if serviceErr == nil {
		t.Error("Expected error when getting deleted dinosaur")
	}

	// Verify events were emitted
	if len(eventEmitter.Events) < 3 {
		t.Errorf("Expected at least 3 events, got %d", len(eventEmitter.Events))
	}

	// Check event types
	expectedEvents := []coreapi.EventType{
		coreapi.CreateEventType,
		coreapi.UpdateEventType,
		coreapi.DeleteEventType,
	}

	for i, expectedEvent := range expectedEvents {
		if i < len(eventEmitter.Events) {
			if eventEmitter.Events[i].EventType != expectedEvent {
				t.Errorf("Expected event %d to be %s, got %s", i, expectedEvent, eventEmitter.Events[i].EventType)
			}
		}
	}
}

// MockEventEmitter for testing
type MockEventEmitter struct {
	Events []MockEvent
}

type MockEvent struct {
	Source    string
	SourceID  string
	EventType coreapi.EventType
}

func (m *MockEventEmitter) EmitEvent(ctx context.Context, source, sourceID string, eventType coreapi.EventType) error {
	m.Events = append(m.Events, MockEvent{
		Source:    source,
		SourceID:  sourceID,
		EventType: eventType,
	})
	return nil
}