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
./trex clone --repo-base github.com/myorg --name my-service   # Custom git repo
```

**Configuration Options:**
- `--name` - Name of the new service (default: "rh-trex")
- `--destination` - Target directory for new instance (default: "/tmp/clone-test")
- `--repo-base` - Git Repository base URL (default: "github.com/openshift-online")

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
   - `plugins/fizzbuzzs/plugin.go` - Plugin with routes, controllers, presenters, and service locator

2. **Updated Files** (automatically modified by generator):
   - `cmd/trex/main.go` - Adds plugin import to trigger auto-registration
   - `pkg/db/migrations/migration_structs.go` - Adds migration to MigrationList automatically
   - `openapi/openapi.yaml` - Adds API references

3. **Regenerated OpenAPI Client** (via `make generate`):
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
1. **Plugin-based architecture** - entities are self-contained with auto-registration
2. **Automatic migration registration** - adds migrations to migration_structs.go automatically
3. **Auto-discovery** - plugins register routes, controllers, and presenters via init() functions
4. **Dynamically detect** command directory structure
5. **Generate correct** snake_case API paths
6. **Use proper** service locator patterns with lock factories
7. **Create complete** test suites with proper factory methods
8. **Maintain consistency** with existing codebase patterns
9. **Generate event-driven controllers** with idempotent handlers

**Zero manual steps required** - the generator is fully automated!

### Post-Generation Workflow

After running the generator, simply build and test:

```bash
# 1. Build the binary
make binary

# 2. Set up the database
make db/teardown
make db/setup

# 3. Run migrations (your new migration is already registered)
./trex migrate

# 4. Run the server (routes and controllers are auto-registered via plugin)
make run-no-auth

# 5. Test the new entity
curl -X POST http://localhost:8000/api/rh-trex/v1/{kinds} \
  -H "Content-Type: application/json" \
  -d '{"species": "example"}' | jq

curl http://localhost:8000/api/rh-trex/v1/{kinds} | jq
```

No manual file edits required - everything is wired up automatically through the plugin system.

### Adding Custom Fields to Entities

The generator supports two approaches for adding custom fields to entities:

#### Option 1: Specify Fields at Generation Time (Recommended)

Use the `--fields` flag to specify custom fields when generating the entity:

```bash
# All fields nullable by default
go run ./scripts/generator.go --kind Rocket \
  --fields "name:string,fuel_type:string,max_speed:int,active:bool"

# Mix of required and nullable fields
go run ./scripts/generator.go --kind Rocket \
  --fields "name:string:required,fuel_type:string,max_speed:int:optional,active:bool"
```

**Supported Field Types:**
- `string` - Text data
- `int` - 32-bit integer
- `int64` - 64-bit integer
- `bool` - Boolean true/false
- `float` or `float64` - Floating point numbers
- `time` - Timestamp (time.Time)

**Field Nullability:**
- **Default**: Fields are nullable (pointer types like `*string`, `*int`)
- **`:required`**: Makes field non-nullable (base types like `string`, `int`)
- **`:optional`**: Explicitly marks as nullable (same as default)
- Required fields are added to OpenAPI `required` array
- Nullable fields use pointer types in Go structs
- All fields in PatchRequest are pointers (for partial updates)

**Examples:**
```bash
# name is required (string), others nullable (*string, *int)
--fields "name:string:required,description:string,count:int"

# All required (no pointers)
--fields "name:string:required,count:int:required,active:bool:required"

# All nullable (default, with pointers)
--fields "name:string,count:int,active:bool"
```

**Field Naming:**
- Use snake_case when specifying field names (e.g., `fuel_type`, `max_speed`)
- Generator automatically converts to proper casing:
  - Go struct fields: PascalCase (`FuelType`, `MaxSpeed`)
  - JSON/API fields: snake_case (`fuel_type`, `max_speed`)
  - Database columns: snake_case (`fuel_type`, `max_speed`)

The generator automatically adds these fields to:
- API model struct (with correct pointer/non-pointer types)
- Database migration
- OpenAPI specification (with `required` array for non-nullable fields)
- Presenter conversion functions (with proper nil handling)
- Test factories (with pointer helpers for nullable fields)
- Integration tests (with appropriate test values)
- PatchRequest struct (all fields as optional pointers)

#### Option 2: Add Fields Manually Post-Generation

If you need to add fields after the entity is generated, update these 5 files:

**1. API Model** (`pkg/api/{kind}.go`):
```go
type Rocket struct {
    Meta
    Name      string    `json:"name"`
    FuelType  string    `json:"fuel_type"`
    MaxSpeed  int       `json:"max_speed"`
    Active    bool      `json:"active"`
    LaunchDate time.Time `json:"launch_date"`
}

