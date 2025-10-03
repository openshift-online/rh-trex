# Generator Usage

Complete guide to using TRex's entity generator to create CRUD operations.

## Overview

The TRex generator creates complete, production-ready CRUD operations for business entities in minutes. It follows the plugin architecture, so generated entities integrate seamlessly without manual framework edits.

## Basic Usage

Generate a new entity:

```bash
# From your TRex project root
go run ./scripts/generate/main.go --kind Product

# Always run this after entity generation
make generate
```

**Important**: Always run `make generate` after creating entities to update OpenAPI models.

## Entity Naming Conventions

The generator uses your entity name to create consistent naming throughout the system:

```bash
# Input: Product
go run ./scripts/generate/main.go --kind Product
```

**Generated Names:**
- **Go Types**: `Product`, `ProductList`, `ProductPatchRequest`
- **Database Table**: `products` (snake_case, plural)
- **API Paths**: `/products` (kebab-case, plural)
- **JSON Fields**: `product_id`, `created_at` (snake_case)
- **Plugin Directory**: `plugins/product/`
- **Service Variables**: `ProductService`, `ProductServiceLocator`

### Multi-Word Entities

```bash
# Input: ProductOrder
go run ./scripts/generate/main.go --kind ProductOrder
```

**Generated Names:**
- **Go Types**: `ProductOrder`, `ProductOrderList`
- **Database Table**: `product_orders`
- **API Paths**: `/product-orders` 
- **Plugin Directory**: `plugins/productorder/`

## What Gets Generated

For each entity, the generator creates:

### Core Files (Complete CRUD)
```
plugins/productorder/
└── plugin.go              # Complete plugin with all registrations

pkg/
├── api/productorder_types.go        # API models and JSON structures
├── dao/productorder.go              # Database operations with GORM
├── handlers/productorder.go         # HTTP request/response handling
├── services/productorder.go         # Business logic and transactions
└── presenters/productorder.go       # Response formatting

test/
├── factories/productorder.go        # Test data generation
├── mocks/productorder.go           # Mock objects for testing
└── integration/productorder_test.go # End-to-end API tests
```

### Database Integration
```
pkg/db/migrations/
└── YYYYMMDDHHMMSS_add_product_orders.go  # Database schema migration

openapi/
└── openapi.productorder.yaml             # OpenAPI specification
```

### Updated Files
The generator also updates existing files:
- **`openapi/openapi.yaml`** - Adds API paths and schema references
- **`pkg/db/migrations/migration_structs.go`** - Registers new migration

## Generated API Operations

Each entity gets complete REST API operations:

### Standard Endpoints
- **`GET /api/{service}/v1/products`** - List products with pagination
- **`GET /api/{service}/v1/products/{id}`** - Get specific product
- **`POST /api/{service}/v1/products`** - Create new product
- **`PATCH /api/{service}/v1/products/{id}`** - Update existing product
- **`DELETE /api/{service}/v1/products/{id}`** - Delete product

### Query Parameters
- **`page`** - Page number (default: 1)
- **`size`** - Page size (default: 100, max: 500)
- **`search`** - SQL-like search criteria
- **`orderBy`** - Sorting specification
- **`fields`** - Field selection

### Example API Usage
```bash
# List all products
curl http://localhost:8000/api/inventory-api/v1/products

# Create a product
curl -X POST http://localhost:8000/api/inventory-api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Laptop", "price": 999.99}'

# Search products
curl "http://localhost:8000/api/inventory-api/v1/products?search=name like 'Laptop%'"

# Paginate results
curl "http://localhost:8000/api/inventory-api/v1/products?page=2&size=10"
```

## Plugin Architecture Integration

The generated plugin file contains all framework registrations:

```go
// plugins/product/plugin.go
package product

import (
    "github.com/your-org/inventory-api/cmd/inventory-api/environments"
    "github.com/your-org/inventory-api/pkg/services"
)

// Service locator type
type ProductServiceLocator func() services.ProductService

// Helper function for type-safe service access
func ProductService(s *environments.Services) services.ProductService {
    if s == nil {
        return nil
    }
    if obj := s.GetService("Products"); obj != nil {
        locator := obj.(ProductServiceLocator)
        return locator()
    }
    return nil
}

func init() {
    // Service registration - automatic discovery
    registry.RegisterService("Products", func(env interface{}) interface{} {
        return NewProductServiceLocator(env.(*environments.Env))
    })
    
    // HTTP routes registration - automatic discovery
    server.RegisterRoutes("products", func(router *mux.Router, env *environments.Env) {
        // Route setup with authentication and authorization
    })
    
    // Event controller registration - automatic discovery
    server.RegisterController("Products", func(env *environments.Env) {
        // Event handlers for create/update/delete
    })
    
    // API presenters registration - automatic discovery
    presenters.RegisterPath(api.Product{}, "products")
    presenters.RegisterKind(api.Product{}, "Product")
}
```

