# Run TRex Integration Tests

Run integration tests against a real PostgreSQL database (provisioned automatically via testcontainers).

## Instructions

1. Build the binary to ensure everything compiles:
   ```bash
   make binary
   ```

2. Run integration tests:
   ```bash
   make test-integration
   ```

3. To run a specific test or subset:
   ```bash
   TESTFLAGS="-run TestDinosaurGet" make test-integration
   ```

4. To enable database debug logging:
   ```bash
   DB_DEBUG=true make test-integration
   ```

5. Report results with pass/fail summary.

## How Integration Tests Work

- Environment: `OCM_ENV=integration_testing`
- Database: Testcontainers automatically starts a PostgreSQL instance (no manual setup needed)
- Each test calls `test.RegisterIntegration(t)` which returns a `(*Helper, *openapi.APIClient)`
- The database is reset between tests via `helper.DBFactory.ResetDB()`
- A mock JWK cert server handles JWT validation
- Tests use gomega matchers (`Expect(...).To(...)`)

## Test Files

- `test/integration/dinosaurs_test.go` — CRUD, paging, search, advisory lock tests
- `test/integration/controller_test.go` — Concurrent event processing
- `test/integration/metadata_test.go` — Version and build info
- `test/integration/openapi_test.go` — OpenAPI spec and Swagger UI

## Adding New Test Files

Follow the pattern in `test/integration/dinosaurs_test.go`. Every test function should:
1. Call `h, client := test.RegisterIntegration(t)`
2. Create an account: `account := h.NewRandAccount()`
3. Get an authenticated context: `ctx := h.NewAuthenticatedContext(account)`
4. Use `client.DefaultAPI.ApiRhTrexV1...` methods for API calls
