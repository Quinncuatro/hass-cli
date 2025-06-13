TUI Interface Specification
TUI Mode Overview
Interactive terminal user interface for browsing and controlling Home Assistant entities. Designed to be intuitive for users who prefer visual navigation over command-line syntax, with both tree view and dashboard-style layouts available.

Interface Layout Options
Tree View Mode (Default)
Hierarchical display organized by areas and domains:

┌─ Home Assistant CLI v1.0.0 ──────────────────────────────────────────────┐
│ 🏠 homeassistant.local:8123                           [Ctrl+Q] Quit      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│ 📂 Living Room (5)                                                      │
│ ├─ 💡 Lights (3)                                                        │
│ │  ├─ Main Lights        [ON]  85%        [Space] Toggle               │
│ │  ├─ Accent Lights      [OFF]            [Space] Toggle               │
│ │  └─ Table Lamp         [ON]  100%       [Space] Toggle               │
│ ├─ 🌀 Ceiling Fan        [OFF]            [Space] Toggle               │
│ └─ 🌡️ Thermostat          72°F            [Enter] Control               │
│                                                                         │
│ 📂 Bedroom (4)                                                          │
│ ├─ 💡 Lights (2)                                                        │
│ └─ 🌀 Fan                [ON]  Speed 3     [Enter] Control               │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│ Status: Ready | Last Update: 12s ago | [F1] Help | [F2] Dashboard       │
└─────────────────────────────────────────────────────────────────────────┘
Dashboard Mode (F2 to toggle)
Status overview with quick controls:

┌─ Home Assistant Dashboard ───────────────────────────────────────────────┐
│ 🏠 homeassistant.local:8123                           [F2] Tree View     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│ ┌─ Quick Controls ──────┐ ┌─ Climate ────────────────┐ ┌─ Status ──────┐ │
│ │ 💡 All Lights    [ON] │ │ 🌡️ Living Room     72°F │ │ 🔋 Online: 23  │ │
│ │ 🌀 All Fans      [ON] │ │ 🌡️ Bedroom        68°F │ │ ⚠️ Offline: 0  │ │
│ │ 🚪 Garage      [OPEN] │ │ 🌡️ Office         70°F │ │ 🕐 Updated: 8s │ │
│ └───────────────────────┘ └─────────────────────────┘ └───────────────┘ │
│                                                                         │
│ ┌─ Recent Activity ─────────────────────────────────────────────────────┐ │
│ │ 2m ago: Living room lights turned on                                 │ │
│ │ 5m ago: Bedroom fan speed changed to 3                               │ │
│ │ 8m ago: Good Night automation triggered                              │ │
│ └───────────────────────────────────────────────────────────────────────┘ │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│ Status: Ready | [F1] Help | [F2] Tree | [S] Search | [R] Refresh        │
└─────────────────────────────────────────────────────────────────────────┘
Navigation Structure
Three-Panel Layout
Left Panel: Area/Room selector
Center Panel: Entity list for selected area
Right Panel: Entity details (when selected)
Navigation Controls
Arrow Keys: Navigate between panels and items
Tab/Shift+Tab: Switch between panels
Enter: Select/activate item
Space: Toggle entity state (lights, switches)
Escape: Go back/cancel
Ctrl+Q: Quit application
Ctrl+R: Refresh data
F1: Help screen
F5: Force refresh from Home Assistant
Panel Specifications
Left Panel: Area Selector
Purpose: Navigate between areas/rooms
Display: Area name with entity count
Sorting: Alphabetical, with "All" at top
Selection: Highlighted area shows entities in center panel
Special Areas:
"All" - shows all entities
"Automations" - shows all automations
"Scenes" - shows all scenes
"Ungrouped" - entities without area assignment
Center Panel: Entity List
Purpose: Display entities in selected area
Columns:
Icon (emoji representing entity type)
Name (friendly name)
State (current state with color coding)
Value (brightness, temperature, etc.)
Actions (available quick actions)
Sorting: By entity type, then alphabetical
Color Coding:
Green: Active/On states
Red: Inactive/Off states
Yellow: Warning/Unavailable states
Blue: Information/Neutral states
Right Panel: Entity Details (Optional)
Purpose: Detailed view of selected entity
Content:
Full entity information
Available attributes
Recent state changes
Available actions/services
Control widgets for complex entities
Entity Type Representations
Icons and Display
Lights: 💡 with brightness percentage
Switches: 🔌 with on/off state
Sensors: 📊 with current value and unit
Climate: 🌡️ with temperature and mode
Fans: 🌀 with speed percentage
Covers: 🪟 with position percentage
Media Players: 📺 with current activity
Cameras: 📷 with status
Locks: 🔒 with locked/unlocked state
Automations: 🤖 with enabled/disabled state
Scenes: 🎬 with last activated time
Interactive Controls
Quick Actions (Space Bar)
For simple entities that support toggle:

