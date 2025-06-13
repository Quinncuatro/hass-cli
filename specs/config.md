Configuration Management Specification
Configuration File Structure
File Locations
Linux/macOS: ~/.config/hass/config.yaml
Windows: %APPDATA%/hass/config.yaml
Override: --config <path> flag or HASS_CONFIG environment variable
Permissions: 600 (read/write for owner only)
Configuration Schema
yaml
# Home Assistant connection settings
homeassistant:
  url: "http://homeassistant.local:8123"
  token: "your-long-lived-access-token"
  timeout: "10s"
  verify_ssl: true

# User-defined aliases for areas and entities
aliases:
  # Area aliases
  lr: "living room"
  br: "bedroom"
  kit: "kitchen"
  bath: "bathroom"

  # Entity aliases
  coffee: "switch.kitchen_coffee_maker"
  tv: "media_player.living_room_tv"
  thermostat: "climate.main_floor"

# Application preferences
preferences:
  # Default timeout for actions
  default_action_timeout: "5s"

  # Caching settings
  cache_entities: true
  cache_ttl: "5m"
  cache_directory: "~/.cache/hass/"

  # User experience settings
  confirm_destructive: true
  verbose_output: false
  color_output: true

  # TUI settings
  tui_refresh_interval: "30s"
  tui_theme: "default"

# Output formatting preferences
output:
  format: "pretty"  # pretty, json, yaml
  timestamp: false
  show_entity_ids: false

# Discovery settings
discovery:
  enabled: true
  timeout: "5s"
  methods: ["mdns", "common_urls"]

# Security settings
security:
  store_token_in_keyring: true
  keyring_service: "hass-cli"
  allow_insecure_ssl: false
Configuration Commands
Initial Setup (hass config init)
Interactive wizard for first-time configuration:

Welcome to Home Assistant CLI setup!

Step 1: Discover Home Assistant
Searching for Home Assistant instances...
✓ Found: homeassistant.local (192.168.1.100:8123)
✓ Found: hassio.local (192.168.1.101:8123)

Select instance:
  1. homeassistant.local (192.168.1.100:8123)
  2. hassio.local (192.168.1.101:8123)
  3. Enter custom URL
Choice [1]: 1

Step 2: Authentication
You need a long-lived access token from Home Assistant.

How to create a token:
1. Open Home Assistant in your browser
2. Go to Profile → Security → Long-lived access tokens
3. Click "Create Token"
4. Enter a name like "CLI Tool"
5. Copy the token

Enter your token: ********************************

Step 3: Test Connection
Testing connection to homeassistant.local...
✓ Connection successful
✓ Authentication valid
✓ Found 23 entities across 5 areas

Step 4: Optional Setup
Would you like to set up area aliases? [y/N]: y

Enter alias for "Living Room" (or press enter to skip): lr
Enter alias for "Kitchen" (or press enter to skip): kit
Enter alias for "Bedroom" (or press enter to skip): br

Configuration saved to ~/.config/hass/config.yaml
Run 'hass status' to test your setup!
Show Configuration (hass config show)
Display current configuration (with sensitive data masked):

yaml
homeassistant:
  url: "http://homeassistant.local:8123"
  token: "eyJ0eXAi...***MASKED***"
  timeout: "10s"
  verify_ssl: true

aliases:
  lr: "living room"
  kit: "kitchen"
  br: "bedroom"

preferences:
  cache_entities: true
  cache_ttl: "5m"
  confirm_destructive: true
Test Configuration (hass config test)
Validate configuration and test connection:

Testing Home Assistant CLI configuration...

✓ Configuration file found: ~/.config/hass/config.yaml
✓ Configuration syntax valid
✓ Connection to http://homeassistant.local:8123 successful
✓ Authentication valid
✓ API access confirmed
✓ Entity cache refreshed (23 entities loaded)

All tests passed! Your configuration is working correctly.
Configuration Updates
bash
# Set specific configuration values
hass config set homeassistant.url "http://192.168.1.200:8123"
hass config set homeassistant.timeout "15s"
hass config set preferences.cache_ttl "10m"

# Manage aliases
hass config alias add office "home office"
hass config alias remove old_alias
hass config alias list
Environment Variables
Override configuration with environment variables:

HASS_URL: Home Assistant URL
HASS_TOKEN: Access token
HASS_CONFIG: Path to config file
HASS_TIMEOUT: Request timeout
HASS_DEBUG: Enable debug logging
HASS_NO_COLOR: Disable colored output
Security Implementation
Token Storage
Primary: OS keyring (Windows Credential Manager, macOS Keychain, Linux Secret Service)
Fallback: Config file with 600 permissions
Environment: HASS_TOKEN variable (least secure)
Token Validation
Test token on first use and periodically
Handle token expiration gracefully
Prompt for new token when needed
SSL/TLS Handling
Default: Verify SSL certificates
Option: verify_ssl: false for self-signed certificates
Warning messages for insecure connections
Migration and Upgrades
Configuration Migration
Handle configuration format changes between versions:

Backup existing config before migration
Automatic migration of compatible settings
Clear migration messages for breaking changes
Default Configuration Template
yaml
# Default config created when none exists
homeassistant:
  url: ""          # Will prompt for discovery
  token: ""        # Will prompt for token
  timeout: "10s"
  verify_ssl: true

preferences:
  cache_entities: true
  cache_ttl: "5m"
  verbosity: 1     # 0=quiet, 1=normal, 2=verbose
  confirm_destructive: true

aliases: {}        # User-defined shortcuts
Validation Requirements
Config File Validation
Validate YAML syntax
Check required fields are present
Validate URL formats and timeout values
Test Home Assistant connectivity
Verify token authentication
Runtime Configuration
Merge config file, environment variables, and CLI flags
Apply precedence rules (CLI > env > config > defaults)
Validate final configuration before use
