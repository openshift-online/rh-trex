# New Model Generation Plan

## Rules
0. **Critical**: Only successful generation, compilation, and testing. Any test failure is fatal and must be fixed before proceeding.
1. **Validation**: Review model relationships and validate business logic

## Phase 1: Initial Prompt
1. **Prompt**: start new model
2. **Context**: You're a new developer trying to make a model. Be gentle. Make sure every step works.
3. **Inquire**: Ask for data definition
   - Allow an idea for model generation (e.g, "design a data model to track expenses")
   - Use PascalCase
   - Define attributes: Company (Name string required, Address string, StockSymbol string), Employee (Name string required, Email string required, Title string, HiredDate date)
   - Define relationships: Employee N:1 Company (employee.CompanyID, not company.Employees[])
4. **Summarize**: render UML data model, ask for feedback.
5. **Write Model**: After confirmation, update MODEL.md with the new UML data model

## Phase 2: Model Generation
0. **Reminder**: Reminder user that new models are in the current project only (TRex or clone).
   - Track all files generated and code changes made for undo
1. **Generate Entities**: Use TRex generator for each model:
   - `export GOPATH=$HOME/go && export GO111MODULE=on && go run  ./scripts/generator.go --kind $model`
2. **Verify Generation**: Check entities were created in demo project:
   - `ls -la pkg/api/ | grep -E "($model1|$model2)"` - should show 4 .go files
   - `ls -la pkg/services/ | grep -E "($model1)"` - should show 4 service files
4. **Manual Registration**: Complete 4-step post-generator workflow per CLAUDE.md
5. **Populate Fields**: Add business fields from UML model to generated models
6. **OpenAPI Generation**: **CRITICAL** - Generate OpenAPI specs BEFORE migration:
   - `make  generate` - creates OpenAPI bindings
   - This step is REQUIRED before building or running migrations
   - OpenAPI generation creates the required openapi client packages
7. **Database Migration**: Run migration after OpenAPI generation:
   - `make migrate`
   - Ensure database container is running (should be from Phase 2 step 5)
   - Migration creates tables with business fields from populated models

## Phase 4: Testing & Verification
0. **Context**: New models are created and should be fully functional  
1. **Build Verification**: `make  binary` - MUST BUILD SUCCESSFULLY
2. **Test Verification**: Run `make  test` to validate functionality:
   - All tests must pass 
   - Database connection tests validate PostgreSQL integration
   - Service locator tests confirm entity registration
3. **Integration Test Verification**: Run `make  test-integration` to validate full API functionality:
   - Integration tests validate complete REST API workflows
   - Database operations confirmed working
   - All CRUD endpoints functional
   - **CRITICAL**: Database must be running for integration tests to pass
4. **Entity Verification**: All entities generated successfully:
   - New APIs created in ~/projects/src/github.com/openshift-online/<project-name>
   - REST endpoints available at /api/{project-name}/v1/{entity}
   - Full CRUD operations with database integration
   - OpenAPI specifications generated with complete documentation


**CRITICAL**: Build and ALL tests must succeed for working demo. See Rule 0. Use tracked files to undo and reset after failure.

This demonstrates the complete Instant API new model generator.