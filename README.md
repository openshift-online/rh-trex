TRex
---

 **TRex** is RH **T**AP's **R**est **Ex**ample

![Trexxy](rhtap-trex_sm.png)


TRex is a full-featured REST API that persists _dinosaurs_, making it a solid foundation from which developers can quickly bootstrap new services.

Some of the features included are:

* Openapi generation
* CRUD code foundation
* Standard API guidelines, paging, etc.
* Test driven development built-in
* GORM and DB migrations
* OIDC authentication
* Responsive control plane
* Blocking and Non-blocking locks

When looking through the code, anything talking about dinosaurs is business logic, which you
will replace with your business logic. The rest is infrastructure that you will probably want to preserve without change.

It's up to you to port future improvements to this project to your own fork. A goal of this project is to become a
framework with an upgrade path.


## Run for the first time

Before running TRex for the first time, ensure the prerequisites are installed. For more detailed information on each prerequisite, refer to the [prerequisites](./PREREQUISITES.md) document.


### Make a build and run postgres

```sh

# 1. build the project

$ go install gotest.tools/gotestsum@latest
$ make binary

# 2. run a postgres database locally in docker

$ make db/setup
$ make db/login

    root@f076ddf94520:/# psql -h localhost -U trex rh-trex
    psql (14.4 (Debian 14.4-1.pgdg110+1))
    Type "help" for help.

    rh-trex=# \dt
    Did not find any relations.

```

### Run database migrations

The initial migration will create the base data model as well as providing a way to add future migrations.

```shell

# Run migrations
./trex migrate

# Verify they ran in the database
$ make db/login

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

### Test the application

```shell

make test
make test-integration

```

### Running the Service

The service will be available at `http://localhost:8000`

#### Option 1: Run Without Authentication (Recommended for Local Development)

For quick testing and development, you can run the service with authentication disabled:

```shell
make run-no-auth
```

This starts the service with `--enable-authz=false --enable-jwt=false`, allowing you to test the API without tokens.

**Test the API:**

```shell
# List all dinosaurs
curl http://localhost:8000/api/rh-trex/v1/dinosaurs | jq

# Create a new dinosaur
curl -X POST http://localhost:8000/api/rh-trex/v1/dinosaurs \
  -H "Content-Type: application/json" \
  -d '{"species": "Tyrannosaurus"}' | jq

# Get a specific dinosaur (replace {id} with actual ID)
curl http://localhost:8000/api/rh-trex/v1/dinosaurs/{id} | jq
```

#### Option 2: Run With Authentication (Production-like)

Start the service with authentication enabled:

```shell
make run
```

Authentication in the default configuration is done through the RedHat SSO. You need:
- A Red Hat customer portal user in the right account (created as part of the onboarding doc)
- An access token from https://console.redhat.com/openshift/token
- The `ocm` CLI tool available at https://console.redhat.com/openshift/downloads

**Step 1: Login to your local service**

```shell
ocm login --token=${OCM_ACCESS_TOKEN} --url=http://localhost:8000
```

**Step 2: List all Dinosaurs**

```shell
ocm get /api/rh-trex/v1/dinosaurs
```

Response (empty if no dinosaurs exist yet):
```json
{
  "items": [],
  "kind": "DinosaurList",
  "page": 1,
  "size": 0,
  "total": 0
}
```

**Step 3: Create a new Dinosaur**

```shell
ocm post /api/rh-trex/v1/dinosaurs << EOF
{
    "species": "foo"
}
EOF
```

**Step 4: Get your Dinosaur**

```shell
ocm get /api/rh-trex/v1/dinosaurs
```

Response:
```json
{
  "items": [
    {
      "created_at": "2023-10-26T08:15:54.509653Z",
      "href": "/api/rh-trex/v1/dinosaurs/2XIENcJIi9t2eBblhWVCtWLdbDZ",
      "id": "2XIENcJIi9t2eBblhWVCtWLdbDZ",
      "kind": "Dinosaur",
      "species": "foo",
      "updated_at": "2023-10-26T08:15:54.509653Z"
    }
  ],
  "kind": "DinosaurList",
  "page": 1,
  "size": 1,
  "total": 1
}
```

