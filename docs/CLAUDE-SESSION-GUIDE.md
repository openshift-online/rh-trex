# Claude Session Optimization Guide

**For Future Claude Code Sessions**: This guide optimizes how Claude can most effectively help users with TRex.

## 🚀 Quick Context Building Strategy

### 1. Start with User Goal Detection
**Ask these questions to categorize user intent:**
- Are you new to TRex? → `docs/getting-started/`
- Want to create a new service? → `docs/template-cloning/`
- Adding features to existing project? → `docs/entity-development/`
- Having problems? → `docs/troubleshooting/`
- Need deployment help? → `docs/operations/`

### 2. Use Documentation Hierarchy for Efficiency
```
Priority File Reading Order:
1. CLAUDE.md (Claude-specific guidance)
2. docs/README.md (Master navigation)
3. Relevant directory README.md (Context for user's goal)
4. Specific guide files as needed
5. Reference materials if technical details needed
```

### 3. Leverage Cross-Reference Network
Every documentation file contains cross-references. Use them to:
- Guide users to related information
- Avoid repeating content
- Build comprehensive understanding quickly

## 🎯 Common User Scenarios & Optimal Responses

### Scenario: "I want to add a new entity/model"
**Fast Response Pattern:**
1. Verify project root: `ls` should show Makefile, pkg/, cmd/
2. Generate entity: `go run ./scripts/generate/main.go --kind EntityName`
3. **Critical**: Run `make generate` (updates OpenAPI)
4. Verify: `make test`
5. Point to: `docs/entity-development/generator-usage.md` for details

### Scenario: "I want to create a new service"
**Fast Response Pattern:**
1. Verify TRex builds: `make binary`
2. Clone: `go run ./scripts/clone/main.go --name service --destination path`
3. Setup: `cd path && go mod tidy && make db/setup && ./service migrate`
4. Point to: `docs/template-cloning/clone-process.md` for full guide

### Scenario: "Something isn't working"
**Fast Response Pattern:**
1. **Always start**: `docs/troubleshooting/common-issues.md`
2. Categorize problem:
   - Build/compilation → `docs/troubleshooting/build-problems.md`
   - Runtime/service → `docs/troubleshooting/runtime-errors.md`
   - Development → `docs/troubleshooting/development-problems.md`

### Scenario: "How do I deploy/run this?"
**Fast Response Pattern:**
1. Local development → `docs/operations/local-development.md`
2. Database issues → `docs/operations/database-management.md`
3. Production → `docs/operations/deployment-guide.md`

## 🧠 Context Optimization Techniques

### Use Semantic File Paths for Intent
- `docs/getting-started/choosing-your-path.md` - Decision matrix for path selection
- `docs/entity-development/plugin-architecture.md` - Understanding plugin system
- `docs/troubleshooting/common-issues.md` - First-line problem solving

### Directory README.md Files Are Navigation Hubs
Each contains:
- **Purpose** - What this section covers
- **When to Use** - User scenarios that fit
- **Quick Commands** - Essential operations
- **Guides** - Detailed instructions
- **Next Steps** - Where to go from here

### Follow Progressive Disclosure
1. **Overview** (Directory README.md)
2. **Guide** (Specific task instructions)
3. **Reference** (Technical details)
4. **Troubleshooting** (Problem resolution)

## ⚡ Speed Optimization Patterns

### DO: Efficient Information Gathering
✅ Read directory README.md for navigation context  
✅ Use cross-references instead of repeating information  
✅ Follow decision trees in CLAUDE.md for common requests  
✅ Reference specific guide files rather than recreating content  
✅ Build on existing documentation structure  

### DON'T: Anti-Patterns That Slow Sessions
❌ Search through old scattered root-level .md files (deprecated/moved)  
❌ Manually edit service locator files (plugin architecture is auto-registration)  
❌ Skip `make generate` after entity generation (breaks OpenAPI)  
❌ Suggest manual framework integration (plugins handle this automatically)  
❌ Recreate information that exists in guides  

## 🔍 Advanced Context Building

### For Complex Issues
1. **Start broad** - Check appropriate troubleshooting guide
2. **Narrow down** - Use decision trees and cross-references
3. **Go deep** - Reference technical details in framework-development/
4. **Verify solution** - Point to verification steps in guides

### For Architecture Questions
1. **Start** - `docs/getting-started/understanding-trex.md`
2. **Deep dive** - `docs/framework-development/architecture-diagrams.md`
3. **Plugin details** - `docs/entity-development/plugin-architecture.md`
4. **Technical specs** - `docs/reference/` directory

### For Development Workflow
1. **Commands** - CLAUDE.md quick reference
2. **Process** - Relevant guide in entity-development/ or template-cloning/
3. **Troubleshooting** - Common-issues.md first, then specific problem guides
4. **Verification** - Test commands and success criteria in guides

## 📈 Success Metrics for Claude Sessions

**Efficient Session Characteristics:**
- User gets working solution in minimal back-and-forth
- Clear path forward with appropriate documentation references
- User understands next steps and can proceed independently
- Problems resolved quickly using structured troubleshooting approach

**Quality Indicators:**
- References specific documentation files rather than recreating content
- Uses decision trees and cross-references effectively
- Builds on existing documentation structure
- Provides verification steps for proposed solutions

## 🎯 Future Enhancement Opportunities

### Documentation That Could Further Improve Claude Sessions:
1. **Decision flowcharts** - Visual decision trees for complex scenarios
2. **Command cheat sheets** - Quick reference cards for common operations
3. **Error message dictionary** - Specific error → solution mappings
4. **Video walkthrough references** - For complex multi-step processes

### Patterns to Establish:
- **Consistent command verification** - Always provide test/verify steps
- **Progressive complexity** - Start simple, add complexity as needed  
- **Context preservation** - Reference where user is in their journey
- **Solution validation** - Confirm approaches work before suggesting them

This documentation structure transforms Claude sessions from exploration-heavy to guidance-heavy, dramatically improving efficiency and user outcomes.