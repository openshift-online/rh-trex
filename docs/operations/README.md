# Operations

Deploy, run, and maintain TRex-based services.

## For Developers

- **[Local Development](local-development.md)** - Running TRex locally for development
- **[Database Management](database-management.md)** - DB setup, migrations, and troubleshooting

## For DevOps/Platform Teams

- **[Deployment Guide](deployment-guide.md)** - Production deployment patterns
- **[Monitoring & Maintenance](monitoring-maintenance.md)** - Metrics, logging, and health checks

## Quick Commands

### Local Development
```bash
# First time setup
make db/setup          # Start PostgreSQL container
make binary            # Build TRex
./trex migrate         # Run database migrations
make run               # Start the service

# Development workflow
make test              # Run unit tests
make test-integration  # Run integration tests
make db/teardown       # Clean up database
```

### Database Operations
```bash
make db/setup    # Create PostgreSQL container
make db/login    # Access database shell
make db/migrate  # Run migrations
make db/teardown # Remove database
```

### Production Deployment
```bash
make deploy      # Deploy to OpenShift
make binary      # Build for container
```

## Service Architecture

TRex services follow a standard deployment pattern:
- **API Server** - HTTP REST endpoints (port 8000)
- **Health Server** - Health checks and readiness (port 8083)
- **Metrics Server** - Prometheus metrics (port 8080)
- **PostgreSQL** - Primary database
- **Event Processing** - Async via NOTIFY/LISTEN

## Configuration

Services are configured via:
- **Environment Variables** - Runtime configuration
- **ConfigMaps** - Kubernetes/OpenShift configuration
- **Secrets** - Database credentials, tokens

## Monitoring

TRex provides built-in:
- **Health Checks** - Kubernetes readiness/liveness probes
- **Metrics** - Prometheus-compatible metrics
- **Structured Logging** - JSON logs for aggregation
- **Request Tracing** - Operation ID tracking

## Next Steps

- **[Reference](../reference/)** - Configuration options and API specifications
- **[Troubleshooting](../troubleshooting/)** - Common operational issues
- **[Framework Development](../framework-development/)** - Understanding TRex internals