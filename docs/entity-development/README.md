# Entity Development

Add new business entities with full CRUD operations to existing TRex projects.

## When to Use Entity Development

Choose entity development when you:
- Already have a TRex-based project
- Want to add new business objects (User, Product, Order, etc.)
- Need complete API, service, and database layers generated automatically
- Want to follow TRex's plugin architecture patterns

## What You'll Get

For each entity you generate:
- ✅ **API Layer** - RESTful endpoints with authentication
- ✅ **Service Layer** - Business logic with transaction management
- ✅ **DAO Layer** - Database operations with GORM
- ✅ **Database Migration** - Schema changes with rollback support
- ✅ **OpenAPI Specification** - Auto-generated API documentation
- ✅ **Plugin Integration** - Self-registering plugin file
- ✅ **Unit Tests** - Comprehensive test coverage
- ✅ **Integration Tests** - End-to-end API testing
- ✅ **Test Factories** - Mock data generation

## Quick Start

```bash
# From your TRex project root
go run ./scripts/generate/main.go --kind Product
make generate  # Update OpenAPI models
make test      # Verify everything works
```

## Guides in This Section

- **[Generator Usage](generator-usage.md)** - How to use the entity generator
- **[Plugin Architecture](plugin-architecture.md)** - Understanding TRex's plugin system
- **[Customizing Entities](customizing-entities.md)** - Adding business logic and validation
- **[Testing Entities](testing-entities.md)** - Unit and integration testing patterns

## Plugin Architecture Benefits

TRex uses a plugin-based architecture where each entity is self-contained:

- **Single File Registration** - All framework integration in one plugin file
- **Auto-Discovery** - No manual framework edits required
- **Type-Safe Access** - Helper functions provide compile-time checking
- **Easy Maintenance** - Add/remove entities as atomic units

## Development Workflow

1. **Generate Entity** - Use generator for new business objects
2. **Customize Logic** - Add business rules in service layer
3. **Update API** - Modify OpenAPI specs if needed
4. **Test** - Run unit and integration tests
5. **Deploy** - Your entity is automatically included

## Next Steps

- **[Operations](../operations/)** - Deploy and run your enhanced service
- **[Reference](../reference/)** - Technical details on APIs and configuration
- **[Framework Development](../framework-development/)** - Contribute to TRex itself

## Alternative: Template Cloning

If you don't have a TRex project yet, see **[Template Cloning](../template-cloning/)** to create one first.