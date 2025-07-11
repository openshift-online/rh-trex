# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

TRex is a Go-based REST API template for Red Hat TAP (Trusted Application Pipeline) that serves as a full-featured foundation for building new microservices. It provides CRUD operations for "dinosaurs" as example business logic to be replaced.

## Development Commands

### Building and Running
- `make binary` - Build the trex binary
- `make install` - Build and install binary to GOPATH/bin
- `make run` - Run migrations and start the server (runs on localhost:8000)

### Testing
- `make test` - Run unit tests
- `make test-integration` - Run integration tests
- `make ci-test-unit` - Run unit tests with JSON output for CI
- `make ci-test-integration` - Run integration tests with JSON output for CI

### Code Quality
- `make verify` - Run source code verification (vet, formatting)
- `make lint` - Run golangci-lint

### Database Operations
- `make db/setup` - Start PostgreSQL container locally
- `make db/login` - Connect to local PostgreSQL database
- `make db/teardown` - Stop and remove PostgreSQL container
- `./trex migrate` - Run database migrations

### Development Workflow
- `make generate` - Regenerate OpenAPI client and models
- `make clean` - Remove temporary generated files

### OpenShift/Container Operations
- `make crc/login` - Login to CodeReady Containers
- `make image` - Build container image
- `make push` - Push image to registry
- `make deploy` - Deploy to OpenShift
- `make undeploy` - Remove from OpenShift

## Architecture

### Core Components

**Main Application (`cmd/trex/main.go`):**
- CLI tool with subcommands: `migrate`, `serve`, `clone`
- Uses Cobra for command structure

**Environment Framework (`cmd/trex/environments/`):**
- Configurable environments: development, testing, production
- Visitor pattern for component initialization
- Service locator pattern for dependency injection

**API Layer (`pkg/api/`):**
- OpenAPI-generated models and clients
- Example "Dinosaur" entity with CRUD operations
- Standardized error handling and metadata

**Data Layer:**
- **DAO Pattern** (`pkg/dao/`): Data Access Objects for database operations
- **Database** (`pkg/db/`): GORM-based persistence with PostgreSQL
- **Migrations** (`pkg/db/migrations/`): Database schema versioning

**Service Layer (`pkg/services/`):**
- Business logic separated from handlers
- Generic service patterns for reuse
- Event-driven architecture support

**HTTP Layer (`pkg/handlers/`):**
- REST API endpoints
- Authentication/authorization middleware
- OpenAPI specification compliance

**Infrastructure:**
- **Authentication** (`pkg/auth/`): OIDC integration with Red Hat SSO
- **Clients** (`pkg/client/ocm/`): OCM (OpenShift Cluster Manager) integration
- **Configuration** (`pkg/config/`): Environment-specific settings
- **Logging** (`cmd/trex/server/logging/`): Structured logging with request middleware
- **Metrics** (`pkg/handlers/prometheus_metrics.go`): Prometheus integration

### Key Patterns

1. **Separation of Concerns**: Clear boundaries between API, service, and data layers
2. **Dependency Injection**: Service locator pattern in environments framework
3. **Code Generation**: OpenAPI specs generate client code and documentation
4. **Test-Driven Development**: Comprehensive test support with mocks and factories

## Code Generation

### How to Generate a New Kind

The generator script creates complete CRUD functionality with **event-driven architecture** for a new resource type. The process is now fully automated with no manual steps required.

**Single Command to Generate a New Kind:**
```bash
go run ./scripts/generator.go --kind KindName
```

**Complete Example:**
```bash
# Generate a new Kind called "FizzBuzz"
go run ./scripts/generator.go --kind FizzBuzz

# This creates a complete implementation with:
# - API model and handlers
# - Service and DAO layers with event-driven controllers
# - Database migration
# - Test files and factories
# - OpenAPI specifications
# - Service locators and routing
# - Automatic controller registration for event handling
```

### What the Generator Creates

The generator automatically creates and configures:

