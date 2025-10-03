# Instant API Testing Plan

## Role

You are an expert at creating and managing clones of this TRex project. You always follow the Rules.

## Rules

1. **Always** incorporate working steps back into this plan. The plan should improve itself over time.
2. **Always** provide feedback on making prompts more clear and concise.
3. **Always** reset the plan on failure and try again. We want a perfect working clone.
4. **Always** fix bugs in TRex cloning before proceeding. Any failure must be fixed before proceeding.
5. Allow model reset with Instant API reset.

## Initial Prompt

When user says "start new instantapi", ask:

What specific entity or API would you like to create? I can either:

1. **Create a new entity** using the TRex generator (e.g., `Product`, `User`, `Order`)
2. **Clone TRex** to create a new project for instant API development 
3. **Generate a demo entity** to showcase the instant API capabilities

Please specify the entity name or if you'd prefer to clone TRex first.

## Phase 1: Project Cloning & Setup

0. **Context**: Working in TRex project directory. User has downloaded TRex and is ready to clone.
1. **Container Cleanup**: Check for containers using port 5432:
   - `podman ps` - check for PostgreSQL containers
   - If found warn before: `podman stop <container-name> && podman rm <container-name>`
   - This frees port 5432 for clone's integration tests
2. **Build TRex Binary**: Use `make binary` to build trex command
3. **Clone TRex**: Use enhanced cloner: `go run ./scripts/cloner.go --name <project-name> --destination ~/projects/src/github.com/openshift-online/<project-name>`
4. **Create Installation Status**: Write .trex.md to cloned project:
   - Create .trex.md in clone with current progress and next steps
   - Mark Phase 1 as COMPLETED
   - Set Phase 2 as NEXT
5. **Post-Clone Cleanup**: **AUTOMATED** - Enhanced cloner eliminates manual editing:
   - ✅ Go.mod module declaration automatically updated
   - ✅ API URL paths automatically updated 
   - ✅ Database naming automatically fixed (SQL-safe)
   - ✅ Container names automatically updated
   - ✅ Error codes automatically prefixed
   - ✅ Generator included for entity creation
6. **Dependencies**: Run `go mod tidy -C ~/projects/src/github.com/openshift-online/<project-name>`
7. **Database Setup**: Create clone's database:
   - `make -C ~/projects/src/github.com/openshift-online/<project-name> db/setup` - creates PostgreSQL container and database
8. **Verify Clone**: Run `make -C ~/projects/src/github.com/openshift-online/<project-name> binary` - must build successfully
9. **Success**: Only if step 8 succeeds, inform user on next steps in .trex.md.

**CRITICAL**: Build must succeed for working demo. Follow Rule 3 - reset on failure.

## Phase 2: Entity Generation & Validation

**Context**: Clone is created and working. Now demonstrate instant entity creation.

1. **Generate New Entity**: Use auto-detecting generator:
   - `go run -C ~/projects/src/github.com/openshift-online/<project-name> ./scripts/generator.go --kind <EntityName>`
   - Generator auto-detects project name and creates all CRUD files
   - Look for success message: "✅ Successfully generated <EntityName> entity!"

2. **Auto-Execute Make Generate**: When generator shows "📋 Next step: Run 'make generate'":
   - **IMMEDIATELY** run: `make -C ~/projects/src/github.com/openshift-online/<project-name> generate`
   - This creates OpenAPI models for the new entity
   - Look for "writing file" messages showing model generation

3. **Verify Entity Integration**: 
   - Run: `make -C ~/projects/src/github.com/openshift-online/<project-name> binary` - must build successfully
   - Run: `make -C ~/projects/src/github.com/openshift-online/<project-name> test` - verify tests pass

4. **Success Criteria**:
   - ✅ Entity files created (API, DAO, handlers, services, presenters)
   - ✅ OpenAPI models generated automatically  
   - ✅ Project builds without errors
   - ✅ Tests pass (database tests may fail due to naming, this is expected)

**CRITICAL**: Always run `make generate` immediately after generator success message.

This demonstrates the complete Instant API entity creation workflow.
