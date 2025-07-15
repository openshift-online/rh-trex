# TRex Cloning Knowledge Base

This document captures all knowledge about TRex cloning behavior, bugs, fixes, and best practices.

## Cloning Overview

The TRex clone command (`./trex clone --name <project-name> --destination <path>`) creates a new microservice project by copying the TRex template and performing string replacements to customize it for the new project.

## Critical Issues Found and Fixed

### Issue 1: Self-Corruption During Clone
**Problem**: The clone command was processing its own source files during the cloning operation, causing the clone command itself to be corrupted with string replacements.

**Symptoms**: 
- Clone command in destination had invalid Go syntax (e.g., `addtrading-apiCloneSection` with hyphens)
- Import paths in clone command were broken
- Subsequent clones would fail

**Root Cause**: The clone logic in `cmd/trex/clone/cmd.go` was applying string replacements to ALL files, including itself.

**Fix Applied**: Added exclusion logic to skip the clone command directory:
```go
// Skip clone command to prevent self-corruption
if strings.Contains(path, "cmd/trex/clone/") || strings.Contains(path, "/clone/cmd.go") {
    return nil
}
```

### Issue 2: Over-Aggressive Import Path Replacement
**Problem**: The original clone logic used broad string replacement that corrupted core library dependencies.

**Original Logic**:
```go
if strings.Contains(content, "github.com/openshift-online/rh-trex") && !strings.Contains(content, "github.com/openshift-online/rh-trex-core") {
    // ... overly broad replacement
}
```

**Problem**: This approach was unreliable and could miss edge cases where both paths appeared in the same file.

**Fix Applied**: More precise import path replacement:
```go
if strings.Contains(content, "github.com/openshift-online/rh-trex/pkg/") {
    glog.Infof("find/replace required for file: %s", path)
    replacement := fmt.Sprintf("%s/%s", provisionCfg.Repo, strings.ToLower(provisionCfg.Name))
    // Replace specific rh-trex package imports, preserving rh-trex-core
    content = strings.Replace(content, "github.com/openshift-online/rh-trex/pkg/", replacement+"/pkg/", -1)
    content = strings.Replace(content, "github.com/openshift-online/rh-trex/cmd/", replacement+"/cmd/", -1)
}
```

### Issue 3: Clone Command Inclusion in Cloned Project
**Problem**: The cloned project included references to the clone command, which doesn't exist in the destination.

**Symptoms**:
- Build failures due to missing clone package imports
- go mod tidy errors trying to fetch non-existent clone packages

**Fix Required**: Manual cleanup after clone:
1. Remove clone import from main.go
2. Remove clone command registration
3. Update go.mod module name
4. Run go mod tidy

## Clone Process Requirements

### Pre-Clone
1. Ensure TRex binary is built: `make binary`
2. Verify source builds successfully

### Clone Command
```bash
./trex clone --name <project-name> --destination <path>
```

### Post-Clone Required Steps
1. **Remove clone command references**:
   ```bash
   # Edit cmd/<project-name>/main.go to remove:
   # - clone import
   # - clone command registration
   ```

2. **Fix module name**:
   ```bash
   # Edit go.mod to change module name from rh-trex to new project name
   ```

3. **Clean dependencies**:
   ```bash
   go mod tidy -C <destination>
   ```

4. **Fix integration test OpenAPI client imports**:
   ```bash
   # Edit test/integration/dinosaurs_test.go to replace:
   # - Import: "github.com/openshift-online/rh-trex/pkg/client/rh-trex/v1"
   # - With: "github.com/openshift-online/<project-name>/pkg/client/<project-name>/v1"
   # 
   # Replace OpenAPI client API method calls:
   # - ApiRhTrexV1* → Api<ProjectName>V1*
   # - /api/rh-trex/v1/ → /api/<project-name>/v1/
   # 
   # Example for project "k8s":
   # - Import: "github.com/openshift-online/k8s/pkg/client/k8s/v1"
   # - ApiRhTrexV1DinosaursGet → ApiK8sV1DinosaursGet
   # - ApiRhTrexV1DinosaursPost → ApiK8sV1DinosaursPost
   # - ApiRhTrexV1DinosaursIdGet → ApiK8sV1DinosaursIdGet
   # - ApiRhTrexV1DinosaursIdPatch → ApiK8sV1DinosaursIdPatch
   ```

5. **Verify build**:
   ```bash
   make -C <destination> binary
   ```

