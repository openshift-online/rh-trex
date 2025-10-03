# Data Model

## Database Connection

**PostgreSQL Configuration:**
```bash
# Connection via TRex make commands
make db/login                    # Interactive psql session
make db/setup                    # Start PostgreSQL container
make db/teardown                 # Stop PostgreSQL container

# Direct connection (with timeout)
timeout 10s podman exec psql-rhtrex psql -h localhost -U trex -d rhtrex -c "\dt"

# Connection Details (from secrets/)
Host: localhost
Port: 5432
Database: rhtrex
User: trex
Password: [REDACTED]
Container: psql-rhtrex
```

**Common Database Commands:**
```sql
\dt                              -- List tables
\d table_name                    -- Describe table structure  
\di                              -- List indexes
SELECT * FROM migrations;        -- View migration history
```

## UML Diagram

```
  ┌─────────────────────────┐              ┌─────────────────────────┐
  │       Dinosaur          │              │         Event           │
  │ ─────────────────────── │              │ ─────────────────────── │
  │ + api.Meta              │              │ + api.Meta              │
  │   - id: string          │              │   - id: string          │
  │   - created_at: time    │              │   - created_at: time    │
  │   - updated_at: time    │              │   - updated_at: time    │
  │   - deleted_at: *time   │              │   - deleted_at: *time   │
  │   - kind: string        │              │   - kind: string        │
  │   - href: string        │              │   - href: string        │
  │ + species: string       │              │ + source: string        │
  │                         │              │ + source_id: string     │
  │                         │              │ + event_type: string    │
  │                         │              │ + reconciled_date: time │
  └─────────────────────────┘              └─────────────────────────┘
```

## Entity Definitions

### Dinosaur (Existing)
- **Purpose**: Example entity from TRex template
- **Business Fields**:
  - `species`: Type of dinosaur
- **Constraints**:
  - `species` is indexed for searching

### Event (Existing) 
- **Purpose**: Event-driven architecture support for controllers
- **Business Fields**:
  - `source`: Source entity type (e.g., "Dinosaurs")
  - `source_id`: ID of the source entity that triggered the event
  - `event_type`: Type of event (CREATE, UPDATE, DELETE)
  - `reconciled_date`: When the event was processed
- **Constraints**:
  - `source` is indexed for filtering by entity type
  - `source_id` is indexed for finding events by entity
  - `reconciled_date` is indexed for processing order

## Database Schema

### Tables
- `dinosaurs` - Example dinosaur entities (existing)
- `events` - Event-driven architecture events (existing)
- `migrations` - Database migration tracking (existing)

### Indexes
- `dinosaurs.species` - Search index (`idx_dinosaurs_species`)
- `events.source` - Entity type index (`idx_events_source`)
- `events.source_id` - Entity ID index (`idx_events_source_id`)
- `events.reconciled_date` - Processing order index (`idx_events_reconciled_date`)
- All `deleted_at` fields - Soft delete index (e.g., `idx_dinosaurs_deleted_at`)


## API Endpoints

### Existing Endpoints
- `GET /api/rh-trex/v1/dinosaurs` - List dinosaurs
- `GET /api/rh-trex/v1/dinosaurs/{id}` - Get dinosaur
- `POST /api/rh-trex/v1/dinosaurs` - Create dinosaur
- `PATCH /api/rh-trex/v1/dinosaurs/{id}` - Update dinosaur
- `DELETE /api/rh-trex/v1/dinosaurs/{id}` - Delete dinosaur

### Events (Internal)
- Events are created automatically by CRUD operations
- Processed by event-driven controllers via PostgreSQL LISTEN/NOTIFY
- Not directly exposed as REST endpoints