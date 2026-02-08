# Clone TRex Into a New Microservice

Create a new microservice project by cloning the TRex template and replacing all project-specific names.

## Instructions

Ask the user for:
1. **Service name** (kebab-case): e.g., "my-service", "fleet-manager"
2. **Repository base** (optional, default: `github.com/openshift-online`): e.g., "github.com/myorg"
3. **Destination directory** (optional, default: `/tmp/clone-test`)

### Step 1: Build the TRex Binary

```bash
make binary
```

### Step 2: Run the Clone

```bash
./trex clone --name {service-name} --repo-base {repo-base} --destination {destination}
```

The clone command:
- Walks the entire project tree (excluding `.git/`)
- Replaces all TRex-specific strings via ordered placeholder substitution:
  - `github.com/openshift-online/rh-trex` -> `{repo-base}/{service-name}`
  - `rh-trex` -> `{service-name}`
  - `rhtrex` -> `{service-name}` (no hyphens)
  - `trex` -> `{service-name}` (short form)
  - `TRex` -> `{ServiceName}` (capitalized)
  - `ApiRhTrexV1` -> `Api{ServiceTitleCase}V1`
- Renames files and directories containing "trex"

### Step 3: Post-Clone Setup

The clone command prints next steps. Follow them:

```bash
cd {destination}
go mod tidy
go install gotest.tools/gotestsum@latest
make binary
make db/setup
./{service-name} migrate
make test
make test-integration
```

### Step 4: Verify

```bash
make run-no-auth
curl -s http://localhost:8000/api/{service-name}/v1/dinosaurs | jq
```

The cloned project starts with the Dinosaur Kind as example scaffolding, ready to be replaced with real business entities.

## What Gets Cloned

Everything except `.git/`:
- All Go source code with import paths rewritten
- Makefile with project names updated
- OpenAPI specs with API paths updated
- Templates directory (for generating new Kinds in the cloned project)
- Deployment templates (OpenShift YAML)
- Test infrastructure
- Configuration and secrets structure
