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
- `make test` - Run unit tests (ALWAYS run after any code changes)
- `make test-integration` - Run integration tests (run after major changes - slower)
- `make ci-test-unit` - Run unit tests with JSON output for CI
- `make ci-test-integration` - Run integration tests with JSON output for CI

**Testing Guidelines:**
- **ALWAYS run `make test` after any code changes** to ensure nothing breaks
- **Run `make test-integration` after major changes** (new features, refactoring, etc.) as it's slower but more comprehensive

### Code Quality
- `make verify` - Run source code verification (vet, formatting)
- `make lint` - Run golangci-lint

### Database Operations
- `make db/setup` - Start PostgreSQL container locally
- `make db/login` - Connect to local PostgreSQL database
- `make db/teardown` - Stop and remove PostgreSQL container
- `./trex migrate` - Run database migrations

### TRex CLI Commands

The `trex` binary provides three main subcommands for different operational tasks:

#### `trex serve` - Start the API Server
Serves the rh-trex REST API with full authentication, database connectivity, and monitoring capabilities.

**Basic Usage:**
```bash
./trex serve                              # Start server on localhost:8000
./trex serve --api-server-bindaddress :8080  # Custom bind address
```

**Key Configuration Options:**
- **Server Binding:**
  - `--api-server-bindaddress` - API server bind address (default: "localhost:8000")
  - `--api-server-hostname` - Server's public hostname
  - `--enable-https` - Enable HTTPS rather than HTTP
  - `--https-cert-file` / `--https-key-file` - TLS certificate files

- **Database Configuration:**
  - `--db-host-file` - Database host file (default: "secrets/db.host")
  - `--db-name-file` - Database name file (default: "secrets/db.name") 
  - `--db-user-file` - Database username file (default: "secrets/db.user")
  - `--db-password-file` - Database password file (default: "secrets/db.password")
  - `--db-port-file` - Database port file (default: "secrets/db.port")
  - `--db-sslmode` - Database SSL mode: disable | require | verify-ca | verify-full (default: "disable")
  - `--db-max-open-connections` - Maximum open DB connections (default: 50)
  - `--enable-db-debug` - Enable database debug mode

- **Authentication & Authorization:**
  - `--enable-jwt` - Enable JWT authentication validation (default: true)
  - `--enable-authz` - Enable authorization on endpoints (default: true)
  - `--jwk-cert-url` - JWK Certificate URL for JWT validation (default: Red Hat SSO)
  - `--jwk-cert-file` - Local JWK Certificate file
  - `--acl-file` - Access control list file

- **OCM Integration:**
  - `--enable-ocm-mock` - Enable mock OCM clients (default: true)
  - `--ocm-base-url` - OCM API base URL (default: integration environment)
  - `--ocm-token-url` - OCM token endpoint URL (default: Red Hat SSO)
  - `--ocm-client-id-file` - OCM API client ID file (default: "secrets/ocm-service.clientId")
  - `--ocm-client-secret-file` - OCM API client secret file (default: "secrets/ocm-service.clientSecret")
  - `--self-token-file` - OCM API privileged offline SSO token file
  - `--ocm-debug` - Enable OCM API debug logging

- **Monitoring & Health Checks:**
  - `--health-check-server-bindaddress` - Health check server address (default: "localhost:8083")
  - `--enable-health-check-https` - Enable HTTPS for health check server
  - `--metrics-server-bindaddress` - Metrics server address (default: "localhost:8080")
  - `--enable-metrics-https` - Enable HTTPS for metrics server

- **Error Monitoring:**
  - `--enable-sentry` - Enable Sentry error monitoring
  - `--enable-sentry-debug` - Enable Sentry debug mode
  - `--sentry-url` - Sentry instance base URL (default: "glitchtip.devshift.net")
  - `--sentry-key-file` - Sentry key file (default: "secrets/sentry.key")
  - `--sentry-project` - Sentry project ID (default: "53")
  - `--sentry-timeout` - Sentry request timeout (default: 5s)

