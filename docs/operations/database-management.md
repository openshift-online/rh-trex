# Database Management

Complete guide to managing PostgreSQL databases in TRex-based services.

## Database Architecture

TRex uses PostgreSQL with the following architecture:
- **GORM ORM** - Object-relational mapping for Go
- **Versioned Migrations** - Forward and backward schema changes
- **Advisory Locks** - Prevent concurrent migration issues
- **Connection Pooling** - Efficient database connections
- **Event Sourcing** - PostgreSQL NOTIFY/LISTEN for async processing

## Local Development Database

### Initial Setup

```bash
# Start PostgreSQL container
make db/setup

# Verify database is running
make db/login
# You should see: service-name=#
```

The `db/setup` command:
- Creates a PostgreSQL container named `{service}-db`
- Creates database named `{service}`
- Creates user `{service}` with password `{service}`
- Exposes PostgreSQL on port 5432

### Database Commands

```bash
# Start database container
make db/setup

# Access database shell
make db/login

# Run migrations
./trex migrate               # TRex binary
./your-service migrate       # Your cloned service binary

# Stop and remove database
make db/teardown

# View database logs
podman logs {service}-db
# OR
docker logs {service}-db
```

### Database Configuration

Database settings are configured via secrets:

```bash
# Database connection settings
cat secrets/db.host      # Default: localhost
cat secrets/db.port      # Default: 5432
cat secrets/db.name      # Service name (e.g., inventory-api)
cat secrets/db.user      # Service name (e.g., inventory-api)
cat secrets/db.password  # Service name (e.g., inventory-api)
```

## Database Migrations

### How Migrations Work

TRex uses versioned migrations stored in `pkg/db/migrations/`:

```
pkg/db/migrations/
├── migration_structs.go           # Migration registry
├── 201911212019_add_dinosaurs.go  # Example migration
├── 202309020925_add_events.go     # System migration
└── YYYYMMDDHHMMSS_add_products.go # Your migrations
```

### Migration Structure

```go
// pkg/db/migrations/YYYYMMDDHHMMSS_add_products.go
package migrations

import (
    "gorm.io/gorm"
    "github.com/your-org/your-service/pkg/api"
)

type AddProducts struct {
    db *gorm.DB
}

func (m *AddProducts) Id() string {
    return "202312150900_add_products"
}

func (m *AddProducts) Up() error {
    return m.db.AutoMigrate(&api.Product{})
}

func (m *AddProducts) Down() error {
    return m.db.Migrator().DropTable(&api.Product{})
}
```

### Running Migrations

```bash
# Run all pending migrations
./your-service migrate

# Check migration status
./your-service migrate --help

# View applied migrations (in database)
make db/login
your-service=# SELECT * FROM migrations;
```

### Migration Best Practices

**✅ Do:**
- Use timestamps in migration names (YYYYMMDDHHMMSS)
- Make migrations reversible with proper `Down()` methods
- Test migrations on sample data
- Add database indexes for performance
- Use GORM's `AutoMigrate` for simple schema changes

**❌ Don't:**
- Modify existing migration files after they're committed
- Delete data in migrations without backups
- Make migrations that can't be rolled back
- Skip the migration registry in `migration_structs.go`

### Custom Migrations

For complex schema changes beyond `AutoMigrate`:

```go
func (m *AddProducts) Up() error {
    // Auto-migrate the basic structure
    err := m.db.AutoMigrate(&api.Product{})
    if err != nil {
        return err
    }
    
    // Add custom indexes
    err = m.db.Exec("CREATE INDEX idx_products_sku ON products(sku)").Error
    if err != nil {
        return err
    }
    
    // Add constraints
    err = m.db.Exec("ALTER TABLE products ADD CONSTRAINT products_price_positive CHECK (price > 0)").Error
    if err != nil {
        return err
    }
    
    return nil
}
```

## Production Database Setup

### Environment Configuration

For production, use environment variables instead of secrets files:

```bash
# Production environment variables
export DB_HOST=your-db-host.amazonaws.com
export DB_PORT=5432
export DB_NAME=your_service_prod
export DB_USER=your_service_user
export DB_PASSWORD=secure_password
export DB_SSLMODE=require
```