## Customizing Generated Entities

After generation, customize the business logic:

### 1. Business Logic (Service Layer)
```go
// pkg/services/product.go
func (s *productService) Create(ctx context.Context, request *api.Product) (*api.Product, error) {
    // Add your business validation
    if request.Price < 0 {
        return nil, errors.BadRequest("Price cannot be negative")
    }
    
    // Add business logic
    request.Status = "active"
    request.SKU = generateSKU(request.Name)
    
    // Database operation (generated)
    return s.productDao.Create(ctx, request)
}
```

### 2. Database Schema (Migration)
```go
// pkg/db/migrations/YYYYMMDDHHMMSS_add_products.go
func (m *AddProducts) Up() error {
    return m.db.AutoMigrate(&api.Product{})
}

// Add custom database indexes, constraints
func (m *AddProducts) Up() error {
    err := m.db.AutoMigrate(&api.Product{})
    if err != nil {
        return err
    }
    
    // Add custom index
    return m.db.Exec("CREATE INDEX idx_products_sku ON products(sku)").Error
}
```

### 3. API Models (Data Structures)
```go
// pkg/api/product_types.go
type Product struct {
    api.ObjectReference
    Name        string    `json:"name" gorm:"index"`
    Price       float64   `json:"price"`
    SKU         string    `json:"sku" gorm:"uniqueIndex"`
    Status      string    `json:"status" gorm:"default:'active'"`
    Category    string    `json:"category"`
    Description string    `json:"description"`
}
```

## Development Workflow

Standard workflow after generating entities:

```bash
# 1. Generate entity
go run ./scripts/generate/main.go --kind Product
make generate  # Update OpenAPI models

# 2. Run tests to verify generation
make test

# 3. Customize business logic
# Edit pkg/services/product.go
# Edit pkg/api/product_types.go

# 4. Run database migration
./trex migrate

# 5. Test your customizations
make test
make test-integration

# 6. Test API manually
make run
curl http://localhost:8000/api/inventory-api/v1/products
```

## Advanced Usage

### Custom Field Types
```go
// pkg/api/product_types.go
type Product struct {
    api.ObjectReference
    Name       string             `json:"name"`
    Price      decimal.Decimal    `json:"price" gorm:"type:decimal(10,2)"`
    Tags       pq.StringArray     `json:"tags" gorm:"type:text[]"`
    Metadata   map[string]string  `json:"metadata" gorm:"type:jsonb"`
    CreatedBy  string            `json:"created_by"`
}
```

### Custom Validation
```go
// pkg/services/product.go
func (s *productService) validateProduct(product *api.Product) error {
    if product.Name == "" {
        return errors.BadRequest("Product name is required")
    }
    if product.Price <= 0 {
        return errors.BadRequest("Product price must be positive")
    }
    if len(product.SKU) < 3 {
        return errors.BadRequest("SKU must be at least 3 characters")
    }
    return nil
}
```

### Custom Queries
```go
// pkg/dao/product.go
func (d *productDao) FindByCategory(ctx context.Context, category string) ([]*api.Product, error) {
    var products []*api.Product
    err := d.sessionFactory.New(ctx).Where("category = ?", category).Find(&products).Error
    return products, err
}
```

## Next Steps

- **[Plugin Architecture](plugin-architecture.md)** - Understand how plugins work
- **[Customizing Entities](customizing-entities.md)** - Add business logic and validation
- **[Testing Entities](testing-entities.md)** - Unit and integration testing patterns
- **[Operations](../operations/)** - Deploy your service with new entities

## Troubleshooting

### Generation Fails
```bash
# Ensure you're in project root
pwd  # Should show your TRex project directory

# Check entity name format
go run ./scripts/generate/main.go --kind ProductOrder  # PascalCase
```

### Build Fails After Generation
```bash
# Update OpenAPI models (required step)
make generate

# Clean and rebuild
go mod tidy
make binary
```

### Tests Fail
```bash
# Run database migrations first
./trex migrate

# Run tests
make test
```

For detailed troubleshooting, see **[Development Problems](../troubleshooting/development-problems.md)**.