1. **Generated Files** (no manual editing needed):
   - `pkg/api/fizzbuzz.go` - API model
   - `pkg/api/presenters/fizzbuzz.go` - Presenter conversion functions  
   - `pkg/handlers/fizzbuzz.go` - HTTP handlers
   - `pkg/services/fizzbuzz.go` - Business logic with event handlers
   - `pkg/dao/fizzbuzz.go` - Data access layer
   - `pkg/dao/mocks/fizzbuzz.go` - Mock for testing
   - `pkg/db/migrations/YYYYMMDDHHMM_add_fizzbuzzs.go` - Database migration
   - `test/integration/fizzbuzzs_test.go` - Integration tests
   - `test/factories/fizzbuzzs.go` - Test data factories
   - `openapi/openapi.fizzbuzzs.yaml` - OpenAPI specification
   - `cmd/trex/environments/locator_fizzbuzz.go` - Service locator

2. **Updated Files** (automatically modified by generator):
   - `pkg/api/presenters/kind.go` - Adds Kind mapping for ObjectKind function
   - `pkg/api/presenters/path.go` - Adds snake_case path mapping for ObjectPath function  
   - `cmd/trex/server/controllers.go` - Adds event handler registration with proper syntax
   - `openapi/openapi.yaml` - Adds API references