## Import Path Behavior

### Correctly Preserved
- `github.com/openshift-online/rh-trex-core/*` - Always preserved
- External dependencies - Unchanged

### Correctly Replaced
- `github.com/openshift-online/rh-trex/pkg/*` → `<repo>/<project>/pkg/*`
- `github.com/openshift-online/rh-trex/cmd/*` → `<repo>/<project>/cmd/*`

### Manual Module Updates Required
- Module declaration in go.mod
- Clone command cleanup in main.go

## Verification Checklist

After cloning, verify:
- [ ] Project builds successfully (`make binary`)
- [ ] No clone command imports in main.go
- [ ] Correct module name in go.mod
- [ ] Core library imports preserved (rh-trex-core)
- [ ] Project-specific imports updated correctly
- [ ] go mod tidy runs without errors
- [ ] Integration test API references updated (ApiRhTrexV1* → Api<ProjectName>V1*)
- [ ] Integration tests build successfully (`make test-integration`)

## Known Limitations

1. **Clone command is not automatically excluded** - Must be manually removed from cloned project
2. **Module name requires manual update** - go.mod module declaration not automatically updated
3. **Integration test OpenAPI client imports not updated** - test/integration/dinosaurs_test.go still imports TRex OpenAPI client instead of project-specific client, and references ApiRhTrexV1* methods instead of Api<ProjectName>V1*
4. **No validation of successful clone** - Clone command doesn't verify the result builds
5. **Generator pattern matching bug** - Fixed generator has overly broad pattern matching that corrupts all structs ending with "}" instead of targeting specific structures
6. **Incomplete database migrations** - TRex generator creates migration files with only base Model struct, missing all business fields from generated entities

## Best Practices

1. **Always test the clone immediately** after creation with `make binary`
2. **Use descriptive project names** that are valid Go module names
3. **Follow post-clone cleanup steps** religiously
4. **Verify core library dependencies** are preserved
5. **Test entity generation** in cloned project before proceeding

## Testing Clone Quality

```bash
# Complete clone test sequence
./trex clone --name test-project --destination /tmp/test
cd /tmp/test

# Required manual fixes
# 1. Remove clone imports from cmd/test-project/main.go
# 2. Update module in go.mod to github.com/openshift-online/test-project
# 3. Fix integration test OpenAPI client imports in test/integration/dinosaurs_test.go:
#    - Import: "github.com/openshift-online/test-project/pkg/client/test-project/v1"
#    - ApiRhTrexV1DinosaursGet → ApiTestProjectV1DinosaursGet
#    - ApiRhTrexV1DinosaursPost → ApiTestProjectV1DinosaursPost
#    - ApiRhTrexV1DinosaursIdGet → ApiTestProjectV1DinosaursIdGet
#    - ApiRhTrexV1DinosaursIdPatch → ApiTestProjectV1DinosaursIdPatch
go mod tidy
make binary

# Should build successfully without errors
# After OpenAPI generation, integration tests should also build:
make generate
make test-integration
```

## Emergency Debugging

If clone fails:
1. Check for self-corruption in cmd/trex/clone/cmd.go
2. Verify exclusion logic is present
3. Test with clean TRex repository
4. Rebuild TRex binary: `make binary`
5. Validate import path replacement logic

## Generator Service Locator Bug

**Problem**: Fixed generator uses overly broad pattern matching `"}"` which corrupts ALL struct definitions, not just the Services struct.

**Symptoms**: 
- Service locator fields added to ApplicationConfig, Database, Handlers, Clients structs
- Framework.go has service initialization scattered throughout methods
- Build failures due to invalid struct definitions

**Root Cause**: Functions `addServiceLocatorToTypes()` and `addServiceLocatorToFramework()` use generic `"}"` pattern instead of specific struct/method signatures.

**Immediate Fix**: Manual cleanup required:
1. Remove erroneous service locator fields from non-Services structs in types.go
2. Remove scattered service initialization from framework.go except in LoadServices()
3. Verify Services struct has correct service locator fields

**Proper Fix Needed**: Update generator pattern matching to target specific code structures instead of generic closing braces.

## Historical Context

This knowledge was gathered during implementation of Instant API demo for stock trading domain. Multiple clone attempts failed due to self-corruption and import path issues before the fixes were implemented. Rule 3 of INSTANTAPI.md ("Always reset the plan on failure and try again. We want a perfect working demo") drove the systematic debugging and fixing of these issues.