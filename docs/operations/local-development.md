# Running TRex

## Prerequisites

Before running TRex for the first time, ensure the prerequisites are installed. For more detailed information on each prerequisite, refer to the [PREREQUISITES.md](./PREREQUISITES.md) document.

## Build and Setup

### 1. Build the Project

```sh
# Install test dependencies
go install gotest.tools/gotestsum@latest

# Build TRex binary
make binary
```

### 2. Database Setup

```sh
# Run PostgreSQL database locally in Docker
make db/setup

# Login to database container
make db/login
```

Verify PostgreSQL is running:
```sh
root@f076ddf94520:/# psql -h localhost -U trex rh-trex
psql (14.4 (Debian 14.4-1.pgdg110+1))
Type "help" for help.

rh-trex=# \dt
Did not find any relations.
```

### 3. Run Database Migrations

The initial migration creates the base data model and provides a way to add future migrations.

```shell
# Run migrations
./trex migrate

# Verify migrations in the database
make db/login
root@f076ddf94520:/# psql -h localhost -U trex rh-trex
psql (14.4 (Debian 14.4-1.pgdg110+1))
Type "help" for help.

rh-trex=# \dt
                 List of relations
 Schema |    Name    | Type  |        Owner
--------+------------+-------+---------------------
 public | dinosaurs  | table | trex
 public | events     | table | trex
 public | migrations | table | trex
(3 rows)
```

## Testing

```shell
# Run unit tests with coverage
make test

# Run integration tests with coverage
make test-integration

# Generate HTML coverage reports
make coverage-html

# View coverage summary in terminal
make coverage-func
```

For comprehensive testing documentation, see the **[Testing Guide](../reference/testing-guide.md)**.

## Running the Service

```shell
# Start the TRex server
make run
```

The server will start on `http://localhost:8000`.

### Verify Server is Running

```shell
curl http://localhost:8000/api/rh-trex/v1/dinosaurs | jq
```

Expected response (401 because authentication is required):
```json
{
  "kind": "Error",
  "id": "401",
  "href": "/api/rh-trex/errors/401",
  "code": "API-401",
  "reason": "Request doesn't contain the 'Authorization' header or the 'cs_jwt' cookie"
}
```

## Authentication

Authentication uses RedHat SSO. You need a Red Hat customer portal user account and can retrieve tokens from https://console.redhat.com/openshift/token

### Using OCM Tool

The ocm tool is available at https://console.redhat.com/openshift/downloads

#### Login to Local Service
```shell
ocm login --token=${OCM_ACCESS_TOKEN} --url=http://localhost:8000
```

#### Test API Access
```shell
# Get all dinosaurs (empty initially)
ocm get /api/rh-trex/v1/dinosaurs
{
  "items": [],
  "kind": "DinosaurList",
  "page": 1,
  "size": 0,
  "total": 0
}

# Create a new dinosaur
ocm post /api/rh-trex/v1/dinosaurs << EOF
{
    "species": "foo"
}
EOF

# Get your dinosaur
ocm get /api/rh-trex/v1/dinosaurs
{
  "items": [
    {
      "created_at":"2023-10-26T08:15:54.509653Z",
      "href":"/api/rh-trex/v1/dinosaurs/2XIENcJIi9t2eBblhWVCtWLdbDZ",
      "id":"2XIENcJIi9t2eBblhWVCtWLdbDZ",
      "kind":"Dinosaur",
      "species":"foo",
      "updated_at":"2023-10-26T08:15:54.509653Z"
    }
  ],
  "kind":"DinosaurList",
  "page":1,
  "size":1,
  "total":1
}
```

## Deployment to OpenShift Local (CRC)

### Prerequisites
Ensure OpenShift Local (CRC) is running:

```shell
crc status
CRC VM:          Running
OpenShift:       Running (v4.13.12)
RAM Usage:       7.709GB of 30.79GB
Disk Usage:      23.75GB of 32.68GB (Inside the CRC VM)
Cache Usage:     37.62GB
Cache Directory: /home/mturansk/.crc/cache
```

### Deploy to CRC

```shell
# Login to CRC
make crc/login
Logging into CRC
Logged into "https://api.crc.testing:6443" as "kubeadmin" using existing credentials.

You have access to 66 projects, the list has been suppressed. You can list all projects with 'oc projects'

Using project "ocm-mturansk".
Login Succeeded!

# Deploy TRex
make deploy

# Login to deployed service
ocm login --token=${OCM_ACCESS_TOKEN} --url=https://trex.apps-crc.testing --insecure

# Test deployed service
ocm post /api/rh-trex/v1/dinosaurs << EOF
{
    "species": "foo"
}
EOF
```