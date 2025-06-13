# Commands Specification

## Command Parsing Strategy

### Natural Language Command Structure
`hass [area] [entity-type] [action] [value]`

Where:
- **area**: Room/area name (optional, can be inferred)
- **entity-type**: Type of entity (lights, switches, etc.)
- **action**: What to do (on, off, toggle, set, etc.)
- **value**: Optional parameter for action (brightness, temperature, etc.)

### Command Resolution Algorithm
1. Parse command line arguments into tokens
2. Resolve area aliases and partial matches
3. Match entity types to Home Assistant domains
4. Validate action against entity capabilities
5. Execute command with appropriate parameters

## Primary Commands

### 1. Entity Control Commands

#### Light Control
```bash
# Basic on/off
hass living lights on
hass bedroom lights off
hass kitchen lights toggle

# Brightness control
hass living lights 50        # 50% brightness
hass living lights 255       # Full brightness (0-255)
hass living lights dim       # Reduce brightness by 25%
hass living lights bright    # Increase brightness by 25%

# Color control
hass living lights red
hass living lights warm
hass living lights cool
hass living lights rgb 255,0,0
```

#### Switch Control
```bash
hass kitchen coffee on
hass living fan off
hass bedroom humidifier toggle
```

#### Fan Control
```bash
hass living fan on
hass bedroom fan off
hass living fan speed 3         # Set speed level (1-5)
hass living fan 75              # Set speed percentage
hass living fan oscillate on    # Enable oscillation
hass living fan oscillate off   # Disable oscillation
```

#### Garage Door Control
```bash
hass garage door open
hass garage door close
hass garage door toggle
hass garage door status         # Check current position
```

#### Climate Control (HVAC & Portable AC)
```bash
hass living temperature 72
hass bedroom heat 70
hass living cool 75
hass living climate auto
hass bedroom ac on              # Portable AC units
hass bedroom ac temperature 68
hass bedroom ac fan high
```

#### Sensor Reading
```bash
hass status temperature         # All temperature sensors
hass bedroom temperature        # Temperature in bedroom
hass garage door sensor         # Garage door position sensor
```

#### Cover Control
```bash
hass living blinds open
hass bedroom curtains close
hass garage door toggle
hass living blinds 50        # 50% open
```

### 2. Automation Commands

#### Trigger Automation
```bash
hass automation "Good Night"
hass automation good_night    # Underscore version
hass automation list          # List all automations
hass automation status "Good Night"  # Check if automation is enabled
```

### 3. Scene Commands

#### Activate Scene
```bash
hass scene "Movie Time"
hass scene movie_time
hass scene list              # List all scenes
```

### 4. Status Commands

#### Entity Status
```bash
hass status                  # Overview of all areas
hass status living           # All entities in living room
hass status lights           # All lights in all areas
hass status living lights    # Just living room lights
hass status "kitchen coffee" # Specific entity
```

#### System Status
```bash
hass status system           # Home Assistant system info
hass status connection       # Connection test
```

### 5. Discovery Commands

#### Network Discovery
```bash
hass discover                # Auto-discover HA instances
hass discover --scan         # Active network scanning
```

### 6. Configuration Commands

#### Setup and Management
```bash
hass config init             # Interactive setup wizard
hass config show             # Display current config
hass config test             # Test connection and auth
hass config set url "http://192.168.1.100:8123"
hass config set token "abc123..."
```

#### Alias Management
```bash
hass config alias add lr "living room"
hass config alias list
hass config alias remove lr
```

### 7. TUI Command

#### Interactive Mode
```bash
hass tui                     # Launch interactive interface
hass tui --area living       # Start in specific area
```

## Command Parsing Logic

### Area Resolution
1. Check exact match against area names from HA Area Registry
2. Check user-defined aliases from config file
3. Fuzzy matching with similarity threshold (0.6)
4. Smart abbreviation matching (lr → living room, br → bedroom)
5. If no area specified, search all areas

