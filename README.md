# TRex

**TRex** is Red Hat TAP's **T**rusted **R**est e**X**ample - a production-ready microservice template for rapid API development.

![Trexxy](rhtap-trex_sm.png)

## Overview

TRex provides a complete foundation for building enterprise-grade REST APIs with built-in best practices. It demonstrates a full CRUD implementation using "dinosaurs" as example business logic that you replace with your own domain models.

## Key Features

- **OpenAPI Specification**: Auto-generated documentation and client SDKs
- **Layered Architecture**: Clean separation with API, Service, DAO, and Database layers
- **Code Generation**: Full CRUD scaffolding generator for rapid development
- **Production Ready**: OIDC authentication, metrics, logging, and error handling
- **Event-Driven**: Async processing via PostgreSQL NOTIFY/LISTEN
- **Database Management**: GORM ORM with migrations and advisory locking
- **Testing**: Built-in test framework with unit and integration tests
- **Deployment**: Container-ready with OpenShift support

## Goals

1. **Rapid Bootstrapping**: Get from zero to production-ready API in minutes
2. **Best Practices**: Enforce enterprise patterns and standards
3. **Framework Evolution**: Provide an upgrade path for future improvements
4. **Developer Experience**: Minimize boilerplate while maximizing functionality


## Getting Started

### Quick Start Options

**Option 1: Clone TRex for New Project**
```shell
# Build TRex cloning tool
make binary

# Clone TRex template to new project
./trex clone --name my-service --destination ~/projects/src/github.com/my-org/my-service
```

**Option 2: Generate New Entity in Current Project**
```shell
# Generate complete CRUD entity with API, service, DAO layers
go run ./scripts/generator.go --kind Product
```

**Option 3: Run TRex Locally**

See [RUNNING.md](./RUNNING.md) for complete setup and running instructions including:
- Building and database setup
- Running migrations and tests  
- Starting the API server
- Authentication with OCM tool
- OpenShift deployment

### Architecture

TRex follows clean architecture principles with clear separation of concerns:

- **API Layer** (`pkg/handlers/`): HTTP routing and request/response handling
- **Service Layer** (`pkg/services/`): Business logic and transaction management
- **DAO Layer** (`pkg/dao/`): Data access abstraction with GORM
- **Database Layer** (`pkg/db/`): PostgreSQL with migrations and advisory locking

See [ASCIIARCH.md](./ASCIIARCH.md) for detailed architecture diagrams.

### Code Generation

TRex includes a powerful generator that creates complete CRUD operations:

```shell
go run ./scripts/generator.go --kind EntityName
```

**Generates:**
- API models and handlers
- Service layer with business logic
- DAO with database operations
- OpenAPI specifications
- Database migrations
- Unit and integration tests
- Service locator registration

**Note**: See [CLAUDE.md](./CLAUDE.md) for known generator issues and manual fixes required.

### Development Workflow

1. **Generate Entity**: Use generator for new business objects
2. **Customize Logic**: Add business rules in service layer  
3. **Test**: Run unit tests (`make test`) and integration tests (`make test-integration`)
4. **Update API**: Modify OpenAPI specs and run `make generate`
5. **Deploy**: Use `make deploy` for OpenShift or container deployment