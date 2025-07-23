# Common Issues

Frequently encountered problems and their solutions.

## Quick Fixes

### "Port 5432 already in use"
```bash
# Find what's using the port
sudo netstat -tulpn | grep 5432
# OR
sudo ss -tulpn | grep 5432

# Stop system PostgreSQL
sudo systemctl stop postgresql

# Stop Docker/Podman containers
podman ps | grep postgres
podman stop <container-name>

# Try database setup again
make db/setup
```

### "Binary not found" or "command not found"
```bash
# Build the binary first
make binary

# For TRex
./trex migrate
./trex serve

# For cloned service
./your-service migrate
./your-service serve
```

### "Module not found" errors after cloning
```bash
cd ~/projects/your-cloned-service

# Clean and update dependencies
go mod tidy
go clean -cache

# Rebuild
make binary
```

### "Database connection refused"
```bash
# Check if database container is running
podman ps | grep db
# OR
docker ps | grep db

# If not running, start it
make db/setup

# If still failing, check port conflicts (see port 5432 issue above)
```

### "Migration failed" or database schema issues
```bash
# Reset database completely
make db/teardown
make db/setup

# Run migrations
./trex migrate  # or ./your-service migrate

# If still failing, check migration files in pkg/db/migrations/
```

### Tests failing after entity generation
```bash
# Always run this after generating entities
make generate

# Update dependencies
go mod tidy

# Run migrations
./trex migrate

# Run tests
make test
```

## Entity Generation Issues

### "Generator fails to find project"
```bash
# Ensure you're in the project root directory
pwd  # Should show your TRex or cloned service directory
ls   # Should see Makefile, go.mod, pkg/, cmd/, etc.

# Run generator from project root
go run ./scripts/generate/main.go --kind YourEntity
```

### "Build fails after entity generation"
```bash
# ALWAYS run this after generating entities
make generate

# If still failing, check for syntax errors in generated files
# Common locations: pkg/api/, pkg/services/, plugins/
```

### "Plugin not loading" or "Service not found"
```bash
# Check plugin file was created
ls plugins/yourentity/plugin.go

# Verify plugin has init() function
grep -n "func init()" plugins/yourentity/plugin.go

# Check main.go imports the plugin
grep -n "yourentity" cmd/*/main.go
```

## Clone-Related Issues

### "Clone corrupts itself" or "Invalid Go syntax in cloned files"
This indicates the clone process is modifying its own files.

```bash
# Rebuild TRex from clean state
cd rh-trex
git checkout .  # Reset any corrupted files
make binary

# Try clone again
go run ./scripts/clone/main.go --name your-service --destination ~/projects/your-service
```

### "Integration tests fail in cloned project"
```bash
cd ~/projects/your-cloned-service

# The clone may not have updated OpenAPI client imports
# Check test/integration/dinosaurs_test.go
# Look for import like: "github.com/openshift-online/rh-trex/pkg/client/rh-trex/v1"
# Should be: "github.com/your-org/your-service/pkg/client/your-service/v1"

# Also check API method calls:
# ApiRhTrexV1* should be ApiYourServiceV1*
```

### "go mod tidy fails in cloned project"
```bash
cd ~/projects/your-cloned-service

# Check go.mod module name
head -1 go.mod  # Should be: module github.com/your-org/your-service

# If wrong, edit go.mod first line:
# module github.com/your-org/your-service

# Then run:
go mod tidy
```

## Build and Compilation Issues

### "undefined: SomeType" or "package not found"
```bash
# Update dependencies
go mod tidy
go mod download

# If still failing, check imports in the failing file
# Ensure import paths match your project structure
```

### "GORM AutoMigrate fails" or database schema errors
```bash
# Check your entity struct tags
# Example:
type Product struct {
    api.ObjectReference
    Name  string `json:"name" gorm:"index"`  # Correct
    Price string `json:price`                # Missing quotes - WRONG
}

# Run migration manually to see detailed error
./trex migrate
```

### "OpenAPI generation fails"
```bash
# Check openapi/openapi.yaml syntax
# Look for YAML syntax errors

# Regenerate OpenAPI files
make generate

# If still failing, check entity YAML files in openapi/
ls openapi/openapi.*.yaml
```

## Runtime Issues

### "Service starts but API returns 404"
```bash
# Check if routes are registered
# Look in service logs for route registration messages

# Check plugin init() function is being called
# Add debug logging to verify

# Verify API path format
curl http://localhost:8000/api/your-service/v1/entities
```

### "Authentication errors" or "Unauthorized"
```bash
# For development, disable authentication
export OCM_ENV=development

# Or configure authentication
cat secrets/ocm-service.clientId
cat secrets/ocm-service.clientSecret

# Check service logs for auth-related errors
```

### "Database queries fail" or "Table doesn't exist"
```bash
# Run migrations
./trex migrate

# Check if migrations were applied
make db/login
your-service=# \dt  # List tables
your-service=# SELECT * FROM migrations;  # Check migration history
```

## Performance Issues

### "Slow startup" or "Service takes long to respond"
```bash
# Check database connection
# Database container might be resource-constrained

# Increase container resources
podman stop your-service-db
podman run --name your-service-db -p 5432:5432 \
  --memory=1g --cpus=2 \
  -e POSTGRES_DB=your-service \
  -e POSTGRES_USER=your-service \
  -e POSTGRES_PASSWORD=your-service \
  -d postgres:13
```

### "Memory usage high"
```bash
# Check for memory leaks in business logic
# Add monitoring to identify high memory usage areas

# Monitor service resources
docker stats your-service-container
# OR
podman stats your-service-container
```

## Development Workflow Issues

### "Hot reload not working" or "Changes not reflected"
```bash
# Restart the service manually
pkill your-service  # or Ctrl+C in service terminal
make run

# For automatic reload, consider using 'air':
go install github.com/cosmtrek/air@latest
air  # instead of make run
```

### "Test data conflicts" or "Tests interfere with each other"
```bash
# Ensure tests use isolated database transactions
# Check test setup in test/helper.go

# Reset test database
export TEST_DB_NAME=your_service_test
make db/teardown
make db/setup
./trex migrate
```

## Getting More Help

### Enable Debug Logging
```bash
# Set debug log level
export LOG_LEVEL=debug

# Run service to see detailed logs
make run
```

### Check Service Status
```bash
# Health check endpoint
curl http://localhost:8083/health

# Metrics endpoint
curl http://localhost:8080/metrics

# Check for error logs
journalctl -u your-service  # If running as systemd service
```

### Diagnostic Commands
```bash
# Check service version
./trex version

# Check database connectivity
psql -h localhost -U your-service -d your-service -c "SELECT 1;"

# Check container status
podman ps -a
docker ps -a

# Check network connectivity
netstat -tulpn | grep 8000  # API port
netstat -tulpn | grep 5432  # Database port
```

## Prevention Tips

- **Always run `make test` after changes**
- **Run `make generate` after entity generation**
- **Keep dependencies updated with `go mod tidy`**
- **Use clean builds when troubleshooting: `make clean && make binary`**
- **Check logs first** before asking for help
- **Test with minimal examples** to isolate issues

## Next Steps

If common fixes don't resolve your issue:

1. **[Build Problems](build-problems.md)** - Compilation and dependency issues
2. **[Runtime Errors](runtime-errors.md)** - Service startup and configuration problems  
3. **[Development Problems](development-problems.md)** - Generator and testing issues
4. **[Framework Development](../framework-development/)** - Understanding TRex internals