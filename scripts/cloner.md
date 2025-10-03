# TRex Cloner Documentation

## Usage

The TRex cloner is now implemented as a modern CLI tool with two main commands:

```bash
# Build the tools
go build -o trex-tools .

# Clone TRex template to create new project
./trex-tools clone --name <project-name> --destination <path>
./trex-tools clone --name awesomeProject -- ../awesomeProject
go run -C scripts/clone . --name ocm-ai --destination ../ocm-ai

# Generate new entity with CRUD operations
./trex-tools generate --kind <EntityName>
```

### Example Usage

```bash
# Clone TRex to create a new API project
./trex-tools clone --name user-service --destination ~/projects/user-service

# Generate a User entity in an existing project
./trex-tools generate --kind User
```

## Architecture Overview

The cloner uses a **file-category-based processing approach** that eliminates string replacement interference:

### File Categories

1. **ModuleFile** (`go.mod`) - Module declaration replacements
2. **GoSourceFile** (`*.go`) - Import path and package name replacements  
3. **OpenAPIFile** (`openapi/*.yaml`) - API paths and operation ID replacements
4. **InfrastructureFile** (`Makefile`, `Dockerfile`, `templates/*`) - Service names, binary paths, database names
5. **DocumentationFile** (`*.md`) - Project references and documentation
6. **ConfigurationFile** (`*.yaml`, `*.yml`) - General config files
7. **SkipFile** - Files to copy without processing

### Processing Functions

Each file category has a dedicated processing function that applies only the relevant transformations:

- `processModuleFile()` - Updates module declaration
- `processGoSourceFile()` - Updates import paths, preserves rh-trex-core dependencies
- `processOpenAPIFile()` - Updates API paths and operation IDs
- `processInfrastructureFile()` - Updates service names, binary paths, database references
- `processDocumentationFile()` - Updates project references, adds clone metadata

## Key Features

### 🚀 **No Import Dependencies**
- Self-contained module (`trex-tools`) 
- No circular dependencies on main TRex project
- Can be run from any directory

### 🎯 **Context-Aware Replacements**
- **Binary names**: `test-api` → `testapi` (no hyphens)
- **SQL-safe names**: `test-api` → `test_api` (underscores) 
- **CamelCase**: `test-api` → `TestApi` (for OpenAPI operations)

### 🛡️ **Framework Preservation**
- Preserves `rh-trex-core` imports and dependencies
- Only replaces project-specific references
- Maintains template source attribution

### 📁 **Smart File Handling**
- Automatically skips `scripts/` directory to avoid recursion
- Excludes build artifacts (`.git`, `vendor/`, `*.log`)
- Preserves file permissions and directory structure

## Name Transformation Functions

```go
// Project name transformations
toCamelCase("test-api")    // → "TestApi" 
toSqlSafeName("test-api")  // → "test_api"
toBinaryName("test-api")   // → "testapi"
```

## Clone Process Flow

1. **Validate Configuration** - Check name and destination parameters
2. **Create Directory Structure** - Create destination with proper permissions
3. **Walk Source Tree** - Recursively process all files
4. **Categorize Files** - Determine processing type for each file
5. **Apply Transformations** - Use category-specific processing
6. **Generate Metadata** - Add clone information to CLAUDE.md
7. **Provide Next Steps** - Display commands for continued setup

## Post-Clone Setup

The cloner automatically provides next steps:

```bash
cd /path/to/cloned-project &&
go mod tidy &&
make db/setup &&
make binary &&
go run ./scripts/ generate --kind YourEntity &&
make generate &&
make test && make test-integration
```

## Testing

The cloner includes comprehensive tests:

- **File categorization tests** - Verify correct file type detection
- **Processing function tests** - Test each transformation type
- **Name transformation tests** - Validate utility functions
- **Integration tests** - End-to-end transformation validation

Run tests: `go test -v`

## Error Handling

- **Validation**: Checks required parameters before processing
- **Path Safety**: Validates destination directory creation
- **Atomic Operations**: Fails fast on file processing errors
- **Clear Messaging**: Provides specific error context

## Best Practices

### For Clone Users
1. Always use descriptive project names (e.g., `user-api`, `inventory-service`)
2. Choose destination paths outside the source TRex directory
3. Run `go mod tidy` immediately after cloning
4. Follow the provided post-clone setup commands

### For Clone Developers
1. Add new file types to `categorizeFile()` function
2. Create category-specific processing functions for new file types
3. Update tests when adding new transformation logic
4. Preserve framework dependencies (never replace `rh-trex-core`)

## Migration from Legacy Cloner

The new cloner eliminates issues from the previous version:

- ❌ **Old**: String replacement interference causing `rh-rh-trex` corruption
- ✅ **New**: Category-specific processing prevents interference

- ❌ **Old**: Circular dependency on main TRex package  
- ✅ **New**: Self-contained module with no dependencies

- ❌ **Old**: Manual post-clone cleanup required
- ✅ **New**: Automated cleanup with clear next steps

The new architecture is more maintainable, testable, and reliable for creating production-ready TRex clones.