- **Performance Tuning:**
  - `--http-read-timeout` - HTTP server read timeout (default: 5s)
  - `--http-write-timeout` - HTTP server write timeout (default: 30s)
  - `--label-metrics-inclusion-duration` - Telemetry collection timeframe (default: 168h)

#### `trex migrate` - Run Database Migrations
Executes database schema migrations to set up or update the database structure.

**Basic Usage:**
```bash
./trex migrate                           # Run all pending migrations
./trex migrate --enable-db-debug        # Run with database debug logging
```

**Configuration Options:**
- **Database Connection:** (same as serve command)
  - `--db-host-file`, `--db-name-file`, `--db-user-file`, `--db-password-file`
  - `--db-port-file`, `--db-sslmode`, `--db-rootcert`
  - `--db-max-open-connections` - Maximum DB connections (default: 50)
  - `--enable-db-debug` - Enable database debug mode

**Migration Process:**
- Applies all pending migrations in order
- Creates migration tracking table if needed
- Idempotent - safe to run multiple times
- Logs each migration applied

#### `trex clone` - Clone New TRex Instance
Creates a new microservice project based on the TRex template, replacing template content with new service details.

**Basic Usage:**
```bash
./trex clone --name my-service                           # Clone with custom name
./trex clone --name my-service --destination ./my-proj   # Custom destination
./trex clone --repo github.com/myorg --name my-service   # Custom git repo
```

**Configuration Options:**
- `--name` - Name of the new service (default: "rh-trex")
- `--destination` - Target directory for new instance (default: "/tmp/clone-test")
- `--repo` - Git repository path (default: "github.com/openshift-online")

**Clone Process:**
- Creates new directory structure
- Replaces template strings throughout codebase
- Updates Go module paths and imports
- Renames files and directories as needed
- Maintains Git history and structure

#### Common Global Flags
All subcommands support these logging flags:
- `--logtostderr` - Log to stderr instead of files (default: true)
- `--alsologtostderr` - Log to both stderr and files
- `--log_dir` - Directory for log files
- `--stderrthreshold` - Minimum log level for stderr (default: 2)
- `-v, --v` - Log level for verbose logs
- `--vmodule` - Module-specific log levels
- `--log_backtrace_at` - Emit stack trace at specific file:line

**Example Production Server Startup:**
```bash
./trex serve \
  --api-server-bindaddress ":8000" \
  --enable-https \
  --https-cert-file /etc/certs/tls.crt \
  --https-key-file /etc/certs/tls.key \
  --db-sslmode verify-full \
  --enable-sentry \
  --ocm-base-url https://api.openshift.com \
  --disable-ocm-mock
```

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

TRex now uses the **rh-trex-core** library for enhanced error handling, database operations, and controller management. This provides consistent patterns across all TRex-based microservices.

### Core Library Integration

**rh-trex-core Dependency:**
- **Repository**: `github.com/openshift-online/rh-trex-core`
- **Purpose**: Shared framework components for TRex-based microservices
- **Components**: Enhanced error handling, database utilities, controller management
- **Benefits**: Consistent patterns, reduced duplication, centralized improvements

**Enhanced Error Handling (`pkg/services/util.go`):**
- **PII Sanitization**: Automatically redacts sensitive fields in error messages
- **Constraint Detection**: Intelligent database constraint violation handling
- **Core Library Delegation**: TRex error handlers delegate to `rh-trex-core/errors`
- **Type Conversion**: Seamless conversion between core and TRex error types

```go
// Enhanced error handlers that delegate to core library
func handleGetError(resourceType, field string, value interface{}, err error) *errors.ServiceError {
    valueStr := fmt.Sprintf("%v", value)
    coreErr := coreerrors.HandleGetError(resourceType, field, valueStr, err)
    return convertCoreError(coreErr)
}
```

