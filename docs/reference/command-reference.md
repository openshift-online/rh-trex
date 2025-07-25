# Command Reference

Complete reference for all TRex commands and tools.

## TRex Binary Commands

### Build Commands
```bash
# Build TRex binary
make binary

# Install to GOPATH/bin
make install

# Clean build artifacts
make clean

# Build container image
make image
```

### Service Commands
```bash
# Start API server
./trex serve
make run                    # Equivalent to: ./trex migrate && ./trex serve

# Run database migrations
./trex migrate

# Show service version
./trex version

# Show help
./trex --help
./trex serve --help
```

### Database Commands
```bash
# Start PostgreSQL container
make db/setup

# Access database shell
make db/login

# Stop and remove database container
make db/teardown

# Run migrations only
make db/migrate             # Equivalent to: ./trex migrate
```

### Testing Commands
```bash
# Run unit tests (with coverage)
make test

# Run integration tests (with coverage)
make test-integration

# Generate HTML coverage reports
make coverage-html

# Show function-level coverage summary
make coverage-func

# Run CI tests (JSON output)
make ci-test-unit
make ci-test-integration
```

### Code Generation Commands
```bash
# Generate OpenAPI models from specs
make generate

# Generate new entity
go run ./scripts/generate/main.go --kind EntityName

# Clone TRex to new project
go run ./scripts/clone/main.go --name project-name --destination /path
```

## Generator Tool

### Basic Usage
```bash
go run ./scripts/generate/main.go --kind Product
go run ./scripts/generate/main.go --kind ProductOrder
go run ./scripts/generate/main.go --kind UserAccount
```

### Options
```bash
--kind string     Entity name in PascalCase (required)
--help           Show help message
```

### Generated Files
The generator creates these files for each entity:
- `plugins/{entity}/plugin.go` - Plugin with all registrations
- `pkg/api/{entity}_types.go` - API models and JSON structures  
- `pkg/dao/{entity}.go` - Database operations
- `pkg/handlers/{entity}.go` - HTTP request/response handling
- `pkg/services/{entity}.go` - Business logic
- `pkg/presenters/{entity}.go` - Response formatting
- `test/factories/{entity}.go` - Test data generation
- `test/mocks/{entity}.go` - Mock objects
- `test/integration/{entity}_test.go` - API tests
- `pkg/db/migrations/YYYYMMDDHHMMSS_add_{entities}.go` - Migration
- `openapi/openapi.{entity}.yaml` - OpenAPI specification

### Post-Generation Steps
```bash
# ALWAYS run after entity generation
make generate

# Verify generation worked
make test
```

## Clone Tool

### Basic Usage
```bash
go run ./scripts/clone/main.go --name my-service --destination ~/projects/my-service
```

### Options
```bash
--name string         Service name (required)
--destination string  Destination directory (required)
--help               Show help message
```

### What Gets Cloned
- Complete TRex codebase
- Customized for your service name
- Updated Go module declaration
- Customized API paths and database names
- Updated error codes and configuration

### Post-Clone Steps
```bash
cd ~/projects/my-service
go mod tidy
make db/setup
./my-service migrate
make binary
make test
```

## Environment Variables

### Database Configuration
```bash
DB_HOST=localhost           # Database host
DB_PORT=5432               # Database port  
DB_NAME=service_name       # Database name
DB_USER=service_name       # Database user
DB_PASSWORD=service_name   # Database password
DB_SSLMODE=disable         # SSL mode: disable, require, verify-ca, verify-full
```

### Service Configuration
```bash
SERVER_HOST=0.0.0.0        # API server bind address
SERVER_PORT=8000           # API server port
HEALTH_CHECK_SERVER_PORT=8083  # Health check server port
METRICS_SERVER_PORT=8080   # Metrics server port

LOG_LEVEL=info             # Log level: debug, info, warn, error
OCM_ENV=development        # Environment: development, staging, production
```

### Authentication Configuration
```bash
OCM_BASE_URL=https://api.openshift.com  # OCM API base URL
OCM_CLIENT_ID=your-client-id             # OIDC client ID
OCM_CLIENT_SECRET=your-client-secret     # OIDC client secret
OCM_TOKEN=your-token                     # OCM token for development
```

### Testing Configuration
```bash
TEST_DB_NAME=service_test   # Test database name
TEST_SUMMARY_FORMAT=short   # Test output format
GINKGO_EDITOR_INTEGRATION=true  # Enable Ginkgo IDE integration
```

## Configuration Files

### Secrets Directory
```bash
secrets/
├── db.host              # Database host
├── db.port              # Database port
├── db.name              # Database name
├── db.user              # Database user
├── db.password          # Database password
├── ocm-service.clientId # OIDC client ID
├── ocm-service.clientSecret # OIDC client secret
├── ocm-service.token    # OCM token
└── sentry.key           # Sentry error tracking key
```

### OpenAPI Specifications
```bash
openapi/
├── openapi.yaml         # Main OpenAPI specification
├── openapi.dinosaurs.yaml  # Entity specifications
└── openapi.{entity}.yaml   # Generated entity specs
```

### Templates Directory
```bash
templates/
├── generate-plugin.txt      # Plugin template
├── generate-api.txt         # API model template
├── generate-dao.txt         # DAO template
├── generate-handlers.txt    # HTTP handlers template
├── generate-services.txt    # Service layer template
├── generate-presenters.txt  # Presenter template
├── generate-test.txt        # Test template
├── generate-mock.txt        # Mock template
├── generate-migration.txt   # Migration template
└── generate-openapi-kind.txt # OpenAPI template
```

## Service Endpoints

### API Endpoints
```bash
# Health check
GET /health                 # Service health status

# Metrics  
GET /metrics               # Prometheus metrics

# API base path
GET /api/{service}/v1/     # Service API root
```

### Default Ports
- **8000** - Main API server
- **8083** - Health check server  
- **8080** - Metrics server
- **5432** - PostgreSQL database

## Exit Codes

### Success
- **0** - Command completed successfully

### Errors
- **1** - General error
- **2** - Database connection error
- **3** - Migration failure
- **4** - Configuration error
- **5** - Authentication error

## Examples

### Complete Development Workflow
```bash
# Initial setup
git clone https://github.com/openshift-online/rh-trex.git
cd rh-trex
make binary
make db/setup

# Create new entity
go run ./scripts/generate/main.go --kind Product
make generate
make test

# Clone for new service
go run ./scripts/clone/main.go --name inventory-api --destination ~/projects/inventory-api
cd ~/projects/inventory-api
go mod tidy
make db/setup
./inventory-api migrate
make run
```

### Testing Workflow
```bash
# Run all tests
make test
make test-integration

# Run specific test
go test ./pkg/services -v
go test ./test/integration -v -run TestDinosaurs

# Run with coverage
make test-coverage
go tool cover -html=coverage.out
```

### Production Deployment
```bash
# Build for production
make binary
make image

# Run migrations
./trex migrate

# Start service
./trex serve
```

## Next Steps

- **[Configuration Options](configuration-options.md)** - Detailed configuration reference
- **[API Specification](api-specification.md)** - Complete API documentation
- **[Operations](../operations/)** - Service deployment and management
- **[Troubleshooting](../troubleshooting/)** - Command-related issues