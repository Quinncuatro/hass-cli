# Home Assistant CLI Tool - Project Specifications

## Overview

A command-line interface and terminal user interface (TUI) application for managing Home Assistant entities and automations. The tool should provide quick, intuitive commands for common Home Assistant operations.

## Project Structure

```
specs/
├── SPECS.md              # This overview document
├── core.md               # CLI parsing, config management, core application logic
├── homeassistant.md      # Home Assistant API integration and entity management
├── commands.md           # Command definitions and argument parsing
├── tui.md               # Interactive TUI mode specifications
└── config.md            # Configuration file format and management

.project-rules/
├── project-rules-overview.md  # How to use the stdlib approach
├── go-architecture.md         # Package structure and design patterns
├── go-coding-standards.md     # Code style, naming, formatting rules
├── security-rules.md          # Security requirements and best practices
├── go-testing-rules.md        # Testing standards and coverage requirements
├── cli-design-rules.md        # CLI UX principles and error handling
├── project-structure.md       # File organization and build requirements
└── documentation-rules.md     # Documentation and README maintenance requirements
```

## Core Requirements

### Primary Use Cases
1. **Quick entity control**: `hass living lights on` (matches "Living Room")
2. **Automation triggering**: `hass automation "Good Night"`
3. **Entity status checking**: `hass status bedroom`
4. **Interactive mode**: `hass tui` for browsing and controlling entities
5. **Configuration management**: `hass config` for setup and authentication

### Home Assistant Environment
- **Deployment**: Docker container on Proxmox VM (local LAN)
- **Authentication**: Long-lived access tokens
- **Entity Types**: Lights (primary), fans, garage doors, switches, sensors, climate, portable air conditioners
- **Discovery**: Automatic discovery of automations and scenes from HA API

### Technical Requirements
- Single binary deployment (no external dependencies)
- Fast startup time (< 100ms for simple commands)
- Secure credential storage (OS keyring preferred)
- Intelligent caching with local network optimization
- Cross-platform support (Linux, macOS, Windows)
- Comprehensive error handling and user feedback

### User Experience Goals
- **Fuzzy matching**: "living" matches "Living Room" intelligently
- **Area shortcuts**: Support aliases like "lr" for "living room"
- **Ambiguity handling**: Show options when multiple matches exist
- **Verbosity levels**: Three levels (quiet, normal, verbose) with normal as default
- **Smart suggestions**: Context-aware error messages and corrections
- **Comprehensive documentation**: Up-to-date README with installation, configuration, and usage examples

## Success Criteria
- Commands execute in under 500ms on local network
- Zero-configuration discovery of Home Assistant instance
- Handles network failures gracefully
- Provides clear feedback for all operations
- Supports both power users (CLI) and casual users (TUI)
- Documentation enables users to successfully install, configure, and use the tool

## Next Steps
Each specification domain will be detailed in separate markdown files as outlined in the project structure above.