**Dual Controller System (`cmd/trex/server/controllers_core.go`):**
- **Legacy Controllers**: Existing TRex controller system for backward compatibility
- **Core Controllers**: New controller system using rh-trex-core framework
- **Event Bus Adapter**: Bridges TRex events to core library event system
- **Unified Management**: Both systems run concurrently during transition

### Core Components

**Main Application (`cmd/trex/main.go`):**
- CLI tool with subcommands: `migrate`, `serve`, `clone`
- Uses Cobra for command structure
- **Now includes** core library integration for enhanced functionality

**Environment Framework (`cmd/trex/environments/`):**
- Configurable environments: development, testing, production
- Visitor pattern for component initialization
- Service locator pattern for dependency injection
- **Enhanced with** core library session factory integration

**API Layer (`pkg/api/`):**
- OpenAPI-generated models and clients
- Example "Dinosaur" entity with CRUD operations
- **Enhanced error handling** using rh-trex-core patterns
- Standardized metadata and response structures

**Data Layer:**
- **DAO Pattern** (`pkg/dao/`): Data Access Objects for database operations
- **Database** (`pkg/db/`): GORM-based persistence with PostgreSQL
- **Core Integration**: Session factories compatible with rh-trex-core
- **Migrations** (`pkg/db/migrations/`): Database schema versioning

**Service Layer (`pkg/services/`):**
- Business logic separated from handlers
- **Enhanced error handling** with PII sanitization and constraint detection
- **Core library delegation** for consistent error patterns
- Generic service patterns for reuse
- Event-driven architecture support

**HTTP Layer (`pkg/handlers/`):**
- REST API endpoints
- Authentication/authorization middleware
- OpenAPI specification compliance
- **Improved error responses** using core library patterns

**Infrastructure:**
- **Authentication** (`pkg/auth/`): OIDC integration with Red Hat SSO
- **Clients** (`pkg/client/ocm/`): OCM (OpenShift Cluster Manager) integration
- **Configuration** (`pkg/config/`): Environment-specific settings
- **Logging** (`cmd/trex/server/logging/`): Structured logging with request middleware
- **Metrics** (`pkg/handlers/prometheus_metrics.go`): Prometheus integration
- **Core Controllers**: Event-driven processing using rh-trex-core framework

### Key Patterns

1. **Separation of Concerns**: Clear boundaries between API, service, and data layers
2. **Dependency Injection**: Service locator pattern in environments framework
3. **Code Generation**: OpenAPI specs generate client code and documentation
4. **Test-Driven Development**: Comprehensive test support with mocks and factories
5. **Core Library Integration**: Consistent error handling and database patterns
6. **Event-Driven Architecture**: Dual controller system for legacy and core processing
7. **PII Protection**: Automatic sanitization of sensitive data in error messages

## Code Generation

### How to Generate a New Kind

The generator script creates complete CRUD functionality with **event-driven architecture** and **rh-trex-core integration** for a new resource type. Generated services automatically use enhanced error handling from the core library. The process is now fully automated with no manual steps required.

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
# - Enhanced error handling using rh-trex-core patterns
# - Database migration
# - Test files and factories
# - OpenAPI specifications
# - Service locators and routing
# - Automatic controller registration for event handling
# - PII sanitization in error messages
```

### What the Generator Creates

The generator automatically creates and configures:

1. **Generated Files** (no manual editing needed):
   - `pkg/api/fizzbuzz.go` - API model
   - `pkg/api/presenters/fizzbuzz.go` - Presenter conversion functions  
   - `pkg/handlers/fizzbuzz.go` - HTTP handlers
   - `pkg/services/fizzbuzz.go` - Business logic with event handlers and rh-trex-core error handling
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
- `{{.ProjectCamelCase}}` - Project name in CamelCase for dynamic API methods (e.g., "RhTrex", "ChessApi")
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
- **Error Handling**: Uses rh-trex-core enhanced error handling with PII sanitization
- **Context Aware**: Supports request tracing
- **Core Integration**: Automatically delegates to core library error patterns

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

## Core Library Usage (rh-trex-core)

### Overview

TRex now integrates with the **rh-trex-core** library to provide consistent, enhanced functionality across all TRex-based microservices. The core library provides:

- **Enhanced Error Handling**: PII sanitization, constraint detection, standardized error responses
- **Database Utilities**: Session factories, advisory locks, transaction management  
- **Controller Framework**: Event-driven processing with PostgreSQL LISTEN/NOTIFY
- **Type-Safe Operations**: Go generics for CRUD operations and error handling

### Error Handling Integration

**Automatic PII Sanitization:**
```go
// PII fields are automatically redacted in error messages
var piiFields []string = []string{
    "username", "first_name", "last_name", "email", "address",
}