type RocketPatchRequest struct {
    Name       *string    `json:"name,omitempty"`
    FuelType   *string    `json:"fuel_type,omitempty"`
    MaxSpeed   *int       `json:"max_speed,omitempty"`
    Active     *bool      `json:"active,omitempty"`
    LaunchDate *time.Time `json:"launch_date,omitempty"`
}
```

**2. Database Migration** (`pkg/db/migrations/xxx_add_rockets.go`):
```go
func addRockets() *gormigrate.Migration {
    type Rocket struct {
        Model
        Name       string
        FuelType   string
        MaxSpeed   int
        Active     bool
        LaunchDate time.Time
    }
    // ... rest of migration
}
```

**3. OpenAPI Specification** (`openapi/openapi.rockets.yaml`):
```yaml
components:
  schemas:
    Rocket:
      allOf:
        - $ref: 'openapi.yaml#/components/schemas/ObjectReference'
        - type: object
          properties:
            name:
              type: string
            fuel_type:
              type: string
            max_speed:
              type: integer
              format: int32
            active:
              type: boolean
            launch_date:
              type: string
              format: date-time

    RocketPatchRequest:
      type: object
      properties:
        name:
          type: string
        fuel_type:
          type: string
        max_speed:
          type: integer
          format: int32
        active:
          type: boolean
        launch_date:
          type: string
          format: date-time
```

**4. Regenerate OpenAPI Client:**
```bash
make generate
```

**5. Add Validation (Optional)** in `pkg/handlers/{kind}.go`:
```go
func (h RocketHandler) Create(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...

    // Add custom validation
    if rocket.Name == "" {
        errors.GeneralError(r, w, errors.ErrorBadRequest, "name cannot be empty")
        return
    }

    if rocket.MaxSpeed < 0 {
        errors.GeneralError(r, w, errors.ErrorBadRequest, "max_speed must be positive")
        return
    }

    // ... rest of handler ...
}
```

**Important Notes:**
- Always use PascalCase for Go struct field names
- Use snake_case for JSON tags and database columns
- Use pointer types (*string, *int, etc.) in PatchRequest for optional updates
- After adding fields manually, recreate the database for integration tests:
  ```bash
  make db/teardown
  make db/setup
  ./trex migrate
  ```

### Generator Troubleshooting

If you encounter issues after running the generator, check these common problems:

#### Compilation Errors

**Issue**: Build fails with compilation errors
**Solution**: Verify the generator completed successfully and run:
```bash
make binary
```

If errors persist, check that the plugin import was added correctly to `cmd/trex/main.go`.

#### Database Issues

**Issue**: Migration fails or tests fail with "relation does not exist"
**Solution**: Recreate the database to apply new migrations:
```bash
make db/teardown  # Stop and remove PostgreSQL container
make db/setup     # Start fresh PostgreSQL container
./trex migrate    # Apply all migrations
make test-integration  # Run tests with new schema
```

**Note**: Always run `make` commands from the project root directory where the Makefile is located.

#### Cleaning Up Test Generations

When experimenting with the generator, you may need to completely remove a generated Kind. Here's the cleanup process:

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
  plugins/testWidgets/

# Remove OpenAPI client files (generated by make generate)
rm -rf \
  pkg/api/openapi/model_test_widget*.go \
  pkg/api/openapi/docs/TestWidget*.md

# Reset modified files to clean state
git checkout HEAD -- \
  cmd/trex/main.go \
  pkg/db/migrations/migration_structs.go \
  openapi/openapi.yaml

# Regenerate OpenAPI client to remove traces
make generate
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
