Core Application Specification
Application Name
hass - Home Assistant CLI/TUI tool

Command Line Interface Structure
Root Command
bash
hass [global-flags] <command> [command-args]
Global Flags
--config, -c <path>: Custom config file path
--url <url>: Override Home Assistant URL
--token <token>: Override access token
--timeout <duration>: Request timeout (default: 10s)
--verbose, -v: Verbose output (level 2)
--quiet, -q: Quiet output (level 0)
--debug: Debug logging with full API responses
--help, -h: Show help
--version: Show version information
Verbosity Levels
Level 0 (Quiet): Only essential output and errors
Level 1 (Normal): Standard confirmations and status (default)
Level 2 (Verbose): Detailed information including entity IDs and API responses
Primary Commands
Entity Control: hass <area> <entity-type> <action> [value]
Examples: hass living lights on, hass bedroom fan off, hass kitchen dimmer 50
Automation: hass automation <name>
Example: hass automation "Good Night"
Scene: hass scene <name>
Example: hass scene "Movie Time"
Status: hass status [area|entity]
Examples: hass status, hass status living, hass status "living room lights"
TUI Mode: hass tui
Interactive terminal interface
Configuration: hass config <subcommand>
hass config init: Initial setup wizard
hass config show: Display current configuration
hass config test: Test connection to Home Assistant
Discovery: hass discover
Auto-discover Home Assistant instances on network
Configuration Management
Config File Location
Linux/macOS: ~/.config/hass/config.yaml
Windows: %APPDATA%/hass/config.yaml
Override with --config flag or HASS_CONFIG environment variable
Config File Format
yaml
homeassistant:
  url: "http://homeassistant.local:8123"
  token: "your-long-lived-access-token"
  timeout: "10s"

aliases:
  # Area shortcuts
  lr: "living room"
  br: "bedroom"
  kit: "kitchen"
  bath: "bathroom"
  office: "home office"
  garage: "garage"

preferences:
  default_action_timeout: "5s"
  cache_entities: true
  cache_ttl: "5m"
  verbosity: 1  # 0=quiet, 1=normal, 2=verbose
  confirm_destructive: true
  fuzzy_matching: true
  fuzzy_threshold: 0.6  # Minimum similarity score for matches
Application Architecture
Main Components
CLI Parser: Command parsing and validation using cobra-like library
Config Manager: Configuration loading, validation, and management
HA Client: Home Assistant API client with caching
Entity Resolver: Smart entity name matching and resolution
Command Executor: Execute actions on Home Assistant entities
TUI Interface: Interactive terminal interface
Cache Manager: Entity and state caching for performance
Data Flow
Parse command line arguments
Load configuration
Initialize Home Assistant client
Resolve entity names to entity IDs
Execute command via HA API
Display results to user
Error Handling Strategy
Graceful degradation on network issues
Clear error messages with suggested fixes
Retry logic for transient failures
Offline mode with cached data when possible
Performance Requirements
Cold start: < 100ms for simple commands
Entity resolution: < 50ms
API calls: < 500ms on local network
Cache lookup: < 10ms
Security Considerations
Store tokens securely (OS keyring when available)
Validate all user inputs
Support token refresh if using OAuth
Never log sensitive information
Output Formatting
Success Messages
Concise confirmation: "Living room lights turned on"
Optional verbose mode with entity details
Consistent formatting across commands
Error Messages
Clear problem description
Suggested solutions when possible
Exit codes following UNIX conventions
Status Display
Table format for multiple entities
Color coding for states (green=on, red=off, yellow=unavailable)
Human-readable timestamps and values
