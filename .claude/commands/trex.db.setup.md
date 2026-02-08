# Set Up TRex Database

Start a fresh PostgreSQL container and run all migrations. If a database already exists, this will tear it down first.

## Instructions

1. Check if a PostgreSQL container is already running:
   ```bash
   podman ps -a --filter name=psql-rhtrex --format "{{.Names}} {{.Status}}"
   ```

2. If a container exists, tear it down first:
   ```bash
   make db/teardown
   ```

3. Start a fresh PostgreSQL container:
   ```bash
   make db/setup
   ```

4. Wait for PostgreSQL to be ready:
   ```bash
   sleep 2
   ```

5. Build the binary if not already built:
   ```bash
   make binary
   ```

6. Run migrations:
   ```bash
   ./trex migrate
   ```

## When to Use This

- First-time database setup
- After generating a new Kind (new migration added)
- After pulling changes that include new migrations
- When integration tests fail with "relation does not exist" or column errors
- When the database schema is corrupted or out of sync

## Database Connection Details

Connection details are stored in the `secrets/` directory files (`db.host`, `db.port`, `db.name`, `db.user`, `db.password`).

| Setting | Value |
|---------|-------|
| Host | localhost |
| Port | 5432 |
| Database | rhtrex |
| User | trex |
