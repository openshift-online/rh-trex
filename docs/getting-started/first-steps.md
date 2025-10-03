# First Steps

Get your development environment ready for TRex development.

## Prerequisites

You need these tools installed on your system:

### Required Tools
- **Go 1.19+** - Programming language runtime
- **Docker or Podman** - For PostgreSQL database container
- **Make** - Build automation (usually pre-installed on Linux/macOS)
- **Git** - Version control

### Optional but Recommended
- **[gotestsum](https://pkg.go.dev/gotest.tools/gotestsum)** - Better test output formatting
- **[OCM CLI](https://github.com/openshift-online/ocm-cli)** - For authentication testing
- **[jq](https://stedolan.github.io/jq/)** - JSON processing for API testing

## Installation Check

Verify your environment:

```bash
# Check Go version (need 1.19+)
go version

# Check Docker/Podman
docker --version
# OR
podman --version

# Check Make
make --version

# Install gotestsum (optional but recommended)
go install gotest.tools/gotestsum@latest
```

## Clone TRex Repository

```bash
# Clone the repository
git clone https://github.com/openshift-online/rh-trex.git
cd rh-trex

# Verify you can build
make binary
```

If the build succeeds, you should see:
```
Built ./trex successfully
```

## Set Up Development Database

TRex uses PostgreSQL for data storage:

```bash
# Start PostgreSQL container
make db/setup

# Verify database is running
make db/login
# You should see a PostgreSQL prompt: rh-trex=#
```

If you get port conflicts, check for existing containers:
```bash
# Check for containers using port 5432
docker ps | grep 5432
# OR
podman ps | grep 5432

# Stop conflicting containers if needed
docker stop <container-name>
# OR  
podman stop <container-name>
```

## Test Your Setup

Run the complete verification:

```bash
# Run database migrations
./trex migrate

# Run unit tests
make test

# Start the server (in another terminal)
make run

# Test API endpoint (in another terminal)
curl http://localhost:8000/api/rh-trex/v1/dinosaurs
```

You should see:
- Database migrations complete successfully
- All tests pass
- Server starts on port 8000
- API returns JSON response (empty list initially)

## Development Workflow Setup

Create your development workflow:

```bash
# Terminal 1: Database
make db/setup

# Terminal 2: Server (auto-restarts on changes if you have air/modd)
make run

# Terminal 3: Development commands
make test                    # Run tests
go run ./scripts/generate/main.go --kind Product  # Generate entities
make generate               # Update OpenAPI
```

## Common Development Commands

Essential commands you'll use daily:

```bash
# Build and test
make binary                 # Build TRex binary
make test                   # Run unit tests (fast)
make test-integration       # Run integration tests (slower)

# Database management
make db/setup              # Start PostgreSQL container
make db/migrate            # Run database migrations  
make db/login              # Access database shell
make db/teardown           # Stop and remove database

# Code generation
go run ./scripts/generate/main.go --kind EntityName  # Generate new entity
make generate              # Update OpenAPI models from specs

# Project cloning
go run ./scripts/clone/main.go --name my-service --destination ~/projects/my-service
```

## Next Steps

You're ready to choose your development path:

1. **[Understanding TRex](understanding-trex.md)** - Learn the architecture and design principles
2. **[Choosing Your Path](choosing-your-path.md)** - Decide between cloning and entity generation
3. **[Template Cloning](../template-cloning/)** - Create a new microservice project
4. **[Entity Development](../entity-development/)** - Add entities to existing projects

## Troubleshooting Setup

### Database Won't Start
```bash
# Check for port conflicts
netstat -ln | grep 5432
# OR
ss -ln | grep 5432

# Stop conflicting services
sudo systemctl stop postgresql  # If system PostgreSQL is running
```

### Build Fails
```bash
# Update dependencies
go mod tidy
go mod download

# Clean and rebuild
rm -f trex
make binary
```

### Tests Fail
```bash
# Reset database
make db/teardown
make db/setup
./trex migrate

# Run tests again
make test
```

For more troubleshooting help, see **[Troubleshooting Guide](../troubleshooting/)**.