# Home Assistant CLI (hass-cli)

A fast, intuitive command-line interface for controlling your Home Assistant instance. Control lights, fans, climate systems, and more using natural language commands.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## Features

- 🏠 **Natural Language Commands**: `hass living lights on`, `hass bedroom temp 72`
- 🔍 **Fuzzy Entity Matching**: "living" automatically matches "Living Room" entities
- ⚡ **Fast Performance**: Sub-500ms response times on local networks
- 🎯 **Smart Discovery**: Auto-detect Home Assistant instances on your network
- 🤖 **Automation Control**: Trigger automations and activate scenes
- 📊 **Status Monitoring**: Check entity states and system status
- 🎨 **Interactive TUI**: Browse and control entities with a terminal interface (coming soon)
- 🔐 **Secure**: Uses Home Assistant's long-lived access tokens

## Quick Start

```bash
# Build the application
make build

# Set up configuration
./bin/hass config init

# Test your connection
./bin/hass config test

# Start controlling your home!
./bin/hass living lights on
```

## Installation

### Prerequisites

- Go 1.21 or later
- Access to a Home Assistant instance

### Build from Source

```bash
git clone https://github.com/quinncuatro/hass-cli.git
cd hass-cli
make build
```

The binary will be created at `./bin/hass`.

### Install to System

```bash
# Install to $GOPATH/bin (make sure it's in your PATH)
make install

# Or copy binary to /usr/local/bin
sudo cp ./bin/hass /usr/local/bin/
```

## Configuration

### 1. Find Your Home Assistant Instance

```bash
# Auto-discover instances on your network
./bin/hass discover
```

Common URLs:
- `http://homeassistant.local:8123`
- `http://hassio.local:8123`
- `http://localhost:8123` (local installation)
- `http://192.168.1.X:8123` (Docker container)

### 2. Create a Long-Lived Access Token