### Connection Pooling

TRex configures connection pooling automatically:

```go
// pkg/config/db.go - Connection pool settings
type DatabaseConfig struct {
    MaxOpenConnections int    // Default: 20
    MaxIdleConnections int    // Default: 5
    ConnMaxLifetime    time.Duration  // Default: 5 minutes
}
```

### SSL Configuration

For production databases:

```bash
# Enable SSL
export DB_SSLMODE=require

# SSL certificates (if needed)
export DB_SSLCERT=/path/to/client-cert.pem
export DB_SSLKEY=/path/to/client-key.pem
export DB_SSLROOTCERT=/path/to/ca-cert.pem
```

## Database Operations

### Backup and Restore

```bash
# Backup database
pg_dump -h localhost -U your-service your-service > backup.sql

# Restore database
psql -h localhost -U your-service your-service < backup.sql

# Backup with compression
pg_dump -h localhost -U your-service your-service | gzip > backup.sql.gz
gunzip -c backup.sql.gz | psql -h localhost -U your-service your-service
```

### Database Monitoring

Query database statistics:

```sql
-- Connection count
SELECT count(*) FROM pg_stat_activity WHERE datname = 'your_service';

-- Table sizes
SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables WHERE schemaname = 'public';

-- Migration history
SELECT * FROM migrations ORDER BY created_at;

-- Active queries
SELECT pid, now() - pg_stat_activity.query_start AS duration, query 
FROM pg_stat_activity 
WHERE (now() - pg_stat_activity.query_start) > interval '5 minutes';
```

### Performance Tuning

Common PostgreSQL tuning for TRex services:

```sql
-- Add indexes for common queries
CREATE INDEX idx_products_created_at ON products(created_at);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_name_gin ON products USING gin(to_tsvector('english', name));

-- Analyze table statistics
ANALYZE products;

-- View query plans
EXPLAIN ANALYZE SELECT * FROM products WHERE status = 'active';
```

## Event Sourcing

TRex uses PostgreSQL's NOTIFY/LISTEN for event processing:

### Event Tables

```sql
-- Events table (automatically created)
SELECT * FROM events;

-- Event columns
id         | UUID primary key
source     | Event source (e.g., "Products")
event_type | Event type (e.g., "create", "update", "delete") 
payload    | JSONB event data
created_at | Timestamp
```

### Event Processing

Events are automatically created for CRUD operations:

```go
// When a product is created, an event is automatically generated:
{
    "id": "uuid",
    "source": "Products", 
    "event_type": "create",
    "payload": {"product_id": "123", "name": "Laptop", ...},
    "created_at": "timestamp"
}
```

## Troubleshooting

### Database Won't Start

```bash
# Check for port conflicts
netstat -ln | grep 5432
ss -ln | grep 5432

# Check for existing containers
podman ps -a | grep postgres
docker ps -a | grep postgres

# Stop conflicting services
sudo systemctl stop postgresql  # System PostgreSQL
podman stop existing-container
```

### Migration Failures

```bash
# Check migration status
./your-service migrate

# Manual migration rollback (careful!)
make db/login
your-service=# DELETE FROM migrations WHERE id = 'failed_migration_id';

# Fix migration and try again
# Edit pkg/db/migrations/failed_migration.go
./your-service migrate
```

### Connection Issues

```bash
# Test database connection
psql -h localhost -U your-service -d your-service -c "SELECT 1;"

# Check service logs
./your-service serve  # Look for database connection errors

# Verify database configuration
cat secrets/db.*
```

### Performance Issues

```sql
-- Find slow queries
SELECT query, mean_time, calls, total_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;

-- Check for missing indexes
SELECT schemaname, tablename, attname, n_distinct, correlation
FROM pg_stats
WHERE tablename = 'your_table';
```

## Next Steps

- **[Local Development](local-development.md)** - Complete development setup
- **[Deployment Guide](deployment-guide.md)** - Production database deployment
- **[Troubleshooting](../troubleshooting/runtime-errors.md)** - Database-related runtime issues
- **[Framework Development](../framework-development/)** - Understanding TRex database patterns