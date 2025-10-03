# Troubleshooting

Common problems and solutions for TRex development and deployment.

## Quick Diagnosis

### Service Won't Start
1. **[Build Problems](build-problems.md)** - Compilation and dependency issues
2. **[Runtime Errors](runtime-errors.md)** - Service startup and configuration problems

### Development Issues
1. **[Development Problems](development-problems.md)** - Generator, cloning, and testing issues
2. **[Common Issues](common-issues.md)** - FAQ and frequent problems

## Common Quick Fixes

### Database Issues
```bash
# Database container conflicts
podman ps | grep postgres
podman stop <container-name> && podman rm <container-name>

# Database connection problems
make db/teardown && make db/setup
./trex migrate
```

### Build Issues
```bash
# Dependency problems
go mod tidy
go mod download

# Missing tools
go install gotest.tools/gotestsum@latest
```

### Generator Issues
```bash
# Service locator problems (old pattern)
# Check cmd/{project}/environments/types.go
# Check cmd/{project}/environments/framework.go

# Plugin registration problems (new pattern)
# Check plugins/{entity}/plugin.go
# Verify init() function exists
```

### Clone Issues
```bash
# Self-corruption during clone
make binary  # Rebuild TRex first
# Follow post-clone cleanup steps

# Import path problems
go mod tidy -C <destination>
# Fix integration test imports manually
```

## Problem Categories

### Build Problems
- Go compilation errors
- Missing dependencies
- Import path issues
- Module resolution failures

### Runtime Errors  
- Service startup failures
- Database connection issues
- Authentication problems
- Configuration errors

### Development Problems
- Generator failures
- Clone corruption
- Test failures
- OpenAPI generation issues

### Integration Issues
- Database migration failures
- Authentication setup
- OpenShift deployment problems
- Event processing issues

## Getting Help

1. **Check the relevant troubleshooting guide** for your problem category
2. **Search [common-issues.md](common-issues.md)** for similar problems
3. **Enable debug logging** and check service logs
4. **Verify prerequisites** are correctly installed
5. **Test with a clean environment** to isolate issues

## Prevention

- **Always run tests** after changes (`make test`)
- **Use clean builds** when troubleshooting (`make clean && make binary`)
- **Follow setup guides exactly** for consistent environments
- **Keep dependencies updated** regularly

## Next Steps

- **[Operations](../operations/)** - Deployment and maintenance guides
- **[Framework Development](../framework-development/)** - Understanding TRex internals
- **[Reference](../reference/)** - Complete technical documentation