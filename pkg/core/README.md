# TRex Core Library

The TRex Core Library provides reusable patterns and frameworks for building REST API microservices. This core library can be used to create new projects that benefit from ongoing improvements without copying code.

## Architecture

### Core Components

- **`pkg/core/api/`** - Base API types and patterns
- **`pkg/core/services/`** - Generic CRUD service patterns
- **`pkg/core/controllers/`** - Event-driven controller framework
- **`pkg/core/dao/`** - Generic data access patterns
- **`pkg/core/generator/`** - Code generation framework
- **`pkg/core/template/`** - Project template system

## Key Features

### 1. Generic CRUD Services

```go
// Create a service for any resource type
type User struct {
    api.Meta
    Name  string `json:"name"`
    Email string `json:"email"`
}

userDAO := dao.NewBaseDAO[User](db)
userService := services.NewBaseCRUDService[User](userDAO, eventEmitter, "Users")

// Automatically provides: Get, Create, Replace, Delete, List, FindByIDs, OnUpsert, OnDelete
```

### 2. Event-Driven Controllers

```go
// Auto-register controllers for any service
controllers.AutoRegisterCRUDController(controllerManager, userService, "Users")

// Automatically handles CREATE/UPDATE/DELETE events
```

### 3. Generic DAOs

```go
// Create a DAO for any GORM model
userDAO := dao.NewBaseDAO[User](db)

// Automatically provides: Get, Create, Replace, Delete, List, Count, FindByIDs
```

### 4. Project Templates

```go
// Create new projects using the core library
projectConfig := template.ProjectConfig{
    Name:        "my-service",
    Module:      "github.com/myorg/my-service",
    Resources:   []ResourceConfig{...},
}

generator := template.NewProjectTemplate(projectConfig)
generator.Generate("./my-service")
```

## Usage Pattern

### For New Projects

Instead of copying TRex code, create projects that depend on the core library:

```go
// go.mod
module github.com/myorg/my-service

require (
    github.com/openshift-online/rh-trex v1.0.0
    // other dependencies
)

// main.go
import (
    "github.com/openshift-online/rh-trex/pkg/core/generator"
    "github.com/openshift-online/rh-trex/pkg/core/controllers"
)

func main() {
    // Create resource factory
    factory := generator.NewResourceFactory(db, controllerMgr, eventEmitter)
    
    // Register resources using core patterns
    userService := generator.RegisterResourceType(factory, "Users", User{})
    
    // All CRUD operations and event handling work automatically
}
```

### Benefits

1. **Shared Evolution**: All projects get improvements automatically
2. **Smaller Codebases**: Projects only contain business logic
3. **Consistency**: All projects use the same patterns
4. **Easier Updates**: `go get -u` updates the framework
5. **Better Testing**: Framework tested once, thoroughly

## Migration Path

### Phase 1: Extract Core (Current)
- Core library created within TRex project
- Patterns established and tested
- Template system created

### Phase 2: Separate Repository
- Move core library to separate repository
- Update TRex to depend on core library
- Create migration guide

### Phase 3: New Project Creation
- Use template system for new projects
- Existing projects can migrate gradually
- All projects benefit from shared improvements

## Development

The core library is designed to be:
- **Generic**: Works with any resource type using Go generics
- **Extensible**: Easy to override default behavior
- **Testable**: Provides clear interfaces for mocking
- **Consistent**: Enforces standard patterns across projects

## Future Enhancements

- Plugin system for custom business logic
- More sophisticated event handling
- Advanced query builders
- Metrics and monitoring integration
- GraphQL support
- OpenAPI generation improvements