// Core library automatically sanitizes these fields
coreErr := coreerrors.HandleGetError("User", "email", "user@example.com", err)
// Result: "User with email='<redacted>' not found"
```

**Enhanced Constraint Detection:**
```go
// Core library intelligently detects database constraint violations
func handleCreateError(resourceType string, err error) *errors.ServiceError {
    coreErr := coreerrors.HandleCreateError(resourceType, err)
    return convertCoreError(coreErr)
}

// Automatically detects:
// - Unique constraint violations → Conflict errors
// - Foreign key violations → BadRequest errors  
// - Not null violations → BadRequest errors
// - Check constraint violations → BadRequest errors
```

**Service Implementation Pattern:**
```go
// Generated services automatically use core library patterns
func (s *sqlKindService) Get(ctx context.Context, id string) (*api.Kind, *errors.ServiceError) {
    kind, err := s.kindDao.Get(ctx, id)
    if err != nil {
        // Delegates to core library with automatic PII sanitization
        return nil, handleGetError("Kind", "id", id, err)
    }
    return kind, nil
}
```

### Core Library Dependencies

**Required Go Module:**
```go
// go.mod
require (
    github.com/openshift-online/rh-trex-core v0.0.0-20250711220747-a9ce95f9f591
    gorm.io/driver/postgres v1.5.9  // Updated for compatibility
    gorm.io/gorm v1.30.0
)
```

**Import Patterns:**
```go
import (
    "github.com/openshift-online/rh-trex/pkg/errors"
    coreerrors "github.com/openshift-online/rh-trex-core/errors"
    coreapi "github.com/openshift-online/rh-trex-core/api"
    corecontrollers "github.com/openshift-online/rh-trex-core/controllers"
    coredb "github.com/openshift-online/rh-trex-core/db"
)
```

### Controller Integration

**Dual Controller System:**
- **Legacy Controllers**: Existing TRex controller system (backward compatibility)
- **Core Controllers**: New rh-trex-core controller framework (future development)
- **Event Bus Adapter**: Seamless integration between systems

**Core Controller Registration:**
```go
// Core controller manager with enhanced functionality
coreManager := corecontrollers.NewControllerManager(
    coredb.NewAdvisoryLockFactory(coreSessionFactory),
    eventBus,
)

