# Getting Started with TRex

TRex is Red Hat's **T**rusted **R**est **EX**ample - a production-ready microservice template for rapid API development.

## What is TRex?

TRex provides a complete foundation for building enterprise-grade REST APIs with built-in best practices. It demonstrates a full CRUD implementation using "dinosaurs" as example business logic that you replace with your own domain models.

## Choose Your Path

### 🏗️ I Want to Create a New Microservice
**→ [Template Cloning](../template-cloning/)**

Perfect if you're starting a new project. TRex will clone itself into a new project with all the plumbing already set up.

**Time Investment:** 10-15 minutes  
**Result:** Complete new microservice ready for your business logic

### 🔧 I Want to Add Entities to an Existing TRex Project  
**→ [Entity Development](../entity-development/)**

Great if you already have a TRex-based project and want to add new business objects (entities) with full CRUD operations.

**Time Investment:** 5 minutes per entity  
**Result:** Complete API, service, and database layers for your entity

### 🎯 I Want to Explore TRex Locally
**→ [Local Development](../operations/local-development.md)**

Best if you want to understand how TRex works before committing to using it for a project.

**Time Investment:** 20-30 minutes  
**Result:** Running TRex instance you can experiment with

## What You'll Get

No matter which path you choose, TRex provides:

- **OpenAPI Specification** - Auto-generated documentation and client SDKs
- **Layered Architecture** - Clean separation: API → Service → DAO → Database
- **Code Generation** - Full CRUD scaffolding for rapid development
- **Production Ready** - OIDC authentication, metrics, logging, error handling
- **Event-Driven** - Async processing via PostgreSQL NOTIFY/LISTEN
- **Database Management** - GORM ORM with migrations and advisory locking
- **Testing** - Built-in unit and integration test framework
- **Deployment** - Container-ready with OpenShift support

## Next Steps

1. **[Understanding TRex](understanding-trex.md)** - Architecture and key concepts
2. **[First Steps](first-steps.md)** - Prerequisites and initial setup
3. **[Choosing Your Path](choosing-your-path.md)** - Detailed comparison of options

Or jump directly to the path that matches your goal:
- **[Template Cloning](../template-cloning/)** for new projects
- **[Entity Development](../entity-development/)** for existing projects
- **[Operations](../operations/)** for running and deploying