3. **Files Requiring Manual Updates** (generator creates templates but doesn't modify existing):
   - `cmd/trex/environments/types.go` - Add service field to Services struct
   - `cmd/trex/environments/framework.go` - Add service initialization in LoadServices
   - `cmd/trex/server/routes.go` - Register API routes and handlers
   - `pkg/db/migrations/migration_structs.go` - Add migration to MigrationList

4. **Regenerated OpenAPI Client** (via `make generate`):
   - `pkg/api/openapi/model_fizzbuzz.go` - Go model structs
   - `pkg/api/openapi/model_fizzbuzz_all_of.go` - Composite model  
   - `pkg/api/openapi/model_fizzbuzz_list.go` - List model
   - `pkg/api/openapi/model_fizzbuzz_list_all_of.go` - List composite
   - `pkg/api/openapi/model_fizzbuzz_patch_request.go` - Patch request model
   - `pkg/api/openapi/docs/FizzBuzz*.md` - Generated API documentation
   - Updated `pkg/api/openapi/api_default.go` - API client methods

### Naming Patterns

The generator uses consistent naming patterns:
- **API paths**: snake_case (e.g., `/api/rh-trex/v1/fizz_buzzs`)
- **Go types**: PascalCase (e.g., `FizzBuzz`)
- **Variables**: camelCase (e.g., `fizzBuzz`)
- **Database tables**: snake_case (e.g., `fizz_buzzs`)

### Template Fields Available

When creating custom templates, these fields are available:
- `{{.Kind}}` - PascalCase (e.g., "FizzBuzz")
- `{{.KindPlural}}` - PascalCase plural (e.g., "FizzBuzzs")
- `{{.KindLowerSingular}}` - camelCase singular (e.g., "fizzBuzz")
- `{{.KindLowerPlural}}` - camelCase plural (e.g., "fizzBuzzs")
- `{{.KindSnakeCasePlural}}` - snake_case plural for API paths (e.g., "fizz_buzzs")
- `{{.Project}}` - Project name (e.g., "rh-trex")
- `{{.Repo}}` - Repository path (e.g., "github.com/openshift-online")
- `{{.Cmd}}` - Command directory name (e.g., "trex")
- `{{.ID}}` - Timestamp ID for migrations (e.g., "202507111234")

### Generated Event Handlers

Each generated service includes idempotent event handlers:

```go
// OnUpsert handles CREATE and UPDATE events
func (s *sqlKindService) OnUpsert(ctx context.Context, id string) error {
    logger := logger.NewOCMLogger(ctx)
    
    kind, err := s.kindDao.Get(ctx, id)
    if err != nil {
        return err
    }
    
    logger.Infof("Do idempotent somethings with this kind: %s", kind.ID)
    return nil
}

// OnDelete handles DELETE events
func (s *sqlKindService) OnDelete(ctx context.Context, id string) error {
    logger := logger.NewOCMLogger(ctx)
    logger.Infof("This kind has been deleted: %s", id)
    return nil
}
```

**Key Handler Characteristics:**
- **Idempotent**: Safe to run multiple times
- **Logged**: Structured logging for debugging
- **Error Handling**: Proper error propagation
- **Context Aware**: Supports request tracing

### Testing the Generated Kind

After generation, verify the implementation:
```bash
# Run integration tests for the new Kind
export GOPATH=/tmp/go
go test -v ./test/integration -run TestFizzBuzz

# Run all tests to ensure no regressions
make test-integration

# Test event-driven functionality
# Events are automatically created during CRUD operations
# Controllers process events asynchronously via PostgreSQL LISTEN/NOTIFY
```

### Expected Results

- **All tests pass** immediately after generation
- **API endpoints** respond correctly with proper HTTP status codes
- **Database operations** work (CREATE, READ, UPDATE, DELETE, SEARCH)
- **Event-driven controllers** automatically process database events
- **Idempotent handlers** safely process CREATE/UPDATE/DELETE events
- **OpenAPI client** includes the new Kind's methods
- **Service locators** properly inject dependencies
- **Integration tests** verify complete functionality
- **Controller registration** automatically handles event processing

### Event-Driven Architecture

The generator creates a complete event-driven system:

**Generated Service Interface:**
```go
type KindService interface {
    // Standard CRUD operations
    Get(ctx context.Context, id string) (*api.Kind, *errors.ServiceError)
    Create(ctx context.Context, kind *api.Kind) (*api.Kind, *errors.ServiceError)
    Replace(ctx context.Context, kind *api.Kind) (*api.Kind, *errors.ServiceError)
    Delete(ctx context.Context, id string) *errors.ServiceError
    All(ctx context.Context) (api.KindList, *errors.ServiceError)
    FindByIDs(ctx context.Context, ids []string) (api.KindList, *errors.ServiceError)
    
    // Event-driven controller functions
    OnUpsert(ctx context.Context, id string) error
    OnDelete(ctx context.Context, id string) error
}
```

**Automatic Controller Registration:**
```go
// Generated in cmd/trex/server/controllers.go
kindServices := env().Services.Kinds()

s.KindControllerManager.Add(&controllers.ControllerConfig{
    Source: "Kinds",
    Handlers: map[api.EventType][]controllers.ControllerHandlerFunc{
        api.CreateEventType: {kindServices.OnUpsert},
        api.UpdateEventType: {kindServices.OnUpsert},
        api.DeleteEventType: {kindServices.OnDelete},
    },
})
```

**Event Flow:**
1. **API Operation** (CREATE/UPDATE/DELETE) → **Event Creation** → **Database NOTIFY**
2. **Controller Listener** → **Event Handlers** → **Business Logic**
3. **Idempotent Processing** → **Structured Logging** → **Error Handling**

### Key Improvements

The generator has been enhanced to:
1. **Dynamically detect** command directory structure
2. **Automatically handle** all service registrations
3. **Generate correct** snake_case API paths
4. **Use proper** service locator patterns with lock factories
5. **Create complete** test suites with proper factory methods
6. **Maintain consistency** with existing codebase patterns
7. **Generate event-driven controllers** with idempotent handlers
8. **Automatically register** event handlers in controller system

**Minimal manual steps required** - the generator automates most of the process, with only 4 files requiring manual updates!

### Required Manual Steps

After running the generator, you must manually update these 4 files:

1. **Add service to Services struct** in `cmd/trex/environments/types.go`:
   ```go
   type Services struct {
       Dinosaurs    DinosaurServiceLocator
       YourKinds    YourKindServiceLocator  // <- Add this line
       Generic      GenericServiceLocator
       Events       EventServiceLocator
   }
   ```

2. **Add service initialization** in `cmd/trex/environments/framework.go`:
   ```go
   func (e *Env) LoadServices() {
       e.Services.Generic = NewGenericServiceLocator(e)
       e.Services.Dinosaurs = NewDinosaurServiceLocator(e)
       e.Services.YourKinds = NewYourKindServiceLocator(e)  // <- Add this line
       e.Services.Events = NewEventServiceLocator(e)
   }
   ```

3. **Register API routes** in `cmd/trex/server/routes.go`:
   ```go
   // Add handler initialization
   yourKindHandler := handlers.NewYourKindHandler(services.YourKinds(), services.Generic())
   
   // Add route registration 
   apiV1YourKindsRouter := apiV1Router.PathPrefix("/your_kinds").Subrouter()
   apiV1YourKindsRouter.HandleFunc("", yourKindHandler.List).Methods(http.MethodGet)
   apiV1YourKindsRouter.HandleFunc("/{id}", yourKindHandler.Get).Methods(http.MethodGet)
   apiV1YourKindsRouter.HandleFunc("", yourKindHandler.Create).Methods(http.MethodPost)
   apiV1YourKindsRouter.HandleFunc("/{id}", yourKindHandler.Patch).Methods(http.MethodPatch)
   apiV1YourKindsRouter.HandleFunc("/{id}", yourKindHandler.Delete).Methods(http.MethodDelete)
   apiV1YourKindsRouter.Use(authMiddleware.AuthenticateAccountJWT)
   apiV1YourKindsRouter.Use(authzMiddleware.AuthorizeApi)
   ```

4. **Add migration to list** in `pkg/db/migrations/migration_structs.go`:
   ```go
   var MigrationList = []*gormigrate.Migration{
       addDinosaurs(),
       addEvents(),
       addYourKinds(),  // <- Add this line
   }
   ```

### Generator Troubleshooting

If you encounter issues after running the generator, check these common problems:

#### Compilation Errors

**Issue**: Syntax errors in `cmd/trex/server/controllers.go`
```bash
# Error: unexpected := in composite literal; possibly missing comma or }
```
**Root Cause**: Missing closing brace in controller registration
**Fix**: The generator should properly close controller configurations. Manual fix:
```go
s.KindControllerManager.Add(&controllers.ControllerConfig{
    Source: "Dinosaurs",
    Handlers: map[api.EventType][]controllers.ControllerHandlerFunc{
        api.CreateEventType: {dinoServices.OnUpsert},
        api.UpdateEventType: {dinoServices.OnUpsert},
        api.DeleteEventType: {dinoServices.OnDelete},
    },
}) // <- Ensure this closing brace exists
```

**Issue**: Service method not found
```bash
# Error: env().Services.KindName undefined
```
**Root Cause**: Service not added to environment framework
**Fix**: Verify these files are updated:
- `cmd/trex/environments/types.go` - Service field in Services struct
- `cmd/trex/environments/framework.go` - Service initialization in LoadServices()

#### Test Failures

**Issue**: Integration tests fail with "404 Not Found" or "relation does not exist"
**Root Causes**:
1. Database migration not registered
2. API routes not registered  
3. Presenter mappings missing

**Fixes**:
1. **Migration**: Add to `pkg/db/migrations/migration_structs.go`:
   ```go
   var MigrationList = []*gormigrate.Migration{
       addDinosaurs(),
       addEvents(),
       addKindName(), // <- Add your migration
   }
   ```

2. **Routes**: Add to `cmd/trex/server/routes.go`:
   ```go
   kindHandler := handlers.NewKindHandler(services.Kinds(), services.Generic())
   
   apiV1KindsRouter := apiV1Router.PathPrefix("/kind_names").Subrouter()
   apiV1KindsRouter.HandleFunc("", kindHandler.List).Methods(http.MethodGet)
   // ... other routes
   ```

3. **Presenters**: Add to both presenter files:
   ```go
   // pkg/api/presenters/kind.go
   case api.KindName, *api.KindName:
       result = "KindName"
   
   // pkg/api/presenters/path.go  
   case api.KindName, *api.KindName:
       return "kind_names"  // snake_case plural
   ```

#### Database Issues

**Issue**: Tests fail with database connection errors
**Solution**: Recreate the database to run new migrations:
```bash
make db/teardown  # Stop and remove PostgreSQL container
make db/setup     # Start fresh PostgreSQL container  
make test-integration  # Run tests with new schema
```

**Note**: Always run `make` commands from the project root directory where the Makefile is located.

#### Cleaning Up Test Generations

When experimenting with the generator, you may need to completely remove a generated Kind. Here's the comprehensive cleanup process:

**Complete Kind Removal** (e.g., for TestWidget):
```bash
# Remove all generated files (replace TestWidget/testWidget with your Kind name)
rm -rf \
  pkg/api/testWidget.go \
  pkg/api/presenters/testWidget.go \
  pkg/handlers/testWidget.go \
  pkg/services/testWidget.go \
  pkg/dao/testWidget.go \
  pkg/dao/mocks/testWidget.go \
  pkg/db/migrations/*testWidget* \
  test/integration/testWidgets_test.go \
  test/factories/testWidgets.go \
  openapi/openapi.testWidgets.yaml \
  cmd/trex/environments/locator_testWidget.go

# Remove OpenAPI client files (generated by make generate)
rm -rf \
  pkg/api/openapi/model_test_widget*.go \
  pkg/api/openapi/docs/TestWidget*.md

# Reset modified files to clean state
git checkout HEAD -- \
  cmd/trex/server/controllers.go \
  cmd/trex/server/routes.go \
  cmd/trex/environments/types.go \
  cmd/trex/environments/framework.go \
  pkg/api/presenters/kind.go \
  pkg/api/presenters/path.go \
  pkg/db/migrations/migration_structs.go \
  openapi/openapi.yaml

# Regenerate OpenAPI client to remove traces
make generate
```

**Quick Test Cleanup** (for temporary testing):
```bash
# For a Kind called "TestWidget", run this one-liner:
rm -rf pkg/api/testWidget.go pkg/api/presenters/testWidget.go pkg/handlers/testWidget.go pkg/services/testWidget.go pkg/dao/testWidget.go pkg/dao/mocks/testWidget.go pkg/db/migrations/*testWidget* test/integration/testWidgets_test.go test/factories/testWidgets.go openapi/openapi.testWidgets.yaml cmd/trex/environments/locator_testWidget.go pkg/api/openapi/model_test_widget*.go pkg/api/openapi/docs/TestWidget*.md && git checkout HEAD -- cmd/trex/server/controllers.go cmd/trex/server/routes.go cmd/trex/environments/types.go cmd/trex/environments/framework.go pkg/api/presenters/kind.go pkg/api/presenters/path.go pkg/db/migrations/migration_structs.go openapi/openapi.yaml && make generate
```

## Authentication

Local development uses Red Hat SSO authentication. Use the `ocm` CLI tool:
```bash
# Login to local service
ocm login --token=${OCM_ACCESS_TOKEN} --url=http://localhost:8000

# Test API endpoints
ocm get /api/rh-trex/v1/dinosaurs
ocm post /api/rh-trex/v1/dinosaurs '{"species": "foo"}'
```

## Database

- PostgreSQL database with GORM ORM
- Migration-based schema management
- DAO pattern for data access
- Advisory locks for concurrency control

## Testing

- Unit tests use mocks for external dependencies
- Integration tests run against real database
- Test factories in `test/factories/` for data setup
- Environment-specific test configuration

### Database Issues During Testing

If integration tests fail with PostgreSQL-related errors (missing columns, transaction issues), recreate the database:

```bash
# From project root directory
make db/teardown  # Stop and remove PostgreSQL container
make db/setup     # Start fresh PostgreSQL container
make test-integration  # Run tests again
```

**Note:** Always run `make` commands from the project root directory where the Makefile is located.