// Register with both systems during transition
s.CoreControllerManager.RegisterController(&corecontrollers.ControllerConfig{
    Source: "Dinosaurs",
    Handlers: map[coreapi.EventType][]corecontrollers.ControllerHandlerFunc{
        coreapi.CreateEventType: {func(ctx context.Context, id string) error {
            return dinoServices.OnUpsert(ctx, id)
        }},
    },
})
```

### Migration Guidelines

**For New Resources:**
- Generated code automatically uses rh-trex-core patterns
- No manual migration needed for new Kinds
- Enhanced error handling included by default

**For Existing Services:**
- Error handlers already migrated to use core library
- Gradual migration of controller system in progress
- Backward compatibility maintained throughout transition

**Best Practices:**
- Use `handleGetError`, `handleCreateError`, etc. for all database operations
- Import core library types with aliases to avoid conflicts
- Test both unit and integration tests after any core library updates
- Follow PII sanitization patterns for sensitive data

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

## TRex Clone Command

The `trex clone` command creates new microservice projects based on the TRex template, with automatic replacement of template content for the new service.

### Clone Command Usage

**Basic Clone:**
```bash
./trex clone --name my-service --destination /tmp/my-service
```

**Clone with Custom Repository:**
```bash
./trex clone --name my-service --repo github.com/myorg --destination ./my-project
```

**Parameters:**
- `--name`: Name of the new service (will replace "rh-trex" throughout the codebase)
- `--repo`: Git repository path (default: "github.com/openshift-online")
- `--destination`: Target directory for the new service (default: "/tmp/clone-test")

### Clone Process

The clone command performs the following transformations:

1. **Project Name Replacement**: Replaces `rh-trex` with the new service name throughout the codebase
2. **Repository Path Updates**: Updates import paths from `github.com/openshift-online/rh-trex` to the specified repository
3. **File and Directory Renaming**: Renames files and directories to match the new service name
4. **Template Variable Substitution**: Updates configuration files, documentation, and build scripts
5. **Core Library Preservation**: Automatically preserves `rh-trex-core` dependency imports

### Core Library Preservation

**Critical Fix Applied**: The clone command now uses **line-by-line replacement logic** to preserve `rh-trex-core` library imports. The clone process:

✅ **Preserves**: Any line containing `rh-trex-core` (imports, dependencies, references)
✅ **Replaces**: Only lines containing TRex project references without core library context
✅ **Maintains**: Correct `github.com/openshift-online/rh-trex-core` dependency

**Protected Patterns:**
- `github.com/openshift-online/rh-trex-core` - Core library import paths
- Import aliases like `coreapi "github.com/openshift-online/rh-trex-core/api"`
- Any go.mod dependency lines with `rh-trex-core`
- Source file imports from core library packages

### Complete Clone Workflow

**1. Clone the Template:**
```bash
./trex clone --name MyService --repo gitlab.com/myorg --destination /tmp/my-service
```

**2. Setup Database:**
```bash
cd /tmp/my-service
make db/setup     # Start PostgreSQL container
```

**3. Build and Migrate:**
```bash
make binary       # Build the service binary
./myservice migrate  # Run database migrations
```

**4. Run Tests:**
```bash
make test         # Run unit tests (should pass immediately)
make test-integration  # Run integration tests (with database)
```

**5. Generate New Resources:**
```bash
go run ./scripts/generator.go --kind MyResource
# Then complete the 4 manual registration steps
make test-integration  # Verify new resource works
```

### Clone Command Fix Details

**Previous Issue**: Clone command incorrectly replaced `rh-trex-core` dependencies, causing build failures.

**Applied Fix**: Modified `/cmd/trex/clone/cmd.go` to use line-by-line replacement that checks each line individually for `rh-trex-core` before applying transformations.

**Result**: All cloned projects now maintain correct core library dependencies and build successfully without manual fixes.

**Verified**: Clone → Build → Database Setup → Test workflow works end-to-end.

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

### Testing Best Practices

**Code Change Testing Protocol:**
- **After ANY code changes:** Run `make test` to verify unit tests pass
- **After major changes:** Run `make test-integration` to verify full system integration
- **After core library updates:** Run both unit and integration tests to ensure compatibility
- **Major changes include:** New features, refactoring, architecture changes, database schema changes, rh-trex-core version updates

**Why This Matters:**
- Unit tests are fast and catch basic regressions immediately
- Integration tests are slower but verify complete system functionality
- This approach balances speed with thoroughness

### Database Issues During Testing

If integration tests fail with PostgreSQL-related errors (missing columns, transaction issues), recreate the database:

```bash
# From project root directory
make db/teardown  # Stop and remove PostgreSQL container
make db/setup     # Start fresh PostgreSQL container
make test-integration  # Run tests again
```

**Note:** Always run `make` commands from the project root directory where the Makefile is located.