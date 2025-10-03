# Understanding TRex

TRex is designed around a simple principle: **minimize the time between having an idea for an API and having a working, production-ready implementation**.

## Core Architecture

### Plugin-Based Entity System

TRex uses a plugin architecture where each business entity (User, Product, Order, etc.) is completely self-contained:

```
plugins/
├── dinosaurs/          # Example entity (replace with your domain)
│   └── plugin.go       # Complete entity definition
├── users/              # Your entities
│   └── plugin.go       # Auto-generated, self-registering
└── products/
    └── plugin.go
```

**Key Benefits:**
- **No Framework Edits** - Adding entities never requires touching framework code
- **Atomic Operations** - Add/remove entire entities as single units
- **Conflict-Free Development** - Multiple developers can add entities simultaneously
- **Auto-Discovery** - Entities register themselves via Go's `init()` functions

### Layered Architecture

Each entity follows clean architecture principles:

```
┌─────────────────────┐
│   HTTP API Layer    │  ← RESTful endpoints, authentication
├─────────────────────┤
│   Service Layer     │  ← Business logic, transactions
├─────────────────────┤
│   DAO Layer         │  ← Database operations, GORM
├─────────────────────┤
│   Database Layer    │  ← PostgreSQL, migrations
└─────────────────────┘
```

## What TRex Provides Out of the Box

### 🔐 Authentication & Authorization
- **JWT Validation** - OIDC token verification
- **Role-Based Access** - Fine-grained permissions
- **Service Accounts** - Machine-to-machine authentication

### 📊 Observability
- **Structured Logging** - JSON logs with correlation IDs
- **Prometheus Metrics** - Standard service metrics
- **Health Checks** - Kubernetes-compatible endpoints
- **Distributed Tracing** - Request tracking across services

### 🗄️ Database Management
- **GORM Integration** - Object-relational mapping
- **Migrations** - Version-controlled schema evolution
- **Advisory Locks** - Prevent concurrent migration issues
- **Event Sourcing** - PostgreSQL NOTIFY/LISTEN for async processing

### 🚀 Developer Experience
- **OpenAPI First** - API documentation and client SDK generation
- **Code Generation** - Complete CRUD scaffolding
- **Testing Framework** - Unit and integration test infrastructure
- **Container Ready** - Multi-stage Docker builds, OpenShift deployment

## Code Generation Magic

When you run:
```bash
go run ./scripts/generate/main.go --kind Product
```

TRex creates:
- **Complete API** - RESTful endpoints with proper HTTP status codes
- **Business Logic** - Service layer with transaction management
- **Database Layer** - DAO with GORM models and operations
- **Database Migration** - Schema changes with rollback support
- **OpenAPI Specification** - Auto-generated API documentation  
- **Test Infrastructure** - Unit tests, integration tests, mock factories
- **Plugin Integration** - Self-registering plugin file

All following established patterns and best practices.

## Template Cloning Power

When you run:
```bash
go run ./scripts/clone/main.go --name inventory-api --destination ~/projects/inventory-api
```

TRex creates:
- **Complete New Project** - Independent Go module
- **Customized Configuration** - Database names, API paths, error codes
- **Clean Integration** - No references to original TRex project
- **Ready to Deploy** - Container images, OpenShift manifests
- **Example Entity** - Replace "dinosaurs" with your domain

## Design Philosophy

### Rapid Bootstrapping
Get from idea to working API in minutes, not hours or days.

### Production Ready by Default
Every generated API includes authentication, logging, metrics, error handling, and testing.

### Plugin Architecture
Entities are self-contained plugins that integrate seamlessly without framework modifications.

### OpenAPI First
API specification drives code generation, ensuring consistency between documentation and implementation.

### Clean Architecture
Clear separation of concerns makes code maintainable and testable.

## When to Use TRex

**✅ Perfect For:**
- New REST API microservices
- CRUD-heavy applications
- Services requiring enterprise features (auth, metrics, logging)
- Teams wanting consistent API patterns
- Rapid prototyping with production-ready output

**❌ Consider Alternatives For:**
- Non-REST APIs (GraphQL, gRPC-only)
- Purely event-driven services (no HTTP API)
- Services with complex business logic that doesn't fit CRUD patterns
- Teams with established, different architectural patterns

## Next Steps

- **[First Steps](first-steps.md)** - Set up your development environment
- **[Choosing Your Path](choosing-your-path.md)** - Detailed comparison of cloning vs entity generation
- **[Template Cloning](../template-cloning/)** - Create a new microservice
- **[Entity Development](../entity-development/)** - Add entities to existing projects