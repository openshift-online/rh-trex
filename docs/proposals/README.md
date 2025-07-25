# TRex Proposals Process

## Overview

The proposals process is used to propose new features, architectural changes, or significant modifications to TRex. When proposals are merged, they are considered accepted.

## When to Use

Use proposals for:
- New features or major functionality
- Architectural or breaking changes
- Process modifications
- Significant refactoring

Use regular PRs for bug fixes, minor enhancements, and documentation updates.

## Process

1. **Create**: Copy [proposal-template.md](proposal-template.md) to `{YYYYMMDD}-proposal.md` format (e.g., `20250115-plugin-versioning.md`)
2. **Submit**: Open PR with your proposal in `docs/proposals/`
3. **Review**: Community provides feedback via PR comments
4. **Decision**: Merged (accepted), closed (rejected), or deferred

## Guidelines

- Be specific about problems and solutions
- Consider backward compatibility and impact
- Provide concrete examples and use cases
- Include implementation plan and timeline

## Resources

- Use [proposal-template.md](proposal-template.md) as your starting point
- Add images/diagrams to the [images/](images/) directory
- See [Kubernetes CAPI examples](https://github.com/kubernetes-sigs/cluster-api-provider-aws/tree/main/docs/proposal) for inspiration