#### Option 3: Deploy to OpenShift Local (CRC)

Use OpenShift Local (CRC) to deploy to a local OpenShift cluster.

**Prerequisites:** Ensure CRC is running locally:

```shell
$ crc status
CRC VM:          Running
OpenShift:       Running (v4.13.12)
RAM Usage:       7.709GB of 30.79GB
Disk Usage:      23.75GB of 32.68GB (Inside the CRC VM)
Cache Usage:     37.62GB
Cache Directory: /home/mturansk/.crc/cache
```

**Deploy to CRC:**

```shell
# 1. Login to CRC
$ make crc/login
Logging into CRC
Logged into "https://api.crc.testing:6443" as "kubeadmin" using existing credentials.

You have access to 66 projects, the list has been suppressed. You can list all projects with 'oc projects'

Using project "ocm-mturansk".
Login Succeeded!

# 2. Deploy the service
$ make deploy

# 3. Login with OCM
$ ocm login --token=${OCM_ACCESS_TOKEN} --url=https://trex.apps-crc.testing --insecure

# 4. Test the deployment
$ ocm post /api/rh-trex/v1/dinosaurs << EOF
{
    "species": "foo"
}
EOF
```

## Run your own service

To create your own service based on TRex, you can use the `clone` command to copy and customize the entire codebase with your service name.

### Clone the code

The clone command will:
- Copy the entire TRex project to a new destination directory
- Replace all occurrences of "trex", "rh-trex", and "TRex" with your new service name
- Update import paths to point to your new repository

```shell
# Build the trex binary first if you haven't built
make binary

# Clone the codebase to create a new service
./trex clone --name my-service --destination /path/to/my-service --repo-base github.com/my-org

# Example:
./trex clone --name rh-birds --destination /tmp/rh-birds --repo-base github.com/openshift-online
```

**Parameters:**
- `--name`: Name of your new service (e.g., "rh-birds", "my-service")
- `--destination`: Directory where the new service code will be created
- `--repo-base`: Your git repository base URL (e.g., "github.com/my-org")

After the clone completes, you'll see a checklist of next steps. Follow these commands to get your new service running:

```shell
# 1. Navigate to your new service directory
cd /path/to/your/new-service

# 2. Install dependencies
go mod tidy

# 3. Build the project
go install gotest.tools/gotestsum@latest
make binary

# 4. Set up the database
make db/setup

# 5. Run database migrations
./your-service-name migrate

# 6. Test the application
make test
make test-integration

# 7. Run your service without authentication required
make run-no-auth

# 8. Verify the service is running
curl http://localhost:8000/api/your-service-name/v1/dinosaurs | jq

# OR 
# 9. Run your service with authentication required
make run

# 10. Verify your application is running with authentication required
curl http://localhost:8000/api/your-service-name/v1/dinosaurs | jq
```


The clone command will output these steps with the correct paths and service names for your convenience.

### Make a new Kind

Generator scripts can be used to auto generate a new Kind. Run the following command to generate a new kind:
```shell
go run ./scripts/generator.go --kind KindName
```

Following manual changes are required to run the application successfully:
- `pkg/api/presenters/kind.go` : Add case statement for the kind
- `pkg/api/presenters/path.go` : Add case statement for the kind
- `pkg/api/presenters/` : Add presenters file (if missing)
- `cmd/trex/environments/service_types.go` : Add new service locator for the kind
- `cmd/trex/environments/types.go` : Add service locator and use `cmd/trex/environments/framework.go` to instantiate
- `cmd/trex/server/routes.go` : Add service routes (if missing)
- Add validation methods in handler if required
- `pkg/db/migrations/migration_structs.go` : Add migration name
- `test/factories.go` : Add helper functions

Here's a reference MR for the same : https://github.com/openshift-online/rh-trex/pull/25
