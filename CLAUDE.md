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

Use the generator script to create new resource types:
```bash
go run ./scripts/generator.go --kind KindName
```

Manual updates required after generation:
- Update `pkg/api/presenters/kind.go` and `pkg/api/presenters/path.go`
- Add service locator in `cmd/trex/environments/`
- Update routes in `cmd/trex/server/routes.go`
- Add migration in `pkg/db/migrations/`
- Create test factories in `test/factories/`

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