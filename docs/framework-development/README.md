# Framework Development

For contributors and maintainers working on TRex itself.

## Architecture & Design

- **[Architecture Diagrams](architecture-diagrams.md)** - Visual system overview
- **[DAO Patterns](dao-patterns.md)** - Data access layer design
- **[Generator Internals](generator-internals.md)** - How code generation works

## Contributing

- **[Contributing Guide](contributing.md)** - Development workflow and standards
- **[Debugging Framework](debugging-framework.md)** - Troubleshooting TRex internals

## Core Concepts

### Plugin Architecture
TRex uses a plugin-based system where business entities register themselves automatically:
- **Auto-Discovery** - Plugins register via `init()` functions
- **Type Safety** - Helper functions provide compile-time checking
- **Conflict-Free** - Multiple developers can add entities without conflicts

### Code Generation
The generator creates complete CRUD operations:
- **Single Template** - Unified plugin template
- **Drop-in Files** - No manual framework edits required
- **Consistent Structure** - Follows established patterns

### Service Locator Pattern
Services are registered dynamically:
- **Registry System** - Thread-safe service registration
- **Dynamic Loading** - Auto-discovery of registered services
- **Helper Functions** - Type-safe service access

## Framework Evolution

### Current State
- Plugin-based entity registration
- Auto-discovery via registry system  
- Single template for entity generation
- Clean separation of concerns

### Future Directions
- Enhanced plugin capabilities
- Additional generator templates
- Improved testing infrastructure
- Extended deployment options

## Development Setup

```bash
# Clone TRex for development
git clone https://github.com/openshift-online/rh-trex.git
cd rh-trex

# Install development dependencies
go install gotest.tools/gotestsum@latest

# Set up development environment
make db/setup
make binary
make test
```

## Testing Framework Changes

```bash
# Unit tests
make test

# Integration tests  
make test-integration

# Test generator
go run ./scripts/generate/main.go --kind TestEntity
make generate
make test

# Test cloner
go run ./scripts/clone/main.go --name test-clone --destination /tmp/test
```

## Code Style

- **Semantic naming** - Clear, purpose-driven names
- **Plugin patterns** - Follow established plugin architecture
- **Error handling** - Comprehensive error messages
- **Documentation** - Code comments and user guides
- **Testing** - Unit and integration test coverage

## Next Steps

- **[Operations](../operations/)** - Deploy and test your changes
- **[Troubleshooting](../troubleshooting/)** - Debug framework issues
- **[Reference](../reference/)** - Technical API documentation