Lights: Toggle on/off
Switches: Toggle on/off
Fans: Toggle on/off
Covers: Toggle open/close
Detailed Controls (Enter Key)
For complex entities requiring parameters:

Lights: Brightness slider, color picker
Climate: Temperature adjustment, mode selection
Media Players: Play/pause, volume, source selection
Covers: Position slider
Automations: Trigger confirmation dialog
Scenes: Activation confirmation dialog
Control Dialogs
Light Control Dialog
┌─ Living Room Main Lights ─────────────────────────────────────────┐
│                                                                   │
│ State: ON                                                         │
│                                                                   │
│ Brightness: [████████████████████████████████████████████▒▒] 85%  │
│ ← → or 0-9 to adjust                                             │
│                                                                   │
│ Color:                                                            │
│ ( ) Warm White  (•) Cool White  ( ) Red  ( ) Blue  ( ) Green     │
│                                                                   │
│ [Turn Off] [Apply] [Cancel]                                       │
└───────────────────────────────────────────────────────────────────┘
Climate Control Dialog
┌─ Living Room Thermostat ──────────────────────────────────────────┐
│                                                                   │
│ Current: 72°F                                                     │
│ Target:  [████████████████████████████████████████████████] 72°F  │
│ ← → to adjust temperature                                         │
│                                                                   │
│ Mode: (•) Heat  ( ) Cool  ( ) Auto  ( ) Off                      │
│                                                                   │
│ [Apply] [Cancel]                                                  │
└───────────────────────────────────────────────────────────────────┘
Status and Feedback
Status Bar (Bottom)
Connection Status: Connected/Disconnected indicator
Last Update: Time since last data refresh
Entity Count: Total entities loaded
Current Action: What's happening now
Help Prompt: Key shortcuts for current context
Loading States
Initial Load: "Loading entities..." with spinner
Refresh: "Refreshing..." with progress indicator
Action Processing: "Turning on lights..." with spinner
Error Handling
Connection Lost: Red status bar with reconnection attempts
Action Failed: Error message in status bar with retry option
Entity Unavailable: Grayed out in list with "Unavailable" state
Keyboard Shortcuts
Global Shortcuts
Ctrl+Q: Quit application
Ctrl+R: Refresh data
F1: Help screen
F5: Force refresh
Tab: Next panel
Shift+Tab: Previous panel
Escape: Cancel/Go back
Navigation Shortcuts
Arrow Keys: Navigate items
Page Up/Down: Navigate by page
Home/End: Go to first/last item
1-9: Quick select (where applicable)
Action Shortcuts
Space: Quick toggle
Enter: Detailed control
T: Toggle (alternative to Space)
R: Refresh current panel
S: Search/filter entities
Help System
Help Screen (F1)
┌─ Home Assistant CLI - Help ───────────────────────────────────────┐
│                                                                   │
│ Navigation:                                                       │
│   ↑↓←→    Navigate items and panels                               │
│   Tab     Switch between panels                                   │
│   Enter   Select/Control entity                                   │
│   Space   Quick toggle (lights, switches)                        │
│   Escape  Go back/Cancel                                          │
│                                                                   │
│ Actions:                                                          │
│   Ctrl+Q  Quit application                                        │
│   Ctrl+R  Refresh data                                            │
│   F5      Force refresh from Home Assistant                       │
│   S       Search entities                                         │
│                                                                   │
│ Entity Types:                                                     │
│   💡 Lights     🔌 Switches    🌡️ Climate    🌀 Fans             │
│   📺 Media      🪟 Covers      📊 Sensors    🤖 Automations       │
│                                                                   │
│ Press any key to close help                                       │
└───────────────────────────────────────────────────────────────────┘
Context-Sensitive Help
Show relevant shortcuts in status bar based on current selection:

Light selected: "Space: Toggle | Enter: Control brightness/color"
Climate selected: "Enter: Adjust temperature and mode"
Automation selected: "Enter: Trigger automation"
Search and Filtering
Search Mode (S key)
Activation: Press 'S' to enter search mode
Input: Type to filter entities by name
Display: Show matching entities across all areas
Exit: Escape to clear search, Enter to select result
Filter Options
By Type: Show only lights, switches, etc.
By State: Show only on/off entities
By Area: Quick area switching
By Recent: Show recently changed entities
Performance Considerations
Refresh Strategy
Background Updates: Refresh entity states every 30 seconds
Smart Refresh: Only update visible entities
User-Triggered: Immediate refresh on Ctrl+R or F5
Connection Recovery: Automatic reconnection with exponential backoff
Memory Usage
Entity Caching: Keep current session entities in memory
Lazy Loading: Load entity details only when needed
Cleanup: Release unused entity data after area changes
