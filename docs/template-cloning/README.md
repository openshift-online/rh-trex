# Template Cloning

Create a new microservice project by cloning the TRex template.

## When to Use Template Cloning

Choose template cloning when you:
- Want to start a completely new microservice project
- Need all TRex infrastructure (auth, database, deployment) set up automatically
- Want to replace "dinosaurs" with your own business domain
- Are building a service that will live independently

## What You'll Get

A complete new project with:
- ✅ Go module configured for your project
- ✅ Database setup and migrations
- ✅ OpenAPI specification customized
- ✅ Authentication and authorization ready
- ✅ Container and deployment configurations
- ✅ CI/CD pipeline setup
- ✅ Example CRUD entity you can replace

## Quick Start

```bash
# From TRex root directory
go run ./scripts/clone/main.go --name my-service --destination ~/projects/my-service
cd ~/projects/my-service
make db/setup && make run
```

## Guides in This Section

- **[Clone Process](clone-process.md)** - Step-by-step cloning instructions
- **[Post-Clone Setup](post-clone-setup.md)** - Required cleanup and customization
- **[Troubleshooting Clones](troubleshooting-clones.md)** - Common issues and solutions
- **[Demo Workflow](demo-workflow.md)** - Instant API demonstration process

## Examples

- **[Clone Examples](clone-examples/)** - Real-world successful clones

## Next Steps After Cloning

1. **Customize your business domain** - Replace dinosaur entities with your models
2. **[Add new entities](../entity-development/)** - Generate CRUD operations for your data
3. **[Deploy your service](../operations/)** - Get it running in production

## Alternative: Entity Development

If you already have a TRex project and just want to add new business entities, see **[Entity Development](../entity-development/)** instead.