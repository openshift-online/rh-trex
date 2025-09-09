# TRex

**TRex** is Red Hat's **T**rusted **R**est **EX**ample - a production-ready microservice template for rapid API development.

![Trexxy](rhtap-trex_sm.png)

## What is TRex?

TRex provides a complete foundation for building enterprise-grade REST APIs with built-in best practices:

- **🚀 Rapid Development** - Generate complete CRUD APIs in minutes
- **🏗️ Plugin Architecture** - Self-contained entities with auto-registration
- **🔒 Production Ready** - OIDC auth, metrics, logging, error handling
- **📊 OpenAPI First** - Auto-generated docs and client SDKs
- **🧪 Testing Built-in** - Unit and integration test frameworks
- **📦 Container Ready** - Docker and OpenShift deployment

**Goal**: Get from zero to production-ready API in minutes, not days.


## Choose Your Path

### 🏗️ I Want to Create a New Microservice
**→ [Template Cloning Guide](docs/template-cloning/)**

Clone TRex into a new project with your business domain:
```bash
go run ./scripts/clone/main.go --name my-service --destination ~/projects/my-service
```

### 🔧 I Want to Add Entities to an Existing Project
**→ [Entity Development Guide](docs/entity-development/)**

Generate complete CRUD operations for new business objects:
```bash
go run ./scripts/generate/main.go --kind Product
```

### 🎯 I Want to Explore TRex First
**→ [Local Development Guide](docs/operations/local-development.md)**

Run TRex locally to understand how it works:
```bash
make db/setup && make run
# Visit http://localhost:8000/api/rh-trex/v1/dinosaurs
```

## Complete Documentation

**📚 [Full Documentation](docs/)** - Organized by user workflow:

- **[Getting Started](docs/getting-started/)** - Choose your path and understand TRex
- **[Template Cloning](docs/template-cloning/)** - Create new microservices  
- **[Entity Development](docs/entity-development/)** - Add entities to existing projects
- **[Operations](docs/operations/)** - Deploy and run services
- **[Reference](docs/reference/)** - Technical specifications and APIs
- **[Troubleshooting](docs/troubleshooting/)** - Common problems and solutions
- **[Framework Development](docs/framework-development/)** - Contributing to TRex
- **[Spec Directory](spec/)** - AI-assisted feature development artifacts

## Architecture Overview

TRex uses a **plugin-based architecture** where each business entity is self-contained:

- **API Layer** - RESTful endpoints with authentication
- **Service Layer** - Business logic with transaction management  
- **DAO Layer** - Database operations with GORM
- **Plugin System** - Auto-registration, no manual framework edits

See **[Architecture Diagrams](docs/framework-development/architecture-diagrams.md)** for detailed technical overview.