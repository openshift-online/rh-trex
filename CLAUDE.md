# CLAUDE.md
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
- `make test` - Run unit tests (ALWAYS run after any code changes)
- `make test-integration` - Run integration tests (run after major changes - slower)
- `make ci-test-unit` - Run unit tests with JSON output for CI
- `make ci-test-integration` - Run integration tests with JSON output for CI

**Testing Guidelines:**
- **ALWAYS run `make test` after any code changes** to ensure nothing breaks
- **Run `make test-integration` after major changes** (new features, refactoring, etc.) as it's slower but more comprehensive
- **Always run `make generate` before running tests if @openapi/openapi.yaml was edited**
- **Always run `make test-integration` after successful Kind generation and `make generate`**
- **When tests fail due to database errors, try `make db/teardown` and `make db/setup`**

### Entity Generation
- `go run ./scripts/generator.go --kind EntityName` - Generate new entity with full CRUD operations
- **CRITICAL**: Current generator has pattern matching bug - see CLONING.md for details
- After generation, manually verify service locator registration in types.go and framework.go
- Generator automatically creates: API model, service layer, DAO, OpenAPI spec, service locator

### Database Operations
- `make db/setup` - Create PostgreSQL container for local development
- `make db/migrate` - Run database migrations
- `make db/teardown` - Remove database container

### TRex Cloning
- `go run ./scripts/cloner.go --name project-name --destination /path` - Clone TRex template for new project
- **NEW**: Cloning logic moved to standalone script (no longer part of trex binary)
- **BENEFIT**: Generated clones automatically exclude scripts directory (cloner.go, generator.go)
- **CRITICAL**: Follow post-clone cleanup steps in CLONING.md
- **CRITICAL**: Generator service locator registration currently has bugs - manual fixes required

## Known Issues

### Generator Service Locator Bug
**Problem**: The enhanced generator (fixed for service locator registration) uses overly broad pattern matching that corrupts multiple struct definitions.

**Impact**: 
- Service locator fields inappropriately added to ApplicationConfig, Database, Handlers, Clients structs
- Framework.go has scattered service initializations instead of centralized LoadServices()
- Results in build failures and corrupted code structure

**Workaround**: Manual cleanup required after entity generation:
1. Remove service locator fields from non-Services structs in cmd/{project}/environments/types.go
2. Remove scattered service initialization from cmd/{project}/environments/framework.go
3. Ensure only LoadServices() method contains service initialization calls
4. Verify Services struct contains correct service locator fields only

**Root Cause**: Functions `addServiceLocatorToTypes()` and `addServiceLocatorToFramework()` in scripts/generator.go use generic `"}"` pattern matching instead of targeting specific code structures.

### Testing Service Locator Fix
After entity generation, verify correct structure:

**types.go Services struct should look like:**
```go
type Services struct {
    Dinosaurs DinosaurServiceLocator
    Generic   GenericServiceLocator  
    Events    EventServiceLocator
    YourEntity YourEntityServiceLocator  // Only new entities here
}
```

**framework.go LoadServices() should look like:**
```go
func (e *Env) LoadServices() {
    e.Services.Generic = NewGenericServiceLocator(e)
    e.Services.Dinosaurs = NewDinosaurServiceLocator(e)
    e.Services.Events = NewEventServiceLocator(e)
    e.Services.YourEntity = NewYourEntityServiceLocator(e)  // Only new entities here
}
```

## Rule 3 Compliance
Per INSTANTAPI.md Rule 3: "Always reset the plan on failure and try again. We want a perfect working demo."

- Generator bugs violate Rule 3 if they prevent successful builds
- Manual cleanup is acceptable as interim solution
- Proper generator fixes required for sustainable development