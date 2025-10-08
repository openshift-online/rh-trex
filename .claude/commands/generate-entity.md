# Generate Entity Command

Generate a complete CRUD entity in the TRex application following the **plugin-based architecture**.

## Plugin System Overview

**NEW APPROACH**: TRex now uses a plugin system for entity registration. Each entity is self-contained in a plugin package that automatically registers:
- Service locators
- HTTP routes
- Event controllers
- Presenter mappings (Kind and Path)

**Key Benefits**:
- ✅ **90% reduction in manual updates** (from 8 files to 3 files)
- ✅ **Single source of truth** - all entity wiring in one plugin file
- ✅ **Auto-discovery** - plugin registers itself via `init()` function
- ✅ **Type-safe** - compile-time checks for service access
- ✅ **Easier maintenance** - self-contained entity logic

**File Count Summary**:
- **1 new plugin file** (`plugins/{entity}s/plugin.go`) - replaces 5 manual update steps
- **10 standard files** (API model, DAO, service, handlers, presenters, migration, OpenAPI, tests, factories)
- **3 manual updates** (main.go import, migration list, OpenAPI refs)

## Instructions

You will guide the user through creating a new entity with all required artifacts. Use the Dinosaur plugin (`plugins/dinosaurs/plugin.go`) as the reference pattern.

### Step 1: Gather Requirements

Ask the user for:
1. **Entity Name** (singular, PascalCase): e.g., "Widget", "Product", "Customer"
2. **Entity Fields**: Additional fields beyond the base Meta fields (ID, CreatedAt, UpdatedAt, DeletedAt)
   - Field name (camelCase in code, snake_case in DB)
   - Field type (string, int, bool, time.Time, etc.)
   - Database constraints (index, unique, etc.)
3. **API Path** (plural, lowercase): e.g., "widgets", "products", "customers"

### Step 2: Create Required Files

Use the TodoWrite tool to track your progress through these steps:

#### 2.1 Plugin Package (`plugins/{entity}s/plugin.go`)

**This is the NEW core file that replaces manual service locator, route registration, and controller setup.**

Create a plugin file that registers:
- Service locator function
- Route registration
- Controller registration
- Presenter mappings (Kind and Path)

Example pattern from `plugins/dinosaurs/plugin.go`:
```go
package widgets

import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/openshift-online/rh-trex/cmd/trex/environments"
    "github.com/openshift-online/rh-trex/cmd/trex/environments/registry"
    "github.com/openshift-online/rh-trex/cmd/trex/server"
    "github.com/openshift-online/rh-trex/pkg/api"
    "github.com/openshift-online/rh-trex/pkg/api/presenters"
    "github.com/openshift-online/rh-trex/pkg/auth"
    "github.com/openshift-online/rh-trex/pkg/controllers"
    "github.com/openshift-online/rh-trex/pkg/dao"
    "github.com/openshift-online/rh-trex/pkg/db"
    "github.com/openshift-online/rh-trex/pkg/handlers"
    "github.com/openshift-online/rh-trex/pkg/services"
)

// Service Locator
type WidgetServiceLocator func() services.WidgetService

func NewWidgetServiceLocator(env *environments.Env) WidgetServiceLocator {
    return func() services.WidgetService {
        return services.NewWidgetService(
            db.NewAdvisoryLockFactory(env.Database.SessionFactory),
            dao.NewWidgetDao(&env.Database.SessionFactory),
            env.Services.Events(),
        )
    }
}

// WidgetService helper function to get the widget service from the registry
func WidgetService(s *environments.Services) services.WidgetService {
    if s == nil {
        return nil
    }
    if obj := s.GetService("Widgets"); obj != nil {
        locator := obj.(WidgetServiceLocator)
        return locator()
    }
    return nil
}

func init() {
    // Service registration
    registry.RegisterService("Widgets", func(env interface{}) interface{} {
        return NewWidgetServiceLocator(env.(*environments.Env))
    })

    // Routes registration
    server.RegisterRoutes("widgets", func(apiV1Router *mux.Router, services server.ServicesInterface, authMiddleware auth.JWTMiddleware, authzMiddleware auth.AuthorizationMiddleware) {
        envServices := services.(*environments.Services)
        widgetHandler := handlers.NewWidgetHandler(WidgetService(envServices), envServices.Generic())

        widgetsRouter := apiV1Router.PathPrefix("/widgets").Subrouter()
        widgetsRouter.HandleFunc("", widgetHandler.List).Methods(http.MethodGet)
        widgetsRouter.HandleFunc("/{id}", widgetHandler.Get).Methods(http.MethodGet)
        widgetsRouter.HandleFunc("", widgetHandler.Create).Methods(http.MethodPost)
        widgetsRouter.HandleFunc("/{id}", widgetHandler.Patch).Methods(http.MethodPatch)
        widgetsRouter.HandleFunc("/{id}", widgetHandler.Delete).Methods(http.MethodDelete)
        widgetsRouter.Use(authMiddleware.AuthenticateAccountJWT)
        widgetsRouter.Use(authzMiddleware.AuthorizeApi)
    })

    // Controller registration
    server.RegisterController("Widgets", func(manager *controllers.KindControllerManager, services *environments.Services) {
        widgetServices := WidgetService(services)

        manager.Add(&controllers.ControllerConfig{
            Source: "Widgets",
            Handlers: map[api.EventType][]controllers.ControllerHandlerFunc{
                api.CreateEventType: {widgetServices.OnUpsert},
                api.UpdateEventType: {widgetServices.OnUpsert},
                api.DeleteEventType: {widgetServices.OnDelete},
            },
        })
    })

    // Presenter registration
    presenters.RegisterPath(api.Widget{}, "widgets")
    presenters.RegisterPath(&api.Widget{}, "widgets")
    presenters.RegisterKind(api.Widget{}, "Widget")
    presenters.RegisterKind(&api.Widget{}, "Widget")
}
```