### Entity Type Resolution
Map common terms to HA domains:
- `lights`, `light`, `lamp`, `bulb` → `light`
- `switches`, `switch`, `outlet`, `plug` → `switch`
- `fan`, `fans`, `ceiling fan` → `fan`
- `climate`, `thermostat`, `hvac`, `ac`, `air conditioning` → `climate`
- `cover`, `covers`, `garage`, `door`, `blinds`, `curtains` → `cover`
- `sensor`, `sensors`, `temperature`, `humidity` → `sensor`
- `automation`, `automations`, `routine` → `automation`
- `scene`, `scenes` → `scene`

### Action Resolution
Map natural language to HA services:
- `on`, `turn_on`, `enable` → `turn_on`
- `off`, `turn_off`, `disable` → `turn_off`
- `toggle`, `switch` → `toggle`
- `dim`, `dimmer` → reduce brightness
- `bright`, `brighter` → increase brightness
- Numeric values → set brightness/temperature/position

### Value Parsing
- **Brightness**: 0-100 (percentage) or 0-255 (HA native)
- **Temperature**: Numeric with optional unit (°F/°C)
- **Position**: 0-100 (percentage for covers)
- **Colors**: Named colors, hex codes, or rgb(r,g,b)

## Error Handling

### Ambiguous Commands
When multiple entities match:
```
$ hass living lights on
Multiple lights found for "living":
  1. Living Room Main Lights (light.living_room_main)
  2. Living Room Accent Lights (light.living_room_accent)
  3. Living Room Table Lamp (light.living_room_table_lamp)

Choose an option [1-3], or be more specific:
  hass living main lights on
  hass lr accent lights on
  hass living table lamp on
```

Verbosity levels for ambiguity:
- **Quiet (0)**: Show numbered list only
- **Normal (1)**: Show list with friendly names and suggestions
- **Verbose (2)**: Include entity IDs and additional context

### Fuzzy Match Suggestions
```
$ hass livng lights on
No exact match for "livng". Did you mean:
  living room → hass living lights on
  living → hass lr lights on

Available areas: living room, bedroom, kitchen, bathroom
```

### Invalid Commands
```
$ hass bedroom refrigerator on
No entity found matching "bedroom refrigerator"
Available entities in bedroom: lights, fan, climate
```

### Offline Mode
```
$ hass living lights on
Cannot connect to Home Assistant
Using cached data - last updated 2 minutes ago
Warning: Command may not execute
```

## Command Shortcuts

### Quick Actions
- `hass lights on` → Turn on all lights
- `hass lights off` → Turn off all lights
- `hass good_night` → Trigger "Good Night" automation (if exists)
- `hass movie_time` → Activate "Movie Time" scene (if exists)

### Batch Operations
```bash
hass living lights,fan,tv off    # Multiple entities
hass downstairs lights on        # All lights in multiple areas
```

## Output Formatting

### Success Messages

#### Quiet Mode (0)
```
✓ On
✓ 72°F
✓ Triggered
```

#### Normal Mode (1) - Default
```
✓ Living room lights turned on (3 entities)
✓ Bedroom temperature set to 72°F
✓ Good Night automation triggered
```

#### Verbose Mode (2)
```
✓ Living room lights turned on
  - Living Room Main Lights (light.living_room_main): on → 85% brightness
  - Living Room Accent Lights (light.living_room_accent): off → on
  - Living Room Table Lamp (light.living_room_table_lamp): on → on (no change)
  API response time: 234ms
```

### Status Display
```
Living Room Status:
  Lights: On (85% brightness)
  Fan: Off
  Temperature: 72°F (heating)
  TV: Playing Netflix
```

### Error Messages
```
✗ Cannot find entity "bedroom refrigerator"
✗ Home Assistant unreachable (timeout after 10s)
✗ Invalid brightness value "500" (must be 0-255)
```