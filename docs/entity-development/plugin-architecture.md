# TRex Plugin Architecture

## Overview

TRex uses a consolidated plugin architecture where each business entity is defined in a single plugin file. This eliminates manual framework editing and provides true drop-in functionality for entity generation.

## Architecture Principles

### Single Plugin File Design
Each entity lives in `plugins/{entity}/plugin.go` containing all framework registrations:
- **Service registration** with the service locator pattern
- **HTTP routes** with authentication and authorization
- **Event controllers** for async processing
- **API presenters** for path and kind mapping

### Auto-Discovery Pattern
All plugins are automatically discovered via Go's `init()` functions. No manual framework edits required when adding new entities.

### Type-Safe Service Access
Plugin files provide helper functions for type-safe service access while maintaining dynamic registration capabilities.

## Plugin Structure

### Directory Layout
```
plugins/
├── dinosaurs/          # Reference implementation
│   └── plugin.go       # All dinosaur registrations
├── {entity}/           # Your entities here
│   └── plugin.go       # Generated plugin file
```

### Plugin File Template
```go
package {entity}

// Service Locator Type
type {Kind}ServiceLocator func() services.{Kind}Service

// Service Helper Function (avoids circular imports)
func {Kind}Service(s *environments.Services) services.{Kind}Service {
    if s == nil {
        return nil
    }
    if obj := s.GetService("{KindPlural}"); obj != nil {
        locator := obj.({Kind}ServiceLocator)
        return locator()
    }
    return nil
}

func init() {
    // Service registration
    registry.RegisterService("{KindPlural}", func(env interface{}) interface{} {
        return New{Kind}ServiceLocator(env.(*environments.Env))
    })
    
    // Routes registration
    server.RegisterRoutes("{kindLowerPlural}", func(...) {
        // Route setup with authentication
    })
    
    // Controller registration  
    server.RegisterController("{KindPlural}", func(...) {
        // Event handler setup
    })
    
    // Presenter registration
    presenters.RegisterPath(api.{Kind}{}, "{kindSnakeCasePlural}")
    presenters.RegisterKind(api.{Kind}{}, "{Kind}")
}
```

## Framework Integration

### Service Registry
- **Dynamic registration**: `registry.RegisterService()` with string keys
- **Type-safe access**: Helper functions in plugin files avoid circular imports
- **Auto-discovery**: `LoadDiscoveredServices()` scans all registered services

### Route Registry
- **HTTP routing**: `server.RegisterRoutes()` with Gorilla Mux integration
- **Middleware**: Built-in JWT authentication and authorization
- **Auto-discovery**: `LoadDiscoveredRoutes()` configures all routes

### Controller Registry
- **Event handling**: `server.RegisterController()` with PostgreSQL NOTIFY/LISTEN
- **Event types**: Create, Update, Delete handlers
- **Auto-discovery**: `LoadDiscoveredControllers()` registers all event handlers

### Presenter Registry
- **Path mapping**: `presenters.RegisterPath()` for URL generation
- **Kind mapping**: `presenters.RegisterKind()` for response type identification
- **Auto-discovery**: Built-in registry pattern in presenter files

## Generator Integration

### Single Template
The generator uses one unified template `templates/generate-plugin.txt` that replaces three previous scattered templates.

### Generation Command
```bash
go run ./scripts/generator.go --kind ProductOrder
# Creates: plugins/productorder/plugin.go (drop-in)
```

### Generated Files
- **Drop-in files** (11): API model, DAO, service, handler, presenter, migration, tests, mocks, OpenAPI spec
- **Edited files** (2): `openapi/openapi.yaml` for API path registration, `pkg/db/migrations/migration_structs.go` for migration registration

## Benefits

### Developer Experience
- **Single source of truth**: All entity logic in one plugin file
- **Atomic operations**: Add/remove entire entities as units
- **Zero framework edits**: No manual file modifications for new entities
- **Easy discovery**: `ls plugins/` shows all business entities

### Architecture Benefits
- **Auto-discovery**: Components discovered via `init()` functions
- **Conflict-free**: Multiple developers can add entities without merge conflicts
- **Template-driven**: Consistent structure via generator
- **Type safety**: Helper functions provide compile-time type checking

### Maintenance Benefits
- **Simplified structure**: 1 file instead of 3 scattered files per entity
- **Better organization**: Logical grouping of entity concerns
- **Clean version control**: Entity changes isolated to plugin directory

## Implementation Details

### Plugin File Requirements
- **Package name**: Must match directory name (singular entity name)
- **File name**: Must be `plugin.go`
- **Init function**: Must contain all framework registrations
- **Imports**: Only required framework packages
- **No circular imports**: Use helper functions for service access

### Framework Requirements
- **Plugin imports**: Main application imports all plugin packages
- **Registry scanning**: Framework auto-discovers registered components
- **Error handling**: Graceful handling of plugin registration failures

## Related Documentation

- **[README.md](./README.md)**: Project overview and getting started guide
- **[CLAUDE.md](./CLAUDE.md)**: Development commands and generator usage
- **[RUNNING.md](./RUNNING.md)**: Setup and deployment instructions
