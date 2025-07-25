# Testing Guide

Comprehensive testing guide for TRex services, including unit tests, integration tests, and code coverage reporting.

## Overview

TRex uses a multi-layered testing approach:
- **Unit Tests** - Test individual components in isolation
- **Integration Tests** - Test complete API workflows with database
- **Code Coverage** - Measure test effectiveness and identify gaps
- **Test Factories** - Generate consistent test data
- **Mocks** - Isolate components during testing

## Quick Start

```bash
# Run all unit tests with coverage
make test

# Run integration tests with coverage  
make test-integration

# Generate HTML coverage reports
make coverage-html

# View function-level coverage summary
make coverage-func
```

## Test Structure

### Unit Tests
Located alongside source code with `_test.go` suffix:
```
pkg/
├── services/
│   ├── dinosaurs.go
│   └── dinosaurs_test.go        # Unit tests for dinosaur service
├── handlers/
│   ├── dinosaurs.go  
│   └── dinosaurs_test.go        # Unit tests for dinosaur handlers
└── dao/
    ├── dinosaur.go
    └── dinosaur_test.go         # Unit tests for dinosaur DAO
```

### Integration Tests
Located in dedicated test directory:
```
test/
├── integration/
│   ├── integration_test.go      # Test setup and helpers
│   ├── dinosaurs_test.go        # End-to-end dinosaur API tests
│   └── controller_test.go       # Framework controller tests
├── factories/
│   ├── factory.go               # Test data factory interface
│   └── dinosaurs.go            # Dinosaur test data generation
└── mocks/
    ├── mocks.go                 # Mock interfaces
    └── ocm.go                   # OCM authentication mocks
```

## Unit Testing

### Writing Unit Tests

Unit tests focus on individual functions and methods in isolation:

```go
// pkg/services/dinosaurs_test.go
func TestDinosaurService_Create(t *testing.T) {
    // Setup
    mockDAO := &mocks.DinosaurDAO{}
    service := services.NewDinosaurService(mockDAO)
    
    // Test data
    request := &api.Dinosaur{
        Species: "T-Rex",
    }
    
    // Mock expectations
    mockDAO.On("Create", mock.Anything).Return(nil)
    
    // Execute
    result, err := service.Create(context.Background(), request)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "T-Rex", result.Species)
    mockDAO.AssertExpectations(t)
}
```

### Running Unit Tests

```bash
# Run all unit tests
make test

# Run tests for specific package
go test ./pkg/services -v

# Run specific test function
go test ./pkg/services -v -run TestDinosaurService_Create

# Run tests with race detection
go test -race ./pkg/services

# Run tests with timeout
go test -timeout 30s ./pkg/services
```

### Unit Test Best Practices

1. **Test One Thing** - Each test should verify one specific behavior
2. **Use Mocks** - Isolate the component being tested
3. **Clear Naming** - Test names should describe the scenario and expected outcome
4. **Setup/Teardown** - Use `t.Setup()` and `t.Cleanup()` for test preparation
5. **Table-Driven Tests** - Use table tests for multiple scenarios

```go
func TestDinosaurValidation(t *testing.T) {
    tests := []struct {
        name          string
        species       string
        expectedError bool
    }{
        {"valid species", "T-Rex", false},
        {"empty species", "", true},
        {"long species", strings.Repeat("A", 300), true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            dinosaur := &api.Dinosaur{Species: tt.species}
            err := validateDinosaur(dinosaur)
            
            if tt.expectedError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Integration Testing

### Test Environment Setup

Integration tests require a database and run against the complete service:

```go
// test/integration/integration_test.go
func TestMain(m *testing.M) {
    // Setup test database
    helper.SetupTestDB()
    
    // Run tests
    code := m.Run()
    
    // Cleanup
    helper.TeardownTestDB()
    
    os.Exit(code)
}
```

### Writing Integration Tests

Integration tests verify complete API workflows:

```go
// test/integration/dinosaurs_test.go
func TestDinosaursAPI(t *testing.T) {
    // Setup test client
    client := helper.NewTestClient()
    
    t.Run("create dinosaur", func(t *testing.T) {
        // Create request
        dinosaur := factories.BuildDinosaur("T-Rex")
        
        // Send request
        response, err := client.Post("/api/rh-trex/v1/dinosaurs", dinosaur)
        require.NoError(t, err)
        
        // Verify response
        assert.Equal(t, http.StatusCreated, response.StatusCode)
        
        var created api.Dinosaur
        err = json.Unmarshal(response.Body, &created)
        require.NoError(t, err)
        
        assert.Equal(t, "T-Rex", created.Species)
        assert.NotEmpty(t, created.ID)
        assert.NotEmpty(t, created.CreatedAt)
    })
    
    t.Run("list dinosaurs", func(t *testing.T) {
        response, err := client.Get("/api/rh-trex/v1/dinosaurs")
        require.NoError(t, err)
        
        assert.Equal(t, http.StatusOK, response.StatusCode)
        
        var list api.DinosaurList
        err = json.Unmarshal(response.Body, &list)
        require.NoError(t, err)
        
        assert.GreaterOrEqual(t, list.Total, 1)
    })
}
```

### Running Integration Tests

```bash
# Setup database first
make db/setup

