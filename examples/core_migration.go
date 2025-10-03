package examples

import (
	"context"

	"github.com/openshift-online/rh-trex/pkg/api"
	corecontrollers "github.com/openshift-online/rh-trex-core/controllers"
	"github.com/openshift-online/rh-trex-core/generator"
	"github.com/openshift-online/rh-trex/pkg/dao"
	"github.com/openshift-online/rh-trex/pkg/db"
	"github.com/openshift-online/rh-trex/pkg/services"
	"gorm.io/gorm"
)

// This example shows how to migrate existing TRex code to use the core framework

func ExampleCoreMigration() {
	// Assume we have a database connection
	var gormDB *gorm.DB
	var sessionFactory *db.SessionFactory
	
	// Step 1: Create core-based DAO
	dinosaurDAO := dao.NewDinosaurCoreDAO(sessionFactory)
	
	// Step 2: Create core-based service
	eventService := services.NewEventService(nil, nil, nil) // existing event service
	lockFactory := db.NewAdvisoryLockFactory(*sessionFactory)
	
	dinosaurService := services.NewDinosaurCoreService(
		dinosaurDAO,
		lockFactory,
		eventService,
	)
	
	// Step 3: Use the core framework for automatic registration
	eventEmitter := services.NewEventEmitterAdapter(eventService)
	
	// Create resource factory for future resources
	resourceFactory := generator.NewResourceFactory(
		gormDB,
		nil, // controller manager (would be injected)
		eventEmitter,
	)
	
	// Step 4: Register resources automatically (for new resources)
	// This would automatically create DAO, service, and controller registration
	_ = generator.RegisterResourceType(resourceFactory, "Dinosaurs", api.Dinosaur{})
	
	// Step 5: Use the service normally
	ctx := context.Background()
	
	// Create a new dinosaur
	newDinosaur := &api.Dinosaur{
		Species: "T-Rex",
	}
	
	created, err := dinosaurService.Create(ctx, newDinosaur)
	if err != nil {
		// Handle error
	}
	
	// Get the dinosaur
	retrieved, err := dinosaurService.Get(ctx, created.ID)
	if err != nil {
		// Handle error
	}
	
	// Use dinosaur-specific methods
	trexes, err := dinosaurService.FindBySpecies(ctx, "T-Rex")
	if err != nil {
		// Handle error
	}
	
	// All of these operations automatically:
	// - Use the core DAO patterns
	// - Emit events through the core event system
	// - Get processed by the core controller framework
	// - Follow consistent error handling patterns
	
	_ = retrieved
	_ = trexes
}

// ExampleNewResourceWithCoreFramework shows how to create a new resource using the core framework
func ExampleNewResourceWithCoreFramework() {
	// Define a new resource type
	type Planet struct {
		api.Meta
		Name     string `json:"name" gorm:"index"`
		Type     string `json:"type" gorm:"index"`
		HasLife  bool   `json:"has_life"`
		Distance float64 `json:"distance"`
	}
	
	// Implement the MetaProvider interface
	func (p *Planet) GetMeta() *api.Meta {
		return &p.Meta
	}
	
	// With the core framework, this is all you need to do:
	var gormDB *gorm.DB
	var eventEmitter services.EventEmitter
	var controllerManager *corecontrollers.ControllerManager
	
	// Create resource factory
	factory := generator.NewResourceFactory(gormDB, controllerManager, eventEmitter)
	
	// Register the resource - this automatically creates:
	// - DAO with all CRUD operations
	// - Service with all CRUD operations and event handling
	// - Controller registration for event processing
	planetService := generator.RegisterResourceType(factory, "Planets", Planet{})
	
	// Use the service immediately
	ctx := context.Background()
	
	newPlanet := &Planet{
		Name:     "Earth",
		Type:     "Terrestrial",
		HasLife:  true,
		Distance: 1.0,
	}
	
	created, err := planetService.Create(ctx, newPlanet)
	if err != nil {
		// Handle error
	}
	
	// All CRUD operations work automatically
	planets, err := planetService.FindByIDs(ctx, []string{created.ID})
	if err != nil {
		// Handle error
	}
	
	// Events are automatically emitted and processed
	// No manual controller registration needed
	// No manual DAO implementation needed
	// No manual service implementation needed
	
	_ = planets
}

// ExampleGradualMigration shows how to gradually migrate from legacy to core framework
func ExampleGradualMigration() {
	// You can run both systems side by side during migration
	
	// Legacy system (existing)
	legacyDinosaurService := services.NewDinosaurService(
		nil, // legacy dependencies
		nil,
		nil,
	)
	
	// Core system (new)
	coreDinosaurService := services.NewDinosaurCoreService(
		nil, // core dependencies
		nil,
		nil,
	)
	
	// Use whichever system you prefer
	ctx := context.Background()
	
	// Legacy way
	legacyDinosaur, err := legacyDinosaurService.Get(ctx, "some-id")
	if err != nil {
		// Handle error
	}
	
	// Core way (same result, but benefits from core framework)
	coreDinosaur, err := coreDinosaurService.Get(ctx, "some-id")
	if err != nil {
		// Handle error
	}
	
	// Both work the same way from the API perspective
	// But the core version automatically benefits from:
	// - Consistent error handling
	// - Automatic event emission
	// - Generic patterns
	// - Future framework improvements
	
	_ = legacyDinosaur
	_ = coreDinosaur
}