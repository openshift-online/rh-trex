# Plan: TRex Clone Validation Test

## Overview
Create an end-to-end test that validates the complete TRex API creation workflow from clone to working API, ensuring the cloned project builds, runs, and passes all validation checks.

## Test Structure
**Location**: `test/e2e/clone_validation_test.go`  
**Pattern**: End-to-end integration test using Go's standard testing library  
**Approach**: Real filesystem operations with temporary directories and cleanup

## Test Steps (Sequential Dependencies)

Each step must succeed before the next step executes. If any step fails, the test stops and cleanup is performed.

### Step 1: Clone Process Validation
**Depends on**: None (initial step)  
**Actions**:
- Execute `go run ./scripts/clone/main.go` with test parameters
- Verify all expected files/directories are created
- Validate file content transformations (module names, import paths, API paths)
**Success Criteria**: Clone completes without errors, all files present with correct content

### Step 2: Post-Clone Setup Validation  
**Depends on**: Step 1 success  
**Actions**:
- `go mod tidy` succeeds
- `make binary` builds successfully  
- `make db/setup` creates database container
- Database migrations run without errors
**Success Criteria**: Binary builds, database container running, migrations complete

### Step 3: Generated Code Validation
**Depends on**: Step 2 success  
**Actions**:
- `go run ./scripts/generate/main.go --kind TestEntity` succeeds
- `make generate` updates OpenAPI specs
- Generated plugin files have correct structure
- Service registration works properly
**Success Criteria**: Entity generation completes, OpenAPI updated, plugin structure valid

### Step 4: API Functionality Validation
**Depends on**: Step 3 success  
**Actions**:
- Service starts on correct port
- API endpoints respond correctly
- CRUD operations work for generated entities
- OpenAPI spec is valid and accessible
**Success Criteria**: Service running, all API endpoints functional, CRUD operations work

### Step 5: Testing Infrastructure Validation
**Depends on**: Step 4 success  
**Actions**:
- `make test` passes all unit tests
- `make test-integration` passes integration tests
- No compilation errors or runtime failures
**Success Criteria**: All tests pass without errors

### Step 6: Cleanup
**Depends on**: Always executes (regardless of previous step success/failure) unless configured to leave resources
**Actions**:
- Stop and remove all podman containers created during test
- Remove temporary directories and files
- Clean up any background processes
- Reset environment to pre-test state
**Success Criteria**: All resources cleaned up, no leaked containers or files

## Testing Libraries
- **Standard Go testing**: Primary framework (`testing` package)
- **Existing TRex test helper**: Leverage `test.Helper` for database/API setup
- **Real filesystem operations**: Use `os.TempDir()` for isolated test environments
- **HTTP clients**: Standard `net/http` for API validation

## Test Implementation Strategy
1. Use temporary directories for isolated clone testing
2. Leverage existing TRex test infrastructure where possible
3. Real database operations using containerized PostgreSQL
4. Comprehensive cleanup in `defer` statements
5. Parallel test execution support where safe
6. Detailed error reporting with context
7. Resource cleanup optional (default: enabled)

## Benefits
- **Validates complete workflow** from template to working API
- **Catches integration issues** that unit tests miss
- **Ensures documentation accuracy** by testing real commands
- **Provides confidence** for production deployments
- **Automated regression testing** for clone process changes