# Run all integration tests
make test-integration

# Run specific integration test
go test ./test/integration -v -run TestDinosaursAPI

# Run integration tests with short flag (skips slow tests)
make test-integration TESTFLAGS="-short"

# Clean up database after testing
make db/teardown
```

### Integration Test Best Practices

1. **Database State** - Each test should clean up after itself
2. **Authentication** - Use test tokens or mock authentication
3. **Real Dependencies** - Test against actual database and HTTP handlers
4. **Error Scenarios** - Test both success and failure cases
5. **Performance** - Use `-short` flag to skip slow tests during development

## Code Coverage

### Coverage Collection

TRex automatically collects code coverage when running tests:

```bash
# Unit test coverage (automatic)
make test
# Creates: coverage-unit.out

# Integration test coverage (automatic)  
make test-integration
# Creates: coverage-integration.out
```

### Coverage Reporting

Generate reports from collected coverage data:

```bash
# Generate HTML reports
make coverage-html
# Creates: coverage-unit.html, coverage-integration.html

# View function-level coverage in terminal
make coverage-func

# Manual report generation
go tool cover -html=coverage-unit.out -o custom-report.html
go tool cover -func=coverage-unit.out
```

### Coverage Output Files

Coverage files are automatically managed:
- **`coverage-unit.out`** - Unit test coverage profile
- **`coverage-integration.out`** - Integration test coverage profile  
- **`coverage-unit.html`** - Unit test HTML report
- **`coverage-integration.html`** - Integration test HTML report

### Coverage Interpretation

**Coverage Percentage**: Shows what percentage of code lines are executed by tests
- **High coverage (>80%)** - Good test coverage, fewer production bugs
- **Medium coverage (60-80%)** - Adequate coverage, room for improvement  
- **Low coverage (<60%)** - Insufficient testing, high risk

**HTML Reports**: Provide line-by-line visualization:
- **Green lines** - Covered by tests
- **Red lines** - Not covered by tests
- **Gray lines** - Not executable (comments, declarations)

**Function Coverage**: Shows per-function coverage percentages:
```bash
github.com/openshift-online/rh-trex/pkg/services/dinosaurs.go:45:    Create      85.7%
github.com/openshift-online/rh-trex/pkg/services/dinosaurs.go:67:    Get         100.0%
github.com/openshift-online/rh-trex/pkg/services/dinosaurs.go:89:    List        92.3%
github.com/openshift-online/rh-trex/pkg/services/dinosaurs.go:112:   Update      75.0%
github.com/openshift-online/rh-trex/pkg/services/dinosaurs.go:134:   Delete      100.0%
```

### Coverage Best Practices

1. **Focus on Critical Paths** - Prioritize business logic and error handling
2. **Quality over Quantity** - High coverage doesn't guarantee good tests
3. **Test Edge Cases** - Cover error conditions and boundary cases
4. **Regular Review** - Monitor coverage trends over time
5. **Team Standards** - Establish minimum coverage thresholds

### Coverage CI Integration

Coverage data is available in CI through existing test targets:
- **`make ci-test-unit`** - Generates coverage-unit.out and unit-test-results.json
- **`make ci-test-integration`** - Generates coverage-integration.out and integration-test-results.json

## Test Data Management

### Test Factories

Use factories to generate consistent test data:

```go
// test/factories/dinosaurs.go
func BuildDinosaur(species string) *api.Dinosaur {
    return &api.Dinosaur{
        Species:   species,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

func BuildDinosaurWithDefaults() *api.Dinosaur {
    return BuildDinosaur("Triceratops")
}
```

### Database Fixtures

For integration tests requiring specific database state:

```go
func setupDinosaurFixtures(t *testing.T) {
    // Create test data
    dinosaurs := []*api.Dinosaur{
        factories.BuildDinosaur("T-Rex"),
        factories.BuildDinosaur("Stegosaurus"),
    }
    
    for _, dinosaur := range dinosaurs {
        _, err := testDB.Create(dinosaur)
        require.NoError(t, err)
    }
    
    // Cleanup after test
    t.Cleanup(func() {
        testDB.Exec("DELETE FROM dinosaurs WHERE species IN (?, ?)", 
                   "T-Rex", "Stegosaurus")
    })
}
```

## Mocking

### Interface Mocking

Use testify/mock for interface mocking:

```go
// test/mocks/dinosaur_dao.go
type DinosaurDAO struct {
    mock.Mock
}

func (m *DinosaurDAO) Create(ctx context.Context, dinosaur *api.Dinosaur) error {
    args := m.Called(ctx, dinosaur)
    return args.Error(0)
}

func (m *DinosaurDAO) Get(ctx context.Context, id string) (*api.Dinosaur, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*api.Dinosaur), args.Error(1)
}
```

### External Service Mocking

Mock external dependencies like authentication:

```go
// test/mocks/ocm.go
func MockOCMAuthentication() *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        response := `{
            "sub": "test-user",
            "email": "test@example.com",
            "preferred_username": "testuser"
        }`
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(response))
    }))
}
```

## Performance Testing

### Load Testing

For performance-critical endpoints:

```go
func BenchmarkDinosaursList(b *testing.B) {
    // Setup
    client := helper.NewTestClient()
    
    // Run benchmark
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        response, err := client.Get("/api/rh-trex/v1/dinosaurs")
        if err != nil {
            b.Fatal(err)
        }
        if response.StatusCode != http.StatusOK {
            b.Fatalf("Expected 200, got %d", response.StatusCode)
        }
    }
}
```

### Memory Profiling

```bash
# Run tests with memory profiling
go test -memprofile=mem.prof ./pkg/services
go tool pprof mem.prof

