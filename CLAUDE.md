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
   - `pkg/handlers/fizzbuzz.go` - HTTP handlers
   - `pkg/services/fizzbuzz.go` - Business logic with event handlers
   - `pkg/dao/fizzbuzz.go` - Data access layer
   - `pkg/dao/mocks/fizzbuzz.go` - Mock for testing
   - `pkg/db/migrations/YYYYMMDDHHMM_add_fizzbuzzs.go` - Database migration
   - `test/integration/fizzbuzzs_test.go` - Integration tests
   - `test/factories/fizzbuzzs.go` - Test data factories
   - `openapi/openapi.fizzbuzzs.yaml` - OpenAPI specification
   - `cmd/trex/environments/locator_fizzbuzz.go` - Service locator

2. **Updated Files** (automatically modified):
   - `pkg/api/presenters/kind.go` - Adds Kind mapping
   - `pkg/api/presenters/path.go` - Adds snake_case path mapping  
   - `cmd/trex/environments/types.go` - Adds service to Services struct
   - `cmd/trex/environments/framework.go` - Adds service initialization
   - `cmd/trex/server/routes.go` - Registers API routes
   - `cmd/trex/server/controllers.go` - Adds event handler registration
   - `pkg/db/migrations/migration_structs.go` - Enables database migration
   - `openapi/openapi.yaml` - Adds API references

3. **Regenerated OpenAPI Client**:
   ```bash
   make generate  # Automatically regenerates client code
   ```

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

**No manual steps are required** - the generator handles everything automatically!

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