**Key Features:**
- **Self-contained**: All entity registration in one file
- **Auto-discovery**: Uses `init()` function for automatic registration
- **No manual edits needed**: Eliminates updates to routes.go, controllers.go, types.go, framework.go
- **Service locator pattern**: Provides type-safe service access
- **Helper function**: `WidgetService()` retrieves service from registry

#### 2.2 API Model (`pkg/api/{entity}_types.go`)

Create the entity struct with:
- Embedded `Meta` struct (provides ID, CreatedAt, UpdatedAt, DeletedAt)
- Custom fields from requirements
- List and Index types
- BeforeCreate hook for ID generation
- PatchRequest struct for updates

Example pattern:
```go
package api

import "gorm.io/gorm"

type Widget struct {
    Meta
    Name        string
    Description string
    Status      string
}

type WidgetList []*Widget
type WidgetIndex map[string]*Widget

func (l WidgetList) Index() WidgetIndex {
    index := WidgetIndex{}
    for _, o := range l {
        index[o.ID] = o
    }
    return index
}

func (w *Widget) BeforeCreate(tx *gorm.DB) error {
    w.ID = NewID()
    return nil
}

type WidgetPatchRequest struct {
    Name        *string `json:"name,omitempty"`
    Description *string `json:"description,omitempty"`
    Status      *string `json:"status,omitempty"`
}
```

#### 2.2 DAO Layer (`pkg/dao/{entity}.go`)

Create interface and implementation with:
- Get(ctx, id) - retrieve by ID
- Create(ctx, entity) - create new record
- Replace(ctx, entity) - update existing record
- Delete(ctx, id) - delete record
- FindByIDs(ctx, ids) - batch retrieval
- All(ctx) - retrieve all records
- Custom finders as needed

Pattern: See `pkg/dao/dinosaur.go`

#### 2.3 DAO Mock (`pkg/dao/mocks/{entity}.go`)

Create mock implementation for testing.

Pattern: See `pkg/dao/mocks/dinosaur.go`

#### 2.4 Service Layer (`pkg/services/{entity}s.go`)

Create interface and implementation with:
- CRUD operations (Get, Create, Replace, Delete, All, FindByIDs)
- Event-driven handlers (OnUpsert, OnDelete)
- Business logic and validation
- Event creation for CREATE, UPDATE, DELETE operations
- Advisory lock for concurrent updates

Pattern: See `pkg/services/dinosaurs.go`