# Run benchmarks with profiling
go test -bench=BenchmarkDinosaursList -memprofile=mem.prof ./test/integration
```

## Continuous Integration

### Test Commands in CI

```bash
# CI unit tests with JSON output and coverage
make ci-test-unit

# CI integration tests with JSON output and coverage
make ci-test-integration

# Generate coverage reports for artifacts
make coverage-html
```

### Test Artifacts

CI generates these test artifacts:
- **`unit-test-results.json`** - Unit test results in JUnit format
- **`integration-test-results.json`** - Integration test results in JUnit format
- **`coverage-unit.out`** - Unit test coverage profile
- **`coverage-integration.out`** - Integration test coverage profile
- **`coverage-unit.html`** - Unit test coverage HTML report
- **`coverage-integration.html`** - Integration test coverage HTML report

## Troubleshooting

### Common Test Issues

1. **Database Connection Errors**
   ```bash
   # Reset database container
   make db/teardown
   make db/setup
   ```

2. **Test Timeouts**
   ```bash
   # Increase timeout for slow tests
   go test -timeout 5m ./test/integration
   ```

3. **Race Conditions**
   ```bash
   # Run with race detector
   go test -race ./pkg/services
   ```

4. **Memory Leaks**
   ```bash
   # Run with memory profiling
   go test -memprofile=mem.prof ./pkg/services
   ```

### Coverage Issues

1. **Low Coverage on Generated Code**
   - Generated files (OpenAPI models) may show low coverage
   - Focus on business logic coverage instead

2. **Integration Test Coverage Not Collected**
   - Ensure database is running: `make db/setup`
   - Check test builds successfully: `make binary`

3. **Coverage Reports Not Generated**
   - Verify coverage files exist: `ls coverage-*.out`
   - Check Go tools are available: `go tool cover -help`

## Testing Workflows

### Development Workflow

```bash
# Daily development testing
make test                    # Quick unit tests
make test-integration       # Full integration tests
make coverage-func          # Check coverage

# Before committing
make test && make test-integration
```

### Entity Generation Testing

```bash
# After generating new entity
go run ./scripts/generate/main.go --kind Product
make generate               # Update OpenAPI models
make test                   # Verify unit tests pass
make test-integration       # Verify integration tests pass
make coverage-html          # Check coverage impact
```

### Pre-Deployment Testing

```bash
# Complete test suite
make test
make test-integration
make coverage-html

# Performance testing  
go test -bench=. ./test/integration

# Load testing (if applicable)
./scripts/load-test.sh
```

## Next Steps

- **[Command Reference](command-reference.md)** - Complete command documentation
- **[Framework Development](../framework-development/)** - Contributing to TRex testing
- **[Troubleshooting](../troubleshooting/)** - Resolving test failures
- **[Operations](../operations/)** - Production testing strategies