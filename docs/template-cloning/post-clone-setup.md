# Post-Clone Setup

Additional customization and setup options after successfully cloning TRex.

## Required Post-Clone Steps

These steps are essential for a working clone:

### 1. Verify Clone Success
```bash
cd ~/projects/your-service
make binary    # Must succeed
make test      # Must pass
```

If either fails, see **[Troubleshooting Clones](troubleshooting-clones.md)**.

### 2. Database Setup
```bash
# Start your service's database
make db/setup

# Run migrations
./your-service migrate

# Verify database works
make db/login
# You should see: your-service=#
```

### 3. Service Configuration
Your cloned service inherits TRex configuration but may need customization:

```bash
# Review configuration files
ls secrets/           # Database and auth configuration
cat go.mod           # Verify module name is correct
```

## Optional Customizations

### Business Domain Replacement

Replace the example "dinosaurs" with your business entities:

```bash
# Remove dinosaur entity (optional - you can keep it as example)
rm -rf plugins/dinosaurs/
rm pkg/api/dinosaur_types.go
rm pkg/dao/dinosaur.go
rm pkg/handlers/dinosaur.go
rm pkg/services/dinosaurs.go
rm pkg/presenters/dinosaur.go
rm test/factories/dinosaurs.go
rm test/integration/dinosaurs_test.go
rm openapi/openapi.dinosaurs.yaml

# Generate your first business entity
go run ./scripts/generate/main.go --kind Product
make generate
```

### Service Metadata

Update service information in key files:

```go
// openapi/openapi.yaml
info:
  title: Your Service API           # Update from "rh-trex Service API"
  description: Your Service API     # Update description
  version: 0.0.1                    # Set your version
```

```dockerfile
# Dockerfile
LABEL name="your-service" \
      description="Your service description" \
      version="0.1.0"
```

### Authentication Configuration

Configure authentication for your service:

```bash
# secrets/ocm-service.clientId - Your OIDC client ID
echo "your-client-id" > secrets/ocm-service.clientId

# secrets/ocm-service.clientSecret - Your OIDC client secret  
echo "your-client-secret" > secrets/ocm-service.clientSecret

# For development, you can disable auth
export OCM_ENV=development
```

### Database Configuration

Customize database settings:

```bash
# secrets/db.name - Database name (already updated by clone)
cat secrets/db.name    # Should show your service name

# secrets/db.host - Database host (default: localhost)
# secrets/db.port - Database port (default: 5432)  
# secrets/db.user - Database user (default: your-service)
# secrets/db.password - Database password (default: your-service)
```

### Error Code Customization

Your clone already has customized error codes, but you can refine them:

```go
// pkg/errors/errors.go - Error code prefixes
const (
    ErrorCodePrefix = "YOUR-SERVICE-MGMT-"  // Already updated by clone
)
```

## Integration with Development Tools

### IDE Configuration

Configure your development environment:

```json
// .vscode/settings.json (VS Code)
{
    "go.toolsEnvVars": {
        "GO111MODULE": "on"
    },
    "go.testEnvVars": {
        "OCM_ENV": "development"
    }
}
```

### Git Repository Setup

Initialize version control for your new service:

```bash
# Initialize git repository
git init
git add .
git commit -m "Initial service clone from TRex"

# Add remote repository
git remote add origin https://github.com/your-org/your-service.git
git push -u origin main
```

### CI/CD Pipeline

Your clone includes CI/CD configuration that may need customization:

```yaml
# .github/workflows/ or ci/ directory
# Update service names, repository URLs, deployment targets
```

## Service-Specific Configuration

### Custom Middleware

Add service-specific middleware:

```go
// cmd/your-service/server/middleware.go
func CustomAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Your custom authentication logic
    }
}
```

### Custom Metrics

Add service-specific metrics:

```go
// pkg/metrics/custom_metrics.go
var (
    ProductCreations = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "products_created_total",
            Help: "Total number of products created",
        },
        []string{"status"},
    )
)
```

### Custom Health Checks

Add service-specific health checks:

```go
// cmd/your-service/server/health.go
func CustomHealthCheck() error {
    // Check external dependencies
    // Check business logic health
    return nil
}
```

## Development Workflow Setup

### Hot Reload (Optional)

Set up automatic reloading during development:

```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Create .air.toml configuration
# Then use: air instead of make run
```

### Testing Configuration

Configure testing environment:

```bash
# Test database setup
export TEST_DB_NAME=your_service_test
make db/setup

# Run tests with coverage
make test-coverage
```

### API Documentation

Generate and serve API documentation:

```bash
# Generate OpenAPI documentation
make generate

# Serve docs locally (if you have swagger-ui)
swagger-ui-serve openapi/openapi.yaml
```

## Production Preparation

### Environment Configuration

Set up environment-specific configuration:

```bash
# config/development.yaml
# config/staging.yaml  
# config/production.yaml
```

### Container Configuration

Customize container build:

```dockerfile
# Dockerfile adjustments for your service
# Multi-stage build optimizations
# Security scanning integration
```

### Deployment Configuration

Update deployment manifests:

```yaml
# deploy/service.yaml - Kubernetes service configuration
# deploy/deployment.yaml - Pod specification
# deploy/configmap.yaml - Configuration management
```

## Next Steps

After post-clone setup:

1. **[Entity Development](../entity-development/)** - Add your business entities
2. **[Operations](../operations/)** - Deploy and run your service  
3. **[Framework Development](../framework-development/)** - Understand TRex internals
4. **[Reference](../reference/)** - Technical specifications

## Verification Checklist

Ensure your clone is ready for development:

- [ ] Service builds successfully (`make binary`)
- [ ] Tests pass (`make test`)
- [ ] Database connects (`make db/setup` and `./your-service migrate`)
- [ ] Service starts (`make run`)
- [ ] API responds (`curl http://localhost:8000/api/your-service/v1/dinosaurs`)
- [ ] OpenAPI documentation generates (`make generate`)
- [ ] Git repository initialized
- [ ] Service name appears in all configurations
- [ ] Database names are service-specific
- [ ] Error codes are service-specific

If any step fails, see **[Troubleshooting Clones](troubleshooting-clones.md)**.