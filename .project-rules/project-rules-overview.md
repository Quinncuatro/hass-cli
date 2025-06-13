# Project Rules Overview

## Purpose
This directory contains technical requirements and coding standards that must be followed when implementing the Home Assistant CLI tool. These rules serve as the "stdlib" that guides AI-assisted development.

## Rule Files Structure

### Go-Specific Rules
- `go-architecture.md` - Application structure and design patterns
- `go-coding-standards.md` - Code style, naming, and formatting
- `go-error-handling.md` - Error handling patterns and practices
- `go-testing.md` - Testing requirements and patterns
- `go-dependencies.md` - Dependency management and selection criteria
- `go-performance.md` - Performance considerations and optimization

### Domain-Specific Rules
- `cli-design.md` - Command-line interface design principles
- `security.md` - Security requirements and best practices
- `api-client.md` - HTTP client implementation standards
- `caching.md` - Caching strategy and implementation
- `configuration.md` - Configuration management requirements

### Project Management Rules
- `project-structure.md` - File organization and naming conventions
- `git-workflow.md` - Git commit conventions and branching
- `documentation.md` - Code documentation and README requirements

## Usage Instructions

### For AI Development Tools
When implementing features, always:
1. Study `SPECS.md` for functional requirements
2. Study `.project-rules/` for technical requirements
3. Implement according to both specifications and rules
4. Run tests and linting to validate compliance
5. Update README.md to reflect new features and usage
6. Auto-commit using conventional commit format

### Rule Priority
1. Security rules are non-negotiable
2. Go language best practices must be followed
3. Project structure must be maintained
4. Performance requirements should be met
5. Documentation must be comprehensive

### Validation
Each rule file specifies how compliance is validated:
- Automated: Linting, testing, static analysis
- Manual: Code review checklist items
- Runtime: Performance benchmarks, integration tests

## Updating Rules
Rules should be updated when:
- New security requirements emerge
- Go language best practices evolve
- Project architecture decisions change
- Performance requirements are adjusted

All rule changes must be documented with rationale and migration guidance.