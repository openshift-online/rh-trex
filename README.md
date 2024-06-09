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

```shell

make run

```

To verify that the server is working use the curl command:

```shell

curl http://localhost:8000/api/rh-trex/v1/dinosaurs | jq

```

That should return a 401 response like this, because it needs authentication:

```
{
  "kind": "Error",
  "id": "401",
  "href": "/api/rh-trex/errors/401",
  "code": "API-401",
  "reason": "Request doesn't contain the 'Authorization' header or the 'cs_jwt' cookie"
}
```


Authentication in the default configuration is done through the RedHat SSO, so you need to login with a Red Hat customer portal user in the right account (created as part of the onboarding doc) and then you can retrieve the token to use below on https://console.redhat.com/openshift/token
To authenticate, use the ocm tool against your local service. The ocm tool is available on https://console.redhat.com/openshift/downloads

#### Login to your local service
```
ocm login --token=${OCM_ACCESS_TOKEN} --url=http://localhost:8000

```

#### Confirm login worked by getting all the Dinosaurs
This will be empty if no Dinosaurs exist yet.

Note that we do not use 'curl' here but instead use 'ocm' which passes the user credentials to the API.

```
ocm get /api/rh-trex/v1/dinosaurs
{
  "items": [],
  "kind": "DinosaurList",
  "page": 1,
  "size": 0,
  "total": 0
}
```

#### Post a new Dinosaur

```shell

ocm post /api/rh-trex/v1/dinosaurs << EOF
{
    "species": "foo"
}
EOF

```

#### Get your Dinosaur

```shell
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

#### Run in CRC

Use OpenShift Local to deploy to a local openshift cluster. Be sure to have CRC running locally:

```shell
$ crc status
CRC VM:          Running
OpenShift:       Running (v4.13.12)
RAM Usage:       7.709GB of 30.79GB
Disk Usage:      23.75GB of 32.68GB (Inside the CRC VM)
Cache Usage:     37.62GB
Cache Directory: /home/mturansk/.crc/cache
```

Log into CRC and try a deployment:

```shell

$ make crc/login
Logging into CRC
Logged into "https://api.crc.testing:6443" as "kubeadmin" using existing credentials.

You have access to 66 projects, the list has been suppressed. You can list all projects with 'oc projects'

Using project "ocm-mturansk".
Login Succeeded!

$ make deploy

$ ocm login --token=${OCM_ACCESS_TOKEN} --url=https://trex.apps-crc.testing --insecure

$ ocm post /api/rh-trex/v1/dinosaurs << EOF
{
    "species": "foo"
}
EOF
```



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