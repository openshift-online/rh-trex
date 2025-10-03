# Clone Process

Step-by-step guide to clone TRex into a new microservice project.

## Overview

The clone process creates a complete, independent microservice by copying TRex and customizing it for your project. This takes about 15-20 minutes from start to working service.

## Prerequisites

- TRex repository cloned locally
- Go 1.19+ installed
- Docker/Podman for database
- Available port 5432 for PostgreSQL

## Step 1: Prepare for Cloning

Clean your environment and verify TRex works:

```bash
# In TRex directory
cd rh-trex

# Stop any existing databases to free port 5432
podman ps | grep postgres
podman stop <container-name> && podman rm <container-name>

# Verify TRex builds successfully
make binary
```

## Step 2: Run the Clone Command

Clone TRex to your new project:

```bash
# Clone TRex to new project
go run ./scripts/clone/main.go --name inventory-api --destination ~/projects/inventory-api

# Example destinations:
# --destination ~/projects/my-service
# --destination /path/to/projects/user-management-api
# --destination ../my-new-service
```

**Command Options:**
- `--name` - Your project name (will be used for Go module, API paths, database names)
- `--destination` - Where to create the new project directory

The cloner will:
- ✅ Copy all TRex files to destination
- ✅ Update Go module name throughout codebase
- ✅ Customize API paths (`/api/rh-trex/v1/` → `/api/inventory-api/v1/`)
- ✅ Update database names (`rh-trex` → `inventory-api`)  
- ✅ Customize error codes (`TREX-MGMT-` → `INVENTORY-API-`)
- ✅ Exclude clone/generator scripts from final project

## Step 3: Post-Clone Setup

Navigate to your new project and set it up:

```bash
cd ~/projects/inventory-api

# Clean up Go dependencies
go mod tidy

# Start the database for your new service
make db/setup

# Run database migrations
./inventory-api migrate

# Build your new service
make binary
```

## Step 4: Verify the Clone

Test that everything works:

```bash
# Run tests to verify functionality
make test

# Start your service
make run
# Should start on http://localhost:8000

# In another terminal, test the API
curl http://localhost:8000/api/inventory-api/v1/dinosaurs
# Should return empty JSON array: {"items":[],"kind":"DinosaurList","page":1,"size":0,"total":0}
```

## Step 5: Customize for Your Domain

Replace the example "dinosaurs" with your business entities:

```bash
# Generate your first business entity
go run ./scripts/generate/main.go --kind Product
make generate  # Update OpenAPI models

# Generate more entities as needed
go run ./scripts/generate/main.go --kind Warehouse
go run ./scripts/generate/main.go --kind Supplier

# Test your new entities
make test
curl http://localhost:8000/api/inventory-api/v1/products
```

## What the Clone Creates

Your new project will have:

```
inventory-api/
├── cmd/inventory-api/          # Service entrypoint (was cmd/trex/)
├── pkg/                        # Business logic (API paths updated)
├── plugins/                    # Entity plugins (ready for your entities)
├── openapi/                    # API specification (service name updated)
├── test/                       # Test infrastructure
├── go.mod                      # Module: github.com/your-org/inventory-api
├── Makefile                    # Build commands (service name updated)
├── Dockerfile                  # Container image (service name updated)
└── secrets/                    # Configuration templates
```

**Key Customizations:**
- **Module Name**: `github.com/your-org/inventory-api`
- **API Paths**: `/api/inventory-api/v1/...`
- **Database Names**: `inventory-api` database
- **Service Binary**: `./inventory-api` command
- **Container Names**: `inventory-api-db`
- **Error Codes**: `INVENTORY-API-MGMT-...`

## Project Structure Benefits

Your cloned project:
- ✅ **Independent** - No dependencies on original TRex
- ✅ **Customized** - All names and paths updated for your domain
- ✅ **Complete** - All TRex features available
- ✅ **Clean** - No clone/generator scripts included
- ✅ **Ready** - Can be deployed immediately

## Development Workflow

After cloning, your daily workflow:

```bash
# Add new business entities
go run ./scripts/generate/main.go --kind Customer
make generate

# Develop business logic
# Edit pkg/services/customer.go
# Edit plugins/customer/plugin.go

# Test changes
make test
make test-integration

# Deploy
make deploy  # Or your deployment process
```

## Next Steps

- **[Post-Clone Setup](post-clone-setup.md)** - Additional customization options
- **[Entity Development](../entity-development/)** - Add your business entities
- **[Operations](../operations/)** - Deploy and run your service
- **[Troubleshooting](troubleshooting-clones.md)** - If something goes wrong

## Common Issues

### Clone Command Fails
```bash
# Rebuild TRex first
make binary

# Ensure destination directory doesn't exist
rm -rf ~/projects/inventory-api

# Try clone again
go run ./scripts/clone/main.go --name inventory-api --destination ~/projects/inventory-api
```

### Database Port Conflicts
```bash
# Find conflicting containers
podman ps | grep 5432

# Stop them
podman stop <container-name>
podman rm <container-name>
```

### Build Failures After Clone
```bash
cd ~/projects/inventory-api

# Clean dependencies
go mod tidy
go clean -cache

# Rebuild
make binary
```

For detailed troubleshooting, see **[Troubleshooting Clones](troubleshooting-clones.md)**.