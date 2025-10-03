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
3. **Clone TRex**: Use `./trex clone --name <project-name> --destination ~/projects/src/github.com/openshift-online/<project-name>`
4. **Create Installation Status**: Write INSTALLSTATUS.md to cloned project:
   - Create INSTALLSTATUS.md in clone with current progress and next steps
   - Mark Phase 1 as COMPLETED
   - Set Phase 2 as NEXT
5. **Post-Clone Cleanup**: Apply CLONING.md fixes:
   - Remove clone command imports from main.go
   - Fix clone's integration tests. ApiRhTrexV1 references the original TRex project. Use the clone's generated openapi client.
   - Update go.mod module name to match project
   - Run `go mod tidy -C ~/projects/src/github.com/openshift-online/<project-name>`
6. **Database Setup**: Create clone's database:
   - `make -C ~/projects/src/github.com/openshift-online/<project-name> db/setup` - creates PostgreSQL container and database
7. **Verify Clone**: Run `make -C ~/projects/src/github.com/openshift-online/<project-name> binary` - must build successfully
8. **Success**: Only if step 7 succeeds, inform user on next steps in INSTALLSTATUS.md.

**CRITICAL**: Build must succeed for working demo. Follow Rule 3 - reset on failure.

This demonstrates the complete Instant API workflow validation.