1. Open your Home Assistant web interface
2. Click on your profile icon (bottom left)
3. Scroll down to **"Long-lived access tokens"**
4. Click **"Create Token"**
5. Enter a name like **"CLI Tool"**
6. **Copy the token immediately** (you'll only see it once!)

### 3. Configure hass-cli

#### Option A: Guided Setup
```bash
./bin/hass config init
```

#### Option B: Manual Configuration

Create the config file:

**Linux/macOS:** `~/.config/hass/config.yaml`  
**Windows:** `%APPDATA%/hass/config.yaml`

```yaml
homeassistant:
  url: "http://homeassistant.local:8123"  # Your Home Assistant URL
  token: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9..."  # Your long-lived token
  timeout: "10s"

# Optional: Create shortcuts for areas
aliases:
  lr: "living room"
  br: "bedroom"
  kit: "kitchen"
  bath: "bathroom"

# Optional: Customize behavior
preferences:
  fuzzy_threshold: 0.6  # How strict entity matching should be (0.0-1.0)

output:
  verbosity: 1  # 0=quiet, 1=normal, 2=verbose
  color: true
```

### 4. Test Your Setup

```bash
./bin/hass config test
```

Expected output:
```
Testing Home Assistant Connection...
===================================
✓ Configuration file found
✓ URL configured: http://homeassistant.local:8123
✓ Token configured
✓ Connection successful
✓ Authentication valid
✓ Home Assistant version: 2024.1.0
✓ Location: Home
✓ Found 47 entities

Entity breakdown:
  light: 12
  switch: 8
  sensor: 15
  climate: 2
  ...

✅ All tests passed! Your configuration is working correctly.
```

## Usage

### Entity Control

Control entities using natural language:

```bash
# Lights
hass living lights on                    # Turn on living room lights
hass bedroom lights off                  # Turn off bedroom lights
hass kitchen lights 50                   # Set kitchen lights to 50% brightness
hass living lights bright               # Increase brightness
hass office lights red                   # Set color

# Climate Control
hass living temperature 72               # Set thermostat to 72°F
hass bedroom heat 68                     # Set heating mode to 68°F
hass living ac on                        # Turn on AC

# Fans
hass bedroom fan on                      # Turn on bedroom fan
hass living fan speed 75                 # Set fan to 75% speed
hass kitchen fan oscillate on           # Enable oscillation

# Switches & Outlets
hass kitchen coffee on                   # Turn on coffee maker
hass living lamp off                     # Turn off lamp

# Covers (blinds, garage doors, etc.)
hass living blinds open                  # Open blinds
hass garage door close                   # Close garage door
hass bedroom curtains 50                 # Set curtains to 50% open
```

### Automation & Scenes

```bash
# List available automations
hass automation

# Trigger an automation
hass automation "Good Night"
hass automation good_night               # Underscores work too

# List available scenes
hass scene

# Activate a scene
hass scene "Movie Time"
hass scene relaxing
```

### Status & Information

```bash
# System overview
hass status

# Entity status
hass status living                       # All entities in living room
hass status lights                       # All lights
hass status "bedroom fan"                # Specific entity
hass status bedroom temperature          # Temperature sensors in bedroom

# Configuration
hass config show                         # Show current configuration
hass config test                         # Test connection
```

### Discovery

```bash
# Find Home Assistant instances
hass discover
```

## Command Reference

### Global Flags

```bash
--config, -c <path>     # Custom config file path
--url <url>             # Override Home Assistant URL
--token <token>         # Override access token
--timeout <duration>    # Request timeout (default: 10s)
--verbose, -v           # Verbose output
--quiet, -q             # Quiet output
--help, -h              # Show help
--version               # Show version
```

### Commands

| Command | Description | Examples |
|---------|-------------|----------|
| `config` | Configuration management | `config init`, `config show`, `config test` |
| `status` | Show entity or system status | `status`, `status living`, `status lights` |
| `automation` | List or trigger automations | `automation`, `automation "Good Night"` |
| `scene` | List or activate scenes | `scene`, `scene "Movie Time"` |
| `discover` | Find Home Assistant instances | `discover` |
| `tui` | Interactive terminal interface | `tui` (coming soon) |
| `help` | Show help information | `help` |
| `version` | Show version information | `version` |

### Entity Control Syntax

```bash
hass [area] [entity-type] [action] [value]
```

**Examples:**
- `hass lights on` - Turn on all lights
- `hass living lights on` - Turn on living room lights
- `hass bedroom fan speed 75` - Set bedroom fan to 75%
- `hass kitchen temperature 72` - Set kitchen thermostat to 72°F

## Troubleshooting

### Connection Issues

**Problem:** `Connection refused` or timeout errors

**Solutions:**
- Verify Home Assistant is running and accessible
- Check the URL in your config file
- Try accessing the URL in a web browser
- For Docker containers, use the host machine's IP, not `localhost`
- Ensure port 8123 is open and accessible

**Problem:** `Authentication failed`

**Solutions:**
- Verify your long-lived access token is correct
- Check the token hasn't expired or been revoked
- Create a new token in Home Assistant (Profile → Security → Long-lived access tokens)
- Ensure the token is copied completely without extra spaces

### Entity Not Found

**Problem:** `No entities found matching criteria`

**Solutions:**
- Check entity names in Home Assistant Developer Tools → States
- Try using the exact entity ID: `hass status light.living_room_main`
- Use broader terms: `hass living lights` instead of `hass living room main lights`
- Check if the entity is in the expected area
- Lower the fuzzy threshold in config: `fuzzy_threshold: 0.4`

### Configuration Issues

**Problem:** `Config file not found`

**Solutions:**
- Run `hass config init` to see the expected path
- Create the directory: `mkdir -p ~/.config/hass`
- Check file permissions (should be readable by your user)

**Problem:** `YAML parsing error`

**Solutions:**
- Validate YAML syntax (use a YAML validator online)
- Check indentation (use spaces, not tabs)
- Ensure strings with special characters are quoted

### Discovery Issues

**Problem:** No instances found during discovery

**Solutions:**
- Ensure you're on the same network as Home Assistant
- Try manual configuration with the IP address
- Check firewall settings
- For Docker, ensure the container is on the correct network

## Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Generate coverage report
make coverage

# Run linting
make lint

# Format code
make fmt
```

### Project Structure

```
hass-cli/
├── cmd/hass/           # Application entry point
├── internal/           # Private application packages
│   ├── app/           # Application coordination
│   ├── cli/           # CLI command handlers
│   ├── client/        # Home Assistant API client
│   ├── config/        # Configuration management
│   ├── entity/        # Entity resolution and matching
│   ├── cache/         # Caching implementation
│   └── tui/           # Terminal UI (future)
├── specs/             # Functional specifications
├── .project-rules/    # Development rules and standards
├── Makefile          # Build automation
└── README.md         # This file
```

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes following the coding standards in `.project-rules/`
4. Run tests and linting: `make check`
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/entity -v
```

## Examples

### Common Workflows

**Morning Routine:**
```bash
#!/bin/bash
hass automation "Good Morning"
hass living lights on
hass kitchen coffee on
hass living temperature 70
```

**Evening Control:**
```bash
# Dim all lights to 25%
hass lights 25

# Turn off all unnecessary devices
hass automation "Good Night"

# Check everything is secure
hass status garage
hass status locks
```

**Status Check:**
```bash
# Quick system overview
hass status

# Check specific areas
hass status living
hass status bedroom
hass status security
```

### Integration with Scripts

```bash
#!/bin/bash
# Check if lights are on before leaving
LIGHTS_STATUS=$(hass status lights --quiet)
if [[ $LIGHTS_STATUS == *"on"* ]]; then
    echo "Warning: Some lights are still on!"
    hass lights off
fi
```

## Roadmap

- [ ] Interactive TUI mode with bubbletea
- [ ] Full network discovery with mDNS
- [ ] WebSocket support for real-time updates
- [ ] Keyring integration for secure token storage
- [ ] Plugin system for custom commands
- [ ] Batch operations and scripting support
- [ ] Configuration import/export
- [ ] Shell completion (bash/zsh/fish)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/quinncuatro/hass-cli/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/quinncuatro/hass-cli/discussions)
- 📖 **Documentation**: [GitHub Wiki](https://github.com/quinncuatro/hass-cli/wiki)

## Acknowledgments

- [Home Assistant](https://www.home-assistant.io/) for the amazing home automation platform
- [Cobra](https://github.com/spf13/cobra) for CLI framework inspiration
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for future TUI implementation