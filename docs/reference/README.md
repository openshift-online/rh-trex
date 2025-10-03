# Reference Documentation

Complete technical reference for TRex APIs, configuration, and commands.

## API Documentation

- **[API Specification](api-specification.md)** - Complete OpenAPI specification
- **[Plugin API](plugin-api.md)** - Plugin development interface

## Configuration

- **[Configuration Options](configuration-options.md)** - All environment variables and settings
- **[Command Reference](command-reference.md)** - CLI commands and options

## Technical Specifications

### OpenAPI Specification
TRex generates complete OpenAPI 3.0 specifications for all entities, including:
- Request/response schemas
- Authentication requirements
- Error responses
- Client SDK generation

### Plugin API
The plugin interface defines how entities integrate with the framework:
- Service registration
- Route registration  
- Controller registration
- Presenter registration

### Configuration System
TRex uses environment-based configuration:
- Database connection settings
- Authentication providers
- Logging configuration
- Metrics collection

### Command Line Interface
TRex provides several CLI commands:
- `trex serve` - Start the API server
- `trex migrate` - Run database migrations
- `scripts/generate/main.go` - Generate entities
- `scripts/clone/main.go` - Clone template

## Standards and Conventions

### Naming Conventions
- **Entities** - PascalCase (User, ProductOrder)
- **Database tables** - snake_case (users, product_orders)
- **API paths** - kebab-case (/users, /product-orders)
- **JSON fields** - snake_case (created_at, user_id)

### HTTP API Patterns
- **REST principles** - Resource-based URLs
- **Standard methods** - GET, POST, PATCH, DELETE
- **Consistent responses** - ObjectReference pattern
- **Error handling** - Standard error format

### Database Patterns
- **Migrations** - Version-controlled schema changes
- **Advisory locks** - Prevent concurrent migration issues
- **GORM integration** - Object-relational mapping
- **Event sourcing** - PostgreSQL NOTIFY/LISTEN

## Integration Points

### Authentication
- **JWT validation** - OIDC token verification
- **Authorization** - Role-based access control
- **Service accounts** - Machine-to-machine auth

### Observability
- **Structured logging** - JSON format with correlation IDs
- **Prometheus metrics** - Standard service metrics
- **Health checks** - Kubernetes-compatible endpoints
- **Distributed tracing** - Operation ID propagation

### Deployment
- **Container images** - Multi-stage Docker builds
- **Kubernetes manifests** - Standard deployment patterns
- **OpenShift integration** - Platform-specific features
- **CI/CD pipelines** - Automated testing and deployment

## Next Steps

- **[Getting Started](../getting-started/)** - Begin using TRex
- **[Framework Development](../framework-development/)** - Contribute to TRex
- **[Troubleshooting](../troubleshooting/)** - Resolve issues