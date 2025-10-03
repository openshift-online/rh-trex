# Choosing Your Path

Detailed comparison to help you choose between Template Cloning and Entity Development.

## Decision Matrix

| Scenario | Template Cloning | Entity Development |
|----------|------------------|-------------------|
| **Starting a new microservice** | ✅ **Perfect** | ❌ No project exists yet |
| **Adding features to existing TRex project** | ❌ Would create duplicate | ✅ **Perfect** |
| **Learning TRex architecture** | ✅ Complete example | ⚠️ Requires existing knowledge |
| **Team has TRex project already** | ❌ Unnecessary duplication | ✅ **Ideal workflow** |
| **Need custom business domain** | ✅ Replace "dinosaurs" | ✅ Add your entities |
| **Time investment** | 15-20 minutes setup | 5 minutes per entity |
| **Result** | Independent microservice | Enhanced existing service |

## Template Cloning: When and Why

### ✅ Use Template Cloning When:
- Starting a completely new microservice project
- Need all TRex infrastructure (auth, database, deployment) set up automatically  
- Want to replace "dinosaurs" with your own business domain
- Building a service that will live independently
- Team doesn't have a TRex project yet

### What You Get:
- **Complete New Project** - Independent Go module with your project name
- **Customized Configuration** - Database names, API paths, error codes updated
- **Clean Separation** - No dependencies on original TRex repository
- **Production Ready** - All TRex features: auth, metrics, logging, deployment
- **Example to Replace** - "Dinosaur" entity shows patterns, ready to be replaced

### Example Workflow:
```bash
# Create new inventory management service
go run ./scripts/clone/main.go --name inventory-api --destination ~/projects/inventory-api
cd ~/projects/inventory-api

# Replace dinosaurs with your domain
go run ./scripts/generate/main.go --kind Product
go run ./scripts/generate/main.go --kind Warehouse
# ... customize business logic
# ... deploy your service
```

## Entity Development: When and Why

### ✅ Use Entity Development When:
- Already have a TRex-based project
- Want to add new business objects (User, Product, Order, etc.)
- Team is actively developing a TRex service
- Need to extend existing service with new capabilities
- Want fastest path to new CRUD operations

### What You Get:
- **Instant CRUD API** - Complete REST endpoints in minutes
- **Plugin Integration** - Self-contained, no framework edits
- **Consistent Patterns** - Follows established project architecture
- **Full Test Coverage** - Unit and integration tests generated
- **OpenAPI Documentation** - Automatically updated specifications

### Example Workflow:
```bash
# In your existing TRex project
go run ./scripts/generate/main.go --kind Customer
make generate  # Update OpenAPI models
make test      # Verify everything works

# Add more entities as needed
go run ./scripts/generate/main.go --kind Order
go run ./scripts/generate/main.go --kind Invoice
# ... business logic development continues
```

## Detailed Comparison

### Time Investment

**Template Cloning:**
- Initial setup: 15-20 minutes
- Post-clone cleanup: 5 minutes  
- First entity generation: 5 minutes
- **Total**: ~25-30 minutes to working service

**Entity Development:**
- Per entity: 5 minutes
- **Total**: 5 minutes per business object

### Learning Curve

**Template Cloning:**
- **Gentler** - Complete working example to explore
- Shows all TRex patterns in context
- Good for understanding overall architecture
- Requires understanding of Go modules and project structure

**Entity Development:**
- **Steeper** - Assumes familiarity with TRex patterns
- Focuses on specific entity development
- Requires existing TRex project knowledge
- Good for focused feature development

### Team Collaboration

**Template Cloning:**
- Each team can have independent TRex-based services
- No coordination needed between services
- Services can evolve independently
- Good for microservice architectures

**Entity Development:**
- Multiple developers can add entities to same service
- Plugin architecture prevents merge conflicts
- Shared service evolution
- Good for service enhancement

### Maintenance

**Template Cloning:**
- Independent services, independent maintenance
- TRex framework updates require manual integration
- Service-specific customizations easier
- Full control over service evolution

**Entity Development:**
- Centralized service maintenance
- TRex framework updates benefit all entities
- Consistent patterns across all entities
- Shared infrastructure and deployment

## Quick Decision Guide

**Choose Template Cloning if you answer "yes" to:**
- Do you need a new, independent microservice?
- Are you starting from scratch?
- Do you want complete control over the service?

**Choose Entity Development if you answer "yes" to:**
- Do you already have a TRex project?
- Are you adding features to an existing service?
- Do you want the fastest path to new CRUD operations?

## Hybrid Approach

Many teams use both approaches:

1. **Start with Template Cloning** - Create your service foundation
2. **Use Entity Development** - Add business entities rapidly
3. **Clone Again** - Create related services when needed

## Next Steps

**Ready to Clone?**
→ **[Template Cloning Guide](../template-cloning/)**

**Ready to Generate?**
→ **[Entity Development Guide](../entity-development/)**

**Want to Explore First?**
→ **[Local Development](../operations/local-development.md)**