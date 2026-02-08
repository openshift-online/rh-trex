# Run TRex Unit Tests

Run unit tests against the package and command code (no database required).

## Instructions

1. Run unit tests:
   ```bash
   make test
   ```

2. To run specific tests:
   ```bash
   TESTFLAGS="-run TestLoadServices" make test
   ```

3. Report results with pass/fail summary.

## How Unit Tests Work

- Environment: `OCM_ENV=unit_testing`
- Database: Uses `dbmocks.NewMockSessionFactory()` (no real database)
- Scope: `./pkg/...` and `./cmd/...`
- Uses `gotestsum` for output formatting

## When to Use This vs `trex.test.integration`

- **Unit tests**: Fast, no external dependencies, test individual packages
- **Integration tests**: Full stack with real database, test API endpoints end-to-end