Key features:
- Use `LockFactory` for advisory locks on updates
- Create events after successful operations
- Implement idempotent OnUpsert and OnDelete handlers

**NOTE**: Service instantiation is now handled in the plugin file, NOT in a separate locator file

#### 2.5 Presenters (`pkg/api/presenters/{entity}.go`)

Create conversion functions:
- `Convert{Entity}(openapi.{Entity}) *api.{Entity}` - OpenAPI to internal model
- `Present{Entity}(*api.{Entity}) openapi.{Entity}` - Internal to OpenAPI model

Pattern: See `pkg/api/presenters/dinosaur.go`

#### 2.6 Handlers (`pkg/handlers/{entity}.go`)

Create HTTP handlers:
- Create - POST endpoint
- Get - GET by ID endpoint
- List - GET collection endpoint with pagination
- Patch - PATCH update endpoint
- Delete - DELETE endpoint

Pattern: See `pkg/handlers/dinosaur.go`

Include validation in handlers using the `validate` pattern.

#### 2.7 Database Migration (`pkg/db/migrations/{timestamp}_add_{entity}s.go`)

Create migration with:
- Timestamp ID (YYYYMMDDHHMM format)
- Inline model definition (never import from pkg/api)
- Migrate function using AutoMigrate
- Rollback function using DropTable

Pattern: See `pkg/db/migrations/201911212019_add_dinosaurs.go`

**IMPORTANT**: Use inline struct definition in migration, not imported types.

#### 2.8 OpenAPI Specification (`openapi/openapi.{entity}s.yaml`)

Create OpenAPI spec with:
- Path definitions for collection and item endpoints
- Schema definitions for entity, list, and patch request
- Parameter definitions (id, page, size, search, orderBy, fields)
- Security requirements (Bearer token)
- Response codes and schemas

Pattern: See `openapi/openapi.dinosaurs.yaml`

#### 2.9 Plugin Import (`cmd/trex/main.go`)

**IMPORTANT**: Add a blank import to ensure the plugin's `init()` function runs:

```go
import (
    _ "github.com/openshift-online/rh-trex/plugins/dinosaurs"
    _ "github.com/openshift-online/rh-trex/plugins/widgets"  // <- Add this
)
```

This triggers the plugin registration when the application starts.

#### 2.10 Test Factory (`test/factories/{entity}s.go`)

Create factory methods:
- `New{Entity}(params) (*api.{Entity}, error)` - create single entity
- `New{Entity}List(prefix, count) ([]*api.{Entity}, error)` - create list

Pattern: See `test/factories/dinosaurs.go`

#### 2.11 Integration Tests (`test/integration/{entity}s_test.go`)

Create integration tests:
- Test{Entity}Get - test GET by ID (200, 404, 401)
- Test{Entity}Post - test CREATE (201, 400)
- Test{Entity}Patch - test UPDATE (200, 404, 400)
- Test{Entity}Paging - test pagination
- Test{Entity}ListSearch - test search functionality

Pattern: See `test/integration/dinosaurs_test.go`

### Step 3: Update Existing Files

**With the plugin system, most manual file updates are ELIMINATED!** Only these files need updates:

#### 3.1 Update `cmd/trex/main.go`

Add plugin import (triggers auto-registration):
```go
import (
    _ "github.com/openshift-online/rh-trex/plugins/widgets"  // <- Add this
)
```

#### 3.2 Update `pkg/db/migrations/migration_structs.go`

Add migration to list:
```go
var MigrationList = []*gormigrate.Migration{
    addDinosaurs(),
    addEvents(),
    addWidgets(),  // <- Add this
}
```

#### 3.3 Update `openapi/openapi.yaml`

Add reference to entity spec:
```yaml
paths:
  $ref:
    - 'openapi.dinosaurs.yaml'
    - 'openapi.widgets.yaml'  # <- Add this
```

### Step 4: Generate OpenAPI Client Code

After creating the OpenAPI spec, run:
```bash
make generate
```

**IMPORTANT: Wait for completion and verify results**

This command takes 2-3 minutes to complete. You MUST:

1. **Run the command and wait for completion:**
   ```bash
   make generate 2>&1 | tee generate.log
   ```

