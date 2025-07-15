# Instant API Testing Plan

## Role

You are an expert API developer with 20 years of experience. Your job is to iterate over this plan until creating
new APIs is Instant and without any errors. Give extra thought as to the user's intent and make suggestions where
appropriate.

## Rules

1. **Always** incorporate working steps back into this plan. Avoid steps that don't work in the future.
2. **Always** provide feedback on making prompts more clear and concise.
3. **Always** reset the plan on failure and try again. We want a perfect working demo.
4. **Always** fix bugs in TRex cloning before proceeding. Any failure must be fixed before proceeding.
5. Allow model and model reset with Instant API reset.

## Phase 1: Project Definition
0. **Context**: Professional demonstration to your boss showing off Instant API™! Everything must always work.
1. **User Input**: Prompt user to describe their project concept and business domain
2. **Model Design**: Create `MODELS.demo.md` with UML diagrams for user's models
3. **Validation**: Review model relationships and validate business logic

## Phase 2: Project Cloning & Setup  
0. **Context**: Clean container environment required for testing
1. **Container Cleanup**: Check for containers using port 5432:
   - `podman ps` - check for PostgreSQL containers
   - If found warn before: `podman stop <container-name> && podman rm <container-name>`
   - This frees port 5432 for clone's integration tests
2. **Build TRex Binary**: Use `make binary` to build trex command
3. **Clone TRex**: Use `./trex clone --name <project-name> --destination ~/projects/src/github.com/openshift-online/<project-name>` 
4. **Post-Clone Cleanup**: Apply CLONING.md fixes:
   - Remove clone command imports from main.go
   - Update go.mod module name to match project
   - Run `go mod tidy -C ~/projects/src/github.com/openshift-online/<project-name>`
5. **Database Setup**: Create clone's database:
   - `make -C ~/projects/src/github.com/openshift-online/<project-name> db/setup` - creates PostgreSQL container and database
6. **Verify Clone**: Run `make -C ~/projects/src/github.com/openshift-online/<project-name> binary` - must build successfully

## Phase 3: Model Generation
0. **Context**: All actions are performed in the cloned project
1. **Copy Models**: Copy @MODELS.demo.md to clone as MODELS.md using `cp /home/mturansk/projects/src/github.com/openshift-online/rh-trex/MODELS.demo.md ~/projects/src/github.com/openshift-online/<project-name>/MODELS.md`
2. **Generate Entities**: Use TRex generator for each model with working directory context:
   - `export GOPATH=$HOME/go && export GO111MODULE=on && go run -C ~/projects/src/github.com/openshift-online/<project-name> ./scripts/generator.go --kind $model`
3. **Verify Generation**: Check entities were created in demo project:
   - `ls -la ~/projects/src/github.com/openshift-online/<project-name>/pkg/api/ | grep -E "($model1|$model2)"` - should show 4 .go files
   - `ls -la ~/projects/src/github.com/openshift-online/<project-name>/pkg/services/ | grep -E "($model1)"` - should show 4 service files
4. **Manual Registration**: Complete 4-step post-generator workflow per CLAUDE.md
5. **Populate Fields**: Add business fields from UML model to generated models
6. **OpenAPI Generation**: **CRITICAL** - Generate OpenAPI specs BEFORE migration:
   - `make -C ~/projects/src/github.com/openshift-online/<project-name> generate` - creates OpenAPI bindings
   - This step is REQUIRED before building or running migrations
   - OpenAPI generation creates the required openapi client packages
7. **Database Migration**: Run migration after OpenAPI generation:
   - `OCM_ENV=development ~/projects/src/github.com/openshift-online/<project-name>/<project-name> migrate`
   - Ensure database container is running (should be from Phase 2 step 5)
   - Migration creates tables with business fields from populated models

## Phase 4: Testing & Verification
0. **Context**: Demo project should now have fully functional API  
1. **Build Verification**: `make -C ~/projects/src/github.com/openshift-online/<project-name> binary` - MUST BUILD SUCCESSFULLY
2. **Test Verification**: Run `make -C ~/projects/src/github.com/openshift-online/<project-name> test` to validate functionality:
   - Most tests should pass (9/10 typical)
   - Database connection tests validate PostgreSQL integration
   - Service locator tests confirm entity registration
3. **Integration Test Verification**: Run `make -C ~/projects/src/github.com/openshift-online/<project-name> test-integration` to validate full API functionality:
   - Integration tests validate complete REST API workflows
   - Database operations confirmed working
   - All CRUD endpoints functional
   - **CRITICAL**: Database must be running for integration tests to pass
4. **Entity Verification**: All entities generated successfully:
   - New APIs created in ~/projects/src/github.com/openshift-online/<project-name>
   - REST endpoints available at /api/{project-name}/v1/{entity}
   - Full CRUD operations with database integration
   - OpenAPI specifications generated with complete documentation

5. **Demo Documentation**: Write final results to demo file:
   - Create `@/demos/{epoch}_{project-name}_demo.md` with complete demonstration results
   - Include project overview, generated entities, API endpoints, and verification status
   - Document any issues encountered and their resolutions
   - Provide summary of Instant API™ workflow performance

**CRITICAL**: Build must succeed for working demo. Follow Rule 3 - reset on failure.

This demonstrates the complete Instant API workflow validation.