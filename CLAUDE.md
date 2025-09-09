# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**TRex** is Red Hat's **T**rusted **R**est **EX**ample - a production-ready microservice template for rapid API development.

**Architecture**: Plugin-based system where entities are self-contained with auto-registration  
**Goal**: Get from zero to production-ready API in minutes, not days  
**Pattern**: Generate complete CRUD operations via `go run ./scripts/generate/main.go --kind EntityName`

## 🧠 For Future Claude Sessions

**Documentation Architecture**: User-journey-based organization in `/docs/` directory:
- **[Getting Started](docs/getting-started/)** - New users: choose path, understand TRex, setup
- **[Template Cloning](docs/template-cloning/)** - Create new microservices from TRex template
- **[Entity Development](docs/entity-development/)** - Add CRUD entities to existing projects  
- **[Operations](docs/operations/)** - Deploy, run, database management
- **[Troubleshooting](docs/troubleshooting/)** - Common issues, build problems, runtime errors
- **[Reference](docs/reference/)** - Complete technical specs, APIs, commands
- **[Framework Development](docs/framework-development/)** - TRex internals, contributing

**Key Insight**: Each directory has README.md navigation hub + specific guides + cross-references

## Quick Development Tasks

### Daily Development Workflow
```bash
# Start development session
make db/setup && make run

# Add new entity
go run ./scripts/generate/main.go --kind Product
make generate  # Update OpenAPI models
make test      # Verify everything works

# Create new project
go run ./scripts/clone/main.go --name my-service --destination ~/projects/my-service
```

### Essential Commands
- `make binary` - Build the trex binary
- `make run` - Run migrations and start server (localhost:8000)
- `make test` - **ALWAYS run after code changes**
- `make test-integration` - Run after major changes (slower)
- `make db/setup` - Create PostgreSQL container  
- `make db/teardown` - Remove database container

### Entity Generation (Plugin Architecture)
- `go run ./scripts/generate/main.go --kind EntityName` - Generate complete CRUD entity
- **Auto-registration**: Plugin architecture - no manual framework edits needed
- **Location**: Entities generated in `/plugins/{entity}/plugin.go`  
- **Always run**: `make generate` after entity creation to update OpenAPI models

### Template Cloning
- `go run ./scripts/clone/main.go --name project-name --destination /path` - Clone TRex for new project
- **Post-clone**: Follow cleanup steps in [troubleshooting-clones.md](docs/template-cloning/troubleshooting-clones.md)
- **Benefits**: Clean clones with no manual service locator issues

### Testing Guidelines
- **ALWAYS** run `make test` after any code changes
- **Run** `make test-integration` after major changes (slower but comprehensive)
- **Always** run `make generate` before tests if `openapi/openapi.yaml` was edited
- **Database issues**: Try `make db/teardown && make db/setup`

## 🤖 Claude Context Optimization

### Decision Trees for Common Requests

**User says "add entity" or "generate CRUD"**:
1. Verify in TRex project root (`ls` for Makefile, pkg/, cmd/)
2. Run: `go run ./scripts/generate/main.go --kind EntityName`
3. **ALWAYS** run: `make generate` (updates OpenAPI models)
4. Verify: `make test`

**User says "clone" or "new service"**:
1. Verify TRex builds: `make binary`
2. Run: `go run ./scripts/clone/main.go --name service --destination path`
3. Navigate to clone: `cd path`
4. Setup: `go mod tidy && make db/setup && ./service migrate`

**User has problems**:
1. **First check**: [docs/troubleshooting/common-issues.md](docs/troubleshooting/common-issues.md)
2. **Build issues**: [docs/troubleshooting/build-problems.md](docs/troubleshooting/build-problems.md)  
3. **Runtime issues**: [docs/troubleshooting/runtime-errors.md](docs/troubleshooting/runtime-errors.md)

### Current Plugin Architecture (Key Points)

**Plugin Pattern**: Each entity is self-contained in `/plugins/{entity}/plugin.go`
- **Auto-Registration**: Uses `init()` functions - no manual framework edits
- **Service Location**: Helper functions provide type-safe access
- **Event Integration**: CRUD operations auto-generate events
- **Complete CRUD**: API + Service + DAO + DB + Tests + OpenAPI generated

**Generator Pattern**: `templates/generate-plugin.txt` creates complete plugin file
**Service Registry**: Dynamic registration with thread-safe lookup

## 🎯 Future Claude Session Optimization

### Information Hierarchy for Fast Context Building

**Start Here Always**: `/docs/README.md` - Master navigation hub
**User Goal Detection**: 
- New to TRex → `/docs/getting-started/`
- Want new service → `/docs/template-cloning/`  
- Adding features → `/docs/entity-development/`
- Deploy/run → `/docs/operations/`
- Problems → `/docs/troubleshooting/`

### Critical Files for Claude Context
1. **This file (CLAUDE.md)** - Claude-specific guidance
2. **docs/README.md** - Master navigation with user journeys
3. **Root README.md** - Project overview with quick paths
4. **docs/getting-started/choosing-your-path.md** - Decision matrix
5. **docs/troubleshooting/common-issues.md** - First-line problem solving
6. **spec/README.md** - AI-assisted development artifacts and spec-driven methodology

### Common Anti-Patterns to Avoid
❌ **Don't** search through scattered root-level .md files (they're deprecated/moved)  
❌ **Don't** manually edit service locator files (plugin architecture is auto-registration)  
❌ **Don't** skip `make generate` after entity generation (required for OpenAPI)  
❌ **Don't** suggest manual framework integration (plugins handle this)

### Semantic Breadcrumbs for Navigation
Each directory README.md contains:
- **Purpose** - What this section covers
- **When to Use** - User scenarios
- **Quick Commands** - Essential operations
- **Next Steps** - Where to go from here
- **Cross-References** - Related sections

### Success Patterns for Fast Assistance
✅ **Do** check directory README.md files for navigation context  
✅ **Do** use cross-references to guide users to related topics  
✅ **Do** follow the decision trees in this file for common requests  
✅ **Do** reference specific guide files rather than repeating information

## Rule 3 Compliance
Per docs/template-cloning/demo-workflow.md Rule 3: "Always reset the plan on failure and try again. We want a perfect working demo."

- Plugin architecture ensures clean builds without manual cleanup
- Auto-registration eliminates service locator configuration errors  
- Generator creates complete, working plugins that integrate seamlessly
- Documentation structure supports rapid problem diagnosis and resolution