2. **Verify the generated files exist:**
   ```bash
   ls -la pkg/api/openapi/model_{entity}*.go
   ls -la pkg/api/openapi/docs/{Entity}*.md
   ```

3. **Check for compilation errors:**
   ```bash
   go build ./cmd/trex
   ```

4. **If generation fails or times out:**
   - Check the Docker/Podman daemon is running
   - Review the full output in generate.log
   - Verify openapi.yaml syntax is valid

**Expected generated files:**
- `pkg/api/openapi/model_{entity}.go`
- `pkg/api/openapi/model_{entity}_all_of.go`
- `pkg/api/openapi/model_{entity}_list.go`
- `pkg/api/openapi/model_{entity}_list_all_of.go`
- `pkg/api/openapi/model_{entity}_patch_request.go`
- `pkg/api/openapi/docs/{Entity}*.md`
- Updated `pkg/api/openapi/api_default.go` with new endpoints

**Do NOT proceed to Step 5 until:**
- [ ] All expected files exist
- [ ] `go build ./cmd/trex` completes without errors
- [ ] Integration test files compile successfully

### Step 5: Verify Compilation

Before running tests, ensure the code compiles:

```bash
# Build the binary to catch any compilation errors
go build -o /tmp/trex ./cmd/trex

# If this fails, review the errors and fix:
# - Missing imports
# - Type mismatches in presenters
# - Undefined constants or types
```

### Step 6: Test the Implementation

```bash
# Run database migrations
make db/teardown
make db/setup
./trex migrate

# Run integration tests
make test-integration

# Run specific entity tests
go test -v ./test/integration -run TestWidget

```
#### Step 6.1: Create a test script

Create a shell script to test the new endpoints for Widget using curl commands

### Key Patterns to Follow

1. **Naming Conventions**:
   - API paths: snake_case plural (e.g., `/widgets`, `/product_categories`)
   - Go types: PascalCase (e.g., `Widget`, `ProductCategory`)
   - Go variables: camelCase (e.g., `widget`, `productCategory`)
   - Database tables: snake_case plural (e.g., `widgets`, `product_categories`)

2. **Event-Driven Architecture**:
   - Create events for all CREATE, UPDATE, DELETE operations
   - Implement idempotent OnUpsert and OnDelete handlers
   - Register handlers in controllers.go

3. **Database Patterns**:
   - Use advisory locks for concurrent updates
   - All entities embed the `Meta` struct
   - Migrations use inline struct definitions
   - Use GORM for ORM operations

4. **API Patterns**:
   - Use OpenAPI specs for all endpoints
   - Follow RESTful conventions
   - Use presenters to convert between OpenAPI and internal models
   - Include proper validation in handlers

5. **Testing Patterns**:
   - Create factory methods for test data
   - Test all CRUD operations
   - Test error cases (404, 400, 401)
   - Test pagination and search

6. **Security**:
   - All endpoints require JWT authentication
   - Use authorization middleware
   - Validate all inputs in handlers

### Checklist

After completing all steps, verify:
- [ ] **Plugin file created** (`plugins/{entity}s/plugin.go`)
- [ ] All 10 other new files created (API model, DAO, service, handlers, presenters, migration, OpenAPI, tests, factories)
- [ ] **Only 3 existing files updated** (main.go, migration_structs.go, openapi.yaml)
- [ ] Plugin import added to `cmd/trex/main.go`
- [ ] OpenAPI client code regenerated (`make generate`)
- [ ] Generated files verified to exist
- [ ] Code compiles without errors (`go build ./cmd/trex`)
- [ ] Database migrations run successfully
- [ ] Integration tests pass
- [ ] API endpoints respond correctly
- [ ] Events are created and processed
- [ ] **Plugin auto-discovery working** (routes, controllers, presenters registered automatically)

## Example Usage

To generate a new "Widget" entity with fields "name" and "description":

1. User provides entity details
2. **You create the plugin file first** (`plugins/widgets/plugin.go`) - this is the NEW core integration point
3. You create all other required files following the patterns
4. **You update only 3 existing files** (main.go import, migration_structs.go, openapi.yaml)
5. You run `make generate` to create OpenAPI client code
6. You verify with tests
7. Create a script file to test the endpoints


Remember to use TodoWrite tool to track progress through all steps!
