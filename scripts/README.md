# TRex Scripts

This directory contains two main tools for TRex development:

## Clone Tool (`scripts/clone/`)

Creates a new project by cloning the TRex template.

**Usage from TRex root directory:**
```bash
go run -C scripts/clone . --name <project-name> --destination <path>
```

**Example:**
```bash
go run -C scripts/clone . --name user-service --destination ~/projects/user-service
```

## Generate Tool (`scripts/generate/`)

Generates new entities with full CRUD operations in an existing TRex project.

**Usage from TRex project root directory:**
```bash
go run -C scripts/generate . --kind <EntityName>
```

**Example:**
```bash
go run -C scripts/generate . --kind User
```

## Architecture

Each tool is a separate Go module with its own:
- `main.go` - CLI interface
- `*.go` - Implementation logic  
- `go.mod` - Independent module definition
- No shared dependencies or conflicts

## Benefits

✅ **Semantic naming** - Clear separation between clone and generate functionality  
✅ **Independent modules** - No package conflicts in the same `main` package  
✅ **Simple execution** - Run with standard `go run` command  
✅ **Clean separation** - Each tool has its own directory and concerns

## Development

To modify either tool:
1. Navigate to the appropriate subdirectory (`scripts/clone/` or `scripts/generate/`)
2. Edit the Go files
3. Test with `go run .`
4